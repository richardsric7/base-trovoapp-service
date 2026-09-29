package aa

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// safeSim applies owner-management calls with Safe v1.4.1's rules
// (linked list with new owners at the head, prevOwner checks, threshold
// bounds) so PlanOwnerChange's output can be checked without a chain.
type safeSim struct {
	owners    []common.Address
	threshold int64
}

func (s *safeSim) apply(t *testing.T, safe common.Address, c Call) {
	t.Helper()
	if c.To != safe {
		t.Fatalf("call to %s, not the Safe", c.To.Hex())
	}
	m, err := contractsABI.MethodById(c.Data[:4])
	if err != nil {
		t.Fatal(err)
	}
	args, _ := m.Inputs.Unpack(c.Data[4:])
	switch m.Name {
	case "addOwnerWithThreshold":
		o, th := args[0].(common.Address), args[1].(*big.Int).Int64()
		for _, x := range s.owners {
			if x == o {
				t.Fatalf("GS204 duplicate owner %s", o.Hex())
			}
		}
		s.owners = append([]common.Address{o}, s.owners...)
		s.setThreshold(t, th)
	case "removeOwner":
		prev, o, th := args[0].(common.Address), args[1].(common.Address), args[2].(*big.Int).Int64()
		idx := -1
		for i, x := range s.owners {
			if x == o {
				idx = i
			}
		}
		if idx < 0 {
			t.Fatalf("GS205 %s not an owner", o.Hex())
		}
		want := sentinelOwners
		if idx > 0 {
			want = s.owners[idx-1]
		}
		if prev != want {
			t.Fatalf("GS205 wrong prevOwner for %s", o.Hex())
		}
		if int64(len(s.owners)-1) < th {
			t.Fatalf("GS201 threshold %d > %d owners", th, len(s.owners)-1)
		}
		s.owners = append(s.owners[:idx:idx], s.owners[idx+1:]...)
		s.setThreshold(t, th)
	case "changeThreshold":
		s.setThreshold(t, args[0].(*big.Int).Int64())
	default:
		t.Fatalf("unexpected %s", m.Name)
	}
}

func (s *safeSim) setThreshold(t *testing.T, th int64) {
	t.Helper()
	if th < 1 || th > int64(len(s.owners)) {
		t.Fatalf("GS201/202 threshold %d with %d owners", th, len(s.owners))
	}
	s.threshold = th
}

func TestPlanOwnerChange(t *testing.T) {
	safe := common.HexToAddress("0x5afe")
	a, b, c, d := common.HexToAddress("0xa"), common.HexToAddress("0xb"), common.HexToAddress("0xc"), common.HexToAddress("0xd")
	cases := []struct {
		name          string
		current       []common.Address
		curTh         int64
		target        []common.Address
		targetTh      int64
		wantCallCount int
	}{
		{"enable approvers", []common.Address{a}, 1, []common.Address{a, b, c}, 2, 3},
		{"only threshold", []common.Address{a, b, c}, 2, []common.Address{a, b, c}, 3, 1},
		{"swap an approver", []common.Address{c, b, a}, 2, []common.Address{a, b, d}, 2, 2},
		{"disable", []common.Address{c, b, a}, 2, []common.Address{a}, 1, 2},
		{"no change", []common.Address{b, a}, 2, []common.Address{a, b}, 2, 0},
	}
	for _, tc := range cases {
		calls, err := PlanOwnerChange(safe, tc.current, tc.curTh, tc.target, tc.targetTh)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if len(calls) != tc.wantCallCount {
			t.Fatalf("%s: %d calls, want %d", tc.name, len(calls), tc.wantCallCount)
		}
		sim := &safeSim{owners: append([]common.Address(nil), tc.current...), threshold: tc.curTh}
		for _, call := range calls {
			sim.apply(t, safe, call)
		}
		if OwnersHash(sim.owners, sim.threshold) != OwnersHash(tc.target, tc.targetTh) {
			t.Fatalf("%s: ended with %v/%d", tc.name, sim.owners, sim.threshold)
		}
	}
	if _, err := PlanOwnerChange(safe, []common.Address{a}, 1, []common.Address{a}, 2); err == nil {
		t.Fatal("an impossible threshold must be refused")
	}
}

func TestViaModuleAndExtraModules(t *testing.T) {
	cfg := ConfigFromEnv(big.NewInt(8453))
	dist, issuing := common.HexToAddress("0xd157"), common.HexToAddress("0x1551")
	call := ViaModule(dist, AddOwner(dist, common.HexToAddress("0xb"), 1))
	m, _ := contractsABI.MethodById(call.Data[:4])
	args, _ := m.Inputs.Unpack(call.Data[4:])
	if call.To != dist || m.Name != "execTransactionFromModule" || args[0].(common.Address) != dist || args[3].(uint8) != 0 {
		t.Fatalf("ViaModule %s %v", m.Name, args)
	}
	owners := []common.Address{common.HexToAddress("0xa")}
	if cfg.SafeAddress(owners, 1, big.NewInt(1)) == cfg.SafeAddress(owners, 1, big.NewInt(1), issuing) {
		t.Fatal("extra modules must change the Safe address")
	}
	key, err := RandomNonceKey()
	if err != nil || key.BitLen() > 192 {
		t.Fatalf("nonce key %v %v", key, err)
	}
}
