package sharedconfig

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The locks live in the database, so goroutines with their own holder
// behave like separate instances; TestKeyLockAcrossProcesses also runs two
// real processes against one database file.

func lockTestDB(t *testing.T, path string) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+path+"?_busy_timeout=10000&_journal_mode=WAL"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&DistributedLock{}, &CallbackDelivery{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSingletonLockIsRenewedWhileWorkRuns(t *testing.T) {
	gc := &GlobalConfig{DB: lockTestDB(t, filepath.Join(t.TempDir(), "locks.db"))}
	var running, overlaps, runs int32
	work := func() {
		if atomic.AddInt32(&running, 1) > 1 {
			atomic.AddInt32(&overlaps, 1)
		}
		atomic.AddInt32(&runs, 1)
		time.Sleep(3500 * time.Millisecond) // much longer than staleAfter
		atomic.AddInt32(&running, -1)
	}
	var wg sync.WaitGroup
	stop := time.Now().Add(5 * time.Second)
	for i := 0; i < 3; i++ { // three "instances" ticking the same job
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(stop) {
				WithSingletonLock(gc, "job", time.Second, work)
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}
	wg.Wait()
	if overlaps != 0 {
		t.Fatalf("the job ran on two instances at once %d times", overlaps)
	}
	if runs < 1 {
		t.Fatal("the job never ran")
	}
}

func TestKeyLockSerializes(t *testing.T) {
	db := lockTestDB(t, filepath.Join(t.TempDir(), "locks.db"))
	var running, overlaps, done int32
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := WithKeyLock(db, "key:safe:0xabc", 30*time.Second, func() error {
				if atomic.AddInt32(&running, 1) > 1 {
					atomic.AddInt32(&overlaps, 1)
				}
				time.Sleep(150 * time.Millisecond)
				atomic.AddInt32(&running, -1)
				atomic.AddInt32(&done, 1)
				return nil
			})
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if overlaps != 0 || done != 6 {
		t.Fatalf("overlaps %d, done %d", overlaps, done)
	}
	// a lock that stays busy for the whole wait is reported, not run
	release := make(chan struct{})
	go WithKeyLock(db, "key:busy", time.Second, func() error { <-release; return nil })
	time.Sleep(300 * time.Millisecond)
	if err := WithKeyLock(db, "key:busy", 500*time.Millisecond, func() error { return nil }); err != ErrLockBusy {
		t.Fatalf("got %v, want ErrLockBusy", err)
	}
	close(release)
}

// TestKeyLockAcrossProcesses runs the critical section from two separate
// processes sharing one database: each appends "start"/"end" markers to a
// file, which must never interleave.
func TestKeyLockAcrossProcesses(t *testing.T) {
	if os.Getenv("LOCK_HELPER_DB") != "" {
		return // the helper itself (see TestLockHelperProcess)
	}
	dir := t.TempDir()
	dbPath, logPath := filepath.Join(dir, "locks.db"), filepath.Join(dir, "log")
	lockTestDB(t, dbPath)
	var cmds []*exec.Cmd
	for i := 0; i < 2; i++ {
		cmd := exec.Command(os.Args[0], "-test.run=^TestLockHelperProcess$")
		cmd.Env = append(os.Environ(), "LOCK_HELPER_DB="+dbPath, "LOCK_HELPER_LOG="+logPath, "LOCK_HELPER_ID="+strconv.Itoa(i))
		cmd.Stdout, cmd.Stderr = os.Stderr, os.Stderr
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		cmds = append(cmds, cmd)
	}
	for _, c := range cmds {
		if err := c.Wait(); err != nil {
			t.Fatal(err)
		}
	}
	raw, _ := os.ReadFile(logPath)
	lines := strings.Fields(string(raw))
	if len(lines) != 2*2*5 {
		t.Fatalf("log has %d entries: %v", len(lines), lines)
	}
	for i := 0; i < len(lines); i += 2 {
		a, b := lines[i], lines[i+1]
		if !strings.HasPrefix(a, "start-") || b != "end-"+strings.TrimPrefix(a, "start-") {
			t.Fatalf("interleaved critical sections at %d: %v", i, lines)
		}
	}
}

func TestLockHelperProcess(t *testing.T) {
	dbPath := os.Getenv("LOCK_HELPER_DB")
	if dbPath == "" {
		t.Skip("helper for TestKeyLockAcrossProcesses")
	}
	db := lockTestDB(t, dbPath)
	id := os.Getenv("LOCK_HELPER_ID")
	appendLine := func(s string) {
		f, err := os.OpenFile(os.Getenv("LOCK_HELPER_LOG"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintln(f, s)
		f.Close()
	}
	for i := 0; i < 5; i++ {
		err := WithKeyLock(db, "key:eoa:0xfunder", time.Minute, func() error {
			appendLine(fmt.Sprintf("start-%s-%d", id, i))
			time.Sleep(100 * time.Millisecond)
			appendLine(fmt.Sprintf("end-%s-%d", id, i))
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestShutdownStopsNewLockedWork(t *testing.T) {
	gc := &GlobalConfig{DB: lockTestDB(t, filepath.Join(t.TempDir(), "locks.db"))}
	defer shuttingDown.Store(false)
	finished := make(chan struct{})
	go WithSingletonLock(gc, "long", time.Second, func() { time.Sleep(500 * time.Millisecond); close(finished) })
	time.Sleep(100 * time.Millisecond)
	BeginShutdown()
	ran := false
	WithSingletonLock(gc, "other", time.Second, func() { ran = true })
	if ran {
		t.Fatal("new work started after shutdown began")
	}
	if !WaitForBackgroundWork(5 * time.Second) {
		t.Fatal("running work was not waited for")
	}
	select {
	case <-finished:
	default:
		t.Fatal("WaitForBackgroundWork returned before the work finished")
	}
}

func TestCallbacksAreRecordedAndRetried(t *testing.T) {
	gc := &GlobalConfig{DB: lockTestDB(t, filepath.Join(t.TempDir(), "cb.db"))}
	var calls int32
	var got []string
	var mu sync.Mutex
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		got = append(got, string(b))
		mu.Unlock()
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusBadGateway) // the partner is down at first
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	gc.SendCallback(srv.URL, []byte(`{"transactionId":"0x1"}`))
	var cb CallbackDelivery
	waitFor(t, func() bool {
		return gc.DB.First(&cb).Error == nil && cb.Attempts == 1
	})
	if cb.Status != CallbackPending || !strings.Contains(cb.LastError, "502") {
		t.Fatalf("after the failed first attempt: %+v", cb)
	}
	// not due yet: nothing is sent
	DeliverDueCallbacks(t.Context(), gc)
	if atomic.LoadInt32(&calls) != 1 {
		t.Fatal("retried before its backoff")
	}
	gc.DB.Model(&CallbackDelivery{}).Where("id = ?", cb.ID).Update("next_attempt_at", time.Now().Add(-time.Second))
	DeliverDueCallbacks(t.Context(), gc)
	gc.DB.First(&cb, "id = ?", cb.ID)
	if cb.Status != CallbackDelivered || cb.Attempts != 2 {
		t.Fatalf("after the retry: %+v", cb)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 || got[0] != `{"transactionId":"0x1"}` || got[1] != got[0] {
		t.Fatalf("bodies sent: %v", got)
	}
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	for i := 0; i < 100; i++ {
		if cond() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("condition not met in time")
}
