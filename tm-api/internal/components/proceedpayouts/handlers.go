package proceedpayouts

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"admin-panel-dashboard/internal/server/response"

	"github.com/gin-gonic/gin"
)

// Handlers are the admin endpoints; Trovo admins only.
type Handlers struct{ S *Service }

func adminOnly(c *gin.Context) (string, bool) {
	if c.GetString("auth_type") != "trovo_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only Trovo admins can manage payouts"})
		return "", false
	}
	email := strings.TrimSpace(c.GetString("trovo_admin_email"))
	if email == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "the admin's email is unknown"})
		return "", false
	}
	return email, true
}

func fail(c *gin.Context, err error) {
	var e *Error
	if errors.As(err, &e) {
		c.JSON(e.Status, gin.H{"error": e.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func payoutID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payout id"})
		return 0, false
	}
	return id, true
}

func filters(c *gin.Context) (Filters, bool) {
	f := Filters{Status: c.Query("status"), AssetID: c.Query("asset")}
	f.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	f.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "50"))
	for _, d := range []struct {
		name string
		dst  **time.Time
		add  time.Duration
	}{{"from", &f.From, 0}, {"to", &f.To, 24 * time.Hour}} {
		if v := strings.TrimSpace(c.Query(d.name)); v != "" {
			t, err := time.Parse("2006-01-02", v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": d.name + " must be YYYY-MM-DD"})
				return f, false
			}
			t = t.Add(d.add) // "to" includes its day
			*d.dst = &t
		}
	}
	return f, true
}

// ReasonRequest carries an action's reason.
type ReasonRequest struct {
	Reason string `json:"reason"`
}

// ReferenceRequest carries how a holder was paid outside the engine.
type ReferenceRequest struct {
	Reference string `json:"reference"`
}

// SweepRequest names the token to sweep from the payout Safe.
type SweepRequest struct {
	Token string `json:"token" binding:"required"`
}

// @Summary List proceeds payouts
// @Description Payouts of authorized stakeholder distributions, newest first.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Produce json
// @Param status query string false "Status (REGISTERED, LOCKED, APPROVED, PAYING, PAUSED, COMPLETED, ...)"
// @Param asset query string false "Tokenized asset ID or code"
// @Param from query string false "Created on or after (YYYY-MM-DD)"
// @Param to query string false "Created on or before (YYYY-MM-DD)"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Page size (max 200)" default(50)
// @Success 200 {object} response.Data{data=[]PayoutView}
// @Router /proceed-payouts [get]
func (h Handlers) List(c *gin.Context) {
	if _, ok := adminOnly(c); !ok {
		return
	}
	f, ok := filters(c)
	if !ok {
		return
	}
	rows, total, err := h.S.List(f)
	if err != nil {
		fail(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	response.JSON(c, http.StatusOK, "Payouts fetched", gin.H{"payouts": rows, "total": total}, nil)
}

// @Summary A proceeds payout
// @Description The payout with its approvals, batches, holder counts and the engine's state.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Produce json
// @Param id path int true "Payout ID"
// @Success 200 {object} response.Data{data=Detail}
// @Router /proceed-payouts/{id} [get]
func (h Handlers) Detail(c *gin.Context) {
	if _, ok := adminOnly(c); !ok {
		return
	}
	id, ok := payoutID(c)
	if !ok {
		return
	}
	d, err := h.S.Detail(id)
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Payout fetched", d, nil)
}

// @Summary A payout's schedule
// @Description The schedule's lines: the fee and VAT, then holders by amount.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Produce json
// @Param id path int true "Payout ID"
// @Param status query string false "PENDING, QUEUED, PAID, FAILED, EXCLUDED or SKIPPED"
// @Param kind query string false "HOLDER, FEE or VAT"
// @Param search query string false "Address or username"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Page size (max 200)" default(50)
// @Success 200 {object} response.Data{data=[]PayoutItem}
// @Router /proceed-payouts/{id}/items [get]
func (h Handlers) Items(c *gin.Context) {
	if _, ok := adminOnly(c); !ok {
		return
	}
	id, ok := payoutID(c)
	if !ok {
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	rows, total, err := h.S.Items(id, ItemFilters{Status: c.Query("status"), Kind: c.Query("kind"), Search: c.Query("search"), Page: page, Limit: limit})
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Schedule fetched", gin.H{"items": rows, "total": total}, nil)
}

// payoutAction runs an action on the payout in the path and answers the
// payout after it.
func (h Handlers) payoutAction(c *gin.Context, message string, act func(id uint64, admin string) (*ProceedPayout, error)) {
	admin, ok := adminOnly(c)
	if !ok {
		return
	}
	id, ok := payoutID(c)
	if !ok {
		return
	}
	p, err := act(id, admin)
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, message, s0(h.S.views([]ProceedPayout{*p})), nil)
}

func s0(v []PayoutView) PayoutView { return v[0] }

func reason(c *gin.Context) (string, bool) {
	var r ReasonRequest
	if c.Request.ContentLength != 0 {
		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return "", false
		}
	}
	return r.Reason, true
}

// @Summary Prepare (or re-prepare) a payout's schedule
// @Description payout-engine snapshots the asset's holders (transfers during preparation included) and locks the schedule for approval. A locked or approved schedule is prepared again and its approvals dropped. The admin who prepares cannot approve.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Param id path int true "Payout ID"
// @Success 200 {object} response.Data{data=PayoutView}
// @Router /proceed-payouts/{id}/prepare [post]
func (h Handlers) Prepare(c *gin.Context) {
	h.payoutAction(c, "Schedule preparation requested", h.S.Prepare)
}

// @Summary Set a payout's processing fee
// @Description FIXED (an amount of the payout token) or PERCENT of the payout (capped by feeCap when above 0). VAT at the asset country's rate is charged on the fee. The fee and VAT are paid from the payout to the fee and VAT wallets. On a locked or approved payout the schedule is prepared again, so the fee is approved with it; the admin who sets the fee cannot approve.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Accept json
// @Param id path int true "Payout ID"
// @Param body body FeeRequest true "Fee"
// @Success 200 {object} response.Data{data=PayoutView}
// @Router /proceed-payouts/{id}/fee [put]
func (h Handlers) SetFee(c *gin.Context) {
	var r FeeRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.payoutAction(c, "Fee set", func(id uint64, admin string) (*ProceedPayout, error) { return h.S.SetFee(id, admin, r) })
}

// @Summary Approve a locked schedule
// @Description Each admin approves once; the payout is approved when the required number of distinct admins (PROCEED_PAYOUT_APPROVALS_REQUIRED, default 2) approved this schedule.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Param id path int true "Payout ID"
// @Success 200 {object} response.Data{data=PayoutView}
// @Router /proceed-payouts/{id}/approve [post]
func (h Handlers) Approve(c *gin.Context) {
	h.payoutAction(c, "Approval recorded", h.S.Approve)
}

// @Summary Reject a schedule
// @Description Sends a locked or approved payout back to REGISTERED and drops its approvals.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Accept json
// @Param id path int true "Payout ID"
// @Param body body ReasonRequest true "Reason"
// @Success 200 {object} response.Data{data=PayoutView}
// @Router /proceed-payouts/{id}/reject [post]
func (h Handlers) Reject(c *gin.Context) {
	r, ok := reason(c)
	if !ok {
		return
	}
	h.payoutAction(c, "Schedule rejected", func(id uint64, admin string) (*ProceedPayout, error) { return h.S.Reject(id, admin, r) })
}

// @Summary Confirm the payout is funded
// @Description payout-engine checks the payout Safe holds what is still to pay (on top of other payouts in progress), that the signers can sign and the executor has gas, then starts paying; otherwise the payout returns to APPROVED with the shortfall in its note.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Param id path int true "Payout ID"
// @Success 200 {object} response.Data{data=PayoutView}
// @Router /proceed-payouts/{id}/confirm-funding [post]
func (h Handlers) ConfirmFunding(c *gin.Context) {
	h.payoutAction(c, "Funding check requested", h.S.ConfirmFunding)
}

// @Summary Pause a payout
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Accept json
// @Param id path int true "Payout ID"
// @Param body body ReasonRequest false "Reason"
// @Success 200 {object} response.Data{data=PayoutView}
// @Router /proceed-payouts/{id}/pause [post]
func (h Handlers) Pause(c *gin.Context) {
	r, ok := reason(c)
	if !ok {
		return
	}
	h.payoutAction(c, "Payout paused", func(id uint64, admin string) (*ProceedPayout, error) { return h.S.Pause(id, admin, r) })
}

// @Summary Resume a paused payout
// @Description The funding check runs again before paying continues.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Param id path int true "Payout ID"
// @Success 200 {object} response.Data{data=PayoutView}
// @Router /proceed-payouts/{id}/resume [post]
func (h Handlers) Resume(c *gin.Context) {
	h.payoutAction(c, "Payout resumed", h.S.Resume)
}

// @Summary Cancel a payout
// @Description Stops the payout for good; holders already paid stay paid. Not while a batch is being mined.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Accept json
// @Param id path int true "Payout ID"
// @Param body body ReasonRequest true "Reason"
// @Success 200 {object} response.Data{data=PayoutView}
// @Router /proceed-payouts/{id}/cancel [post]
func (h Handlers) Cancel(c *gin.Context) {
	r, ok := reason(c)
	if !ok {
		return
	}
	h.payoutAction(c, "Payout cancelled", func(id uint64, admin string) (*ProceedPayout, error) { return h.S.Cancel(id, admin, r) })
}

// @Summary Retry failed transfers
// @Description Puts a completed payout's failed transfers back in the schedule; they are paid after a new funding check.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Param id path int true "Payout ID"
// @Success 200 {object} response.Data{data=PayoutView}
// @Router /proceed-payouts/{id}/retry-failed [post]
func (h Handlers) RetryFailed(c *gin.Context) {
	h.payoutAction(c, "Failed transfers queued again", h.S.RetryFailed)
}

func (h Handlers) itemAction(c *gin.Context, message string, act func(id uint64, itemID, admin string) (*PayoutItem, error)) {
	admin, ok := adminOnly(c)
	if !ok {
		return
	}
	id, ok := payoutID(c)
	if !ok {
		return
	}
	it, err := act(id, c.Param("itemId"), admin)
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, message, it, nil)
}

// @Summary Exclude a holder
// @Description Leaves a pending or failed holder out of the payout; its share stays in the payout Safe. Approvals stand.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Accept json
// @Param id path int true "Payout ID"
// @Param itemId path string true "Schedule item ID"
// @Param body body ReasonRequest true "Reason"
// @Success 200 {object} response.Data{data=PayoutItem}
// @Router /proceed-payouts/{id}/items/{itemId}/exclude [post]
func (h Handlers) ExcludeItem(c *gin.Context) {
	r, ok := reason(c)
	if !ok {
		return
	}
	h.itemAction(c, "Holder excluded", func(id uint64, itemID, admin string) (*PayoutItem, error) {
		return h.S.ExcludeItem(id, itemID, admin, r)
	})
}

// @Summary Include an excluded holder again
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Param id path int true "Payout ID"
// @Param itemId path string true "Schedule item ID"
// @Success 200 {object} response.Data{data=PayoutItem}
// @Router /proceed-payouts/{id}/items/{itemId}/include [post]
func (h Handlers) IncludeItem(c *gin.Context) {
	h.itemAction(c, "Holder included", h.S.IncludeItem)
}

// @Summary Mark a holder paid
// @Description Records that a holder the engine could not pay was paid another way. The payout must be paused or completed with failures.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Accept json
// @Param id path int true "Payout ID"
// @Param itemId path string true "Schedule item ID"
// @Param body body ReferenceRequest true "How it was paid"
// @Success 200 {object} response.Data{data=PayoutItem}
// @Router /proceed-payouts/{id}/items/{itemId}/mark-paid [post]
func (h Handlers) MarkPaid(c *gin.Context) {
	var r ReferenceRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.itemAction(c, "Holder marked paid", func(id uint64, itemID, admin string) (*PayoutItem, error) {
		return h.S.MarkPaid(id, itemID, admin, r.Reference)
	})
}

// @Summary payout-engine's state
// @Description Heartbeat (online), current activity, last error, kill switch and pending sweep.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data{data=EngineStatus}
// @Router /proceed-payouts/engine [get]
func (h Handlers) Engine(c *gin.Context) {
	if _, ok := adminOnly(c); !ok {
		return
	}
	st, err := h.S.Engine()
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Engine state fetched", st, nil)
}

// @Summary Stop payout-engine (kill switch)
// @Description The engine stops all payout work until it is resumed.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Accept json
// @Param body body ReasonRequest true "Reason"
// @Success 200 {object} response.Data{data=EngineStatus}
// @Router /proceed-payouts/engine/halt [post]
func (h Handlers) Halt(c *gin.Context) {
	admin, ok := adminOnly(c)
	if !ok {
		return
	}
	r, ok := reason(c)
	if !ok {
		return
	}
	st, err := h.S.Halt(admin, r)
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Payout engine halted", st, nil)
}

// @Summary Let payout-engine work again
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Success 200 {object} response.Data{data=EngineStatus}
// @Router /proceed-payouts/engine/unhalt [post]
func (h Handlers) Unhalt(c *gin.Context) {
	admin, ok := adminOnly(c)
	if !ok {
		return
	}
	st, err := h.S.Unhalt(admin)
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Payout engine resumed", st, nil)
}

// @Summary Sweep a token from the payout Safe
// @Description Moves the payout Safe's whole balance of the token to the engine's sweep address (PROCEED_PAYOUT_SWEEP_ADDRESS); refused while a payout of that token is being funded or paid.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Accept json
// @Param body body SweepRequest true "Token contract"
// @Success 200 {object} response.Data{data=EngineStatus}
// @Router /proceed-payouts/engine/sweep [post]
func (h Handlers) Sweep(c *gin.Context) {
	admin, ok := adminOnly(c)
	if !ok {
		return
	}
	var r SweepRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	st, err := h.S.Sweep(admin, r.Token)
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Sweep requested", st, nil)
}

// @Summary The payout fee configuration
// @Description The payout fee wallet and the default fee of new payouts (none set: FIXED 0, cap 0), and the VAT wallet.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data{data=FeeConfig}
// @Router /proceed-payouts/fee-config [get]
func (h Handlers) FeeConfig(c *gin.Context) {
	if _, ok := adminOnly(c); !ok {
		return
	}
	cfg, err := h.S.FeeConfig()
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Payout fee configuration fetched", cfg, nil)
}

// @Summary Set the payout fee configuration
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Accept json
// @Param body body FeeConfigRequest true "Fee wallet and defaults"
// @Success 200 {object} response.Data{data=FeeConfig}
// @Router /proceed-payouts/fee-config [put]
func (h Handlers) SetFeeConfig(c *gin.Context) {
	admin, ok := adminOnly(c)
	if !ok {
		return
	}
	var r FeeConfigRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cfg, err := h.S.SetFeeConfig(admin, r)
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Payout fee configuration saved", cfg, nil)
}

// @Summary Payouts report
// @Description Matching payouts with totals per payout currency (authorized, to holders, paid, fees, VAT) and counts per status.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Produce json
// @Param status query string false "Status"
// @Param asset query string false "Tokenized asset ID or code"
// @Param from query string false "From (YYYY-MM-DD)"
// @Param to query string false "To (YYYY-MM-DD)"
// @Success 200 {object} response.Data{data=PayoutsReport}
// @Router /proceed-payouts/reports/payouts [get]
func (h Handlers) PayoutsReport(c *gin.Context) {
	if _, ok := adminOnly(c); !ok {
		return
	}
	f, ok := filters(c)
	if !ok {
		return
	}
	r, err := h.S.PayoutsReport(f)
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Payouts report", r, nil)
}

// @Summary Payout fees and VAT report
// @Description Payout processing fees and the VAT on them as paid to the fee and VAT wallets, with totals per asset and currency.
// @Tags Proceeds payouts
// @Security JwtTokenAuth
// @Produce json
// @Param asset query string false "Tokenized asset ID or code"
// @Param from query string false "From (YYYY-MM-DD)"
// @Param to query string false "To (YYYY-MM-DD)"
// @Success 200 {object} response.Data{data=FeesReport}
// @Router /proceed-payouts/reports/fees [get]
func (h Handlers) FeesReport(c *gin.Context) {
	if _, ok := adminOnly(c); !ok {
		return
	}
	f, ok := filters(c)
	if !ok {
		return
	}
	r, err := h.S.FeesReport(f)
	if err != nil {
		fail(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Payout fees report", r, nil)
}
