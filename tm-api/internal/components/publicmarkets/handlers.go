package publicmarkets

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"admin-panel-dashboard/internal/server/response"

	"github.com/gin-gonic/gin"
)

// Permissions (seeded in internal/db/main.go): reading Public Markets needs
// a Trovo admin; changing it needs PermManage (settings: MANAGE_SETTINGS);
// exchange customers' personal data needs PermPII.
const (
	PermManage   = "MANAGE_PUBLIC_MARKETS"
	PermPII      = "VIEW_PUBLIC_MARKETS_PII"
	PermSettings = "MANAGE_SETTINGS"
)

// Handlers are the admin endpoints. Can reports whether the calling admin
// holds a permission.
type Handlers struct {
	S   *Service
	Can func(c *gin.Context, permission string) bool
}

func (h Handlers) admin(c *gin.Context, permission string) (string, bool) {
	if c.GetString("auth_type") != "trovo_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only Trovo admins can use Public Markets"})
		return "", false
	}
	email := strings.ToLower(strings.TrimSpace(c.GetString("trovo_admin_email")))
	if email == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "the admin's email is unknown"})
		return "", false
	}
	if permission != "" && (h.Can == nil || !h.Can(c, permission)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not have the " + permission + " permission"})
		return "", false
	}
	return email, true
}

func (h Handlers) reader(c *gin.Context) bool            { _, ok := h.admin(c, ""); return ok }
func (h Handlers) manager(c *gin.Context) (string, bool) { return h.admin(c, PermManage) }

func fail(c *gin.Context, err error) {
	var e *Error
	if errors.As(err, &e) {
		c.JSON(e.Status, e)
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func ok(c *gin.Context, message string, data interface{}) {
	response.JSON(c, http.StatusOK, message, data, nil)
}

func page(c *gin.Context) Page {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	return Page{Page: p, Limit: l}
}

func bind(c *gin.Context, v interface{}) bool {
	if err := c.ShouldBindJSON(v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return false
	}
	return true
}

// ReasonRequest carries an action's reason or note.
type ReasonRequest struct {
	Reason string `json:"reason"`
}

func reason(c *gin.Context) (string, bool) {
	var r ReasonRequest
	if c.Request.ContentLength != 0 && !bind(c, &r) {
		return "", false
	}
	return r.Reason, true
}

// ---------------------------------------------------------------- overview & health

// @Summary Public Markets overview
// @Description Asset, order and value totals, what needs attention (approvals, escalations, dead letters, late confirmations, disclosures, halts, drift), the engine's jobs and today's batches.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data{data=Overview}
// @Router /public-markets/overview [get]
func (h Handlers) Overview(c *gin.Context) {
	if h.reader(c) {
		ok(c, "Overview fetched", h.S.Overview())
	}
}

// @Summary Public Markets engine health
// @Description The engine's background jobs (last run, status), recent job requests and queue depths.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data{data=Health}
// @Router /public-markets/health [get]
func (h Handlers) Health(c *gin.Context) {
	if h.reader(c) {
		ok(c, "Health fetched", h.S.Health())
	}
}

// @Summary Inbound partner events
// @Description Webhooks from Custodians and Dealing Members, and what Operations recorded (source MANUAL), with the engine's result.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param source query string false "Source prefix: MANUAL, CUSTODIAN, DEALING_MEMBER"
// @Param kind query string false "settlement, execution, position-feed, corporate-action"
// @Success 200 {object} response.Data
// @Router /public-markets/partner-events [get]
func (h Handlers) PartnerEvents(c *gin.Context) {
	if h.reader(c) {
		rows, total := h.S.PartnerEvents(c.Query("source"), c.Query("kind"), page(c))
		ok(c, "Events fetched", gin.H{"events": rows, "total": total})
	}
}

// ---------------------------------------------------------------- assets

// @Summary Public Market assets
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param market query string false "NGX or FMDQ"
// @Param type query string false "EQUITY or BOND"
// @Param status query string false "SETUP, LIVE or HALTED"
// @Param custodianId query int false "Custodian"
// @Param search query string false "Code, ticker, ISIN or name"
// @Success 200 {object} response.Data{data=[]AssetRow}
// @Router /public-markets/assets [get]
func (h Handlers) Assets(c *gin.Context) {
	if !h.reader(c) {
		return
	}
	cust, _ := strconv.ParseUint(c.Query("custodianId"), 10, 64)
	ok(c, "Assets fetched", h.S.Assets(AssetFilters{Market: c.Query("market"), Type: c.Query("type"), Status: c.Query("status"), Search: c.Query("search"), CustodianID: cust}))
}

// @Summary A Public Market asset
// @Description The asset with supply, Custodian position, beneficial owners (substantial holders flagged), reconciliation runs, prices, corporate actions and its setup steps.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Asset id or code"
// @Success 200 {object} response.Data{data=AssetDetail}
// @Router /public-markets/assets/{id} [get]
func (h Handlers) Asset(c *gin.Context) {
	if !h.reader(c) {
		return
	}
	a, err := h.S.AssetDetail(c.Param("id"))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Asset fetched", a)
}

// @Summary Add a Public Market asset
// @Description Creates the asset in SETUP. Register its token, set a price and record the Custodian position before taking it live.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param body body AssetRequest true "assetCode, market, assetType, isin, instrumentName, custodianId, dealingMemberId, ..."
// @Success 200 {object} response.Data{data=Asset}
// @Router /public-markets/assets [post]
func (h Handlers) CreateAsset(c *gin.Context) {
	by, okk := h.manager(c)
	var r AssetRequest
	if !okk || !bind(c, &r) {
		return
	}
	a, err := h.S.CreateAsset(r, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Asset created", a)
}

// @Summary Edit a Public Market asset
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Asset id or code"
// @Param body body AssetRequest true "Fields to change"
// @Success 200 {object} response.Data{data=Asset}
// @Router /public-markets/assets/{id} [put]
func (h Handlers) UpdateAsset(c *gin.Context) {
	by, okk := h.manager(c)
	var r AssetRequest
	if !okk || !bind(c, &r) {
		return
	}
	a, err := h.S.UpdateAsset(c.Param("id"), r, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Asset updated", a)
}

// ContractRequest registers an asset's token.
type ContractRequest struct {
	ContractAddress    string `json:"contractAddress"`
	IssuingSafeAddress string `json:"issuingSafeAddress"`
}

// @Summary Register an asset's token contract
// @Description Verified on Base: the token's owner must be the issuing Safe and nothing may be minted yet. Its decimals are read from the contract.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Asset id or code"
// @Param body body ContractRequest true "contractAddress and issuingSafeAddress"
// @Success 200 {object} response.Data{data=Asset}
// @Router /public-markets/assets/{id}/contract [post]
func (h Handlers) RegisterContract(c *gin.Context) {
	_, okk := h.manager(c)
	var r ContractRequest
	if !okk || !bind(c, &r) {
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	a, err := h.S.RegisterContract(ctx, c.Param("id"), r.ContractAddress, r.IssuingSafeAddress)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Token contract registered", a)
}

// @Summary Take an asset live
// @Description Opens creation and redemption once the asset's setup steps are done.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Asset id or code"
// @Success 200 {object} response.Data{data=Asset}
// @Router /public-markets/assets/{id}/go-live [post]
func (h Handlers) GoLive(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	a, err := h.S.GoLive(c.Param("id"), by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Asset is live", a)
}

// @Summary Halt an asset
// @Description Stops creation and redemption (holders' balances are not frozen).
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Asset id or code"
// @Param body body ReasonRequest true "reason"
// @Success 200 {object} response.Data{data=Asset}
// @Router /public-markets/assets/{id}/halt [post]
func (h Handlers) Halt(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	why, okk := reason(c)
	if !okk {
		return
	}
	a, err := h.S.Halt(c.Param("id"), why, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Asset halted", a)
}

// @Summary Resume a halted asset
// @Description Asks the engine to reconcile the asset now and reopen it only if the run matches.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Asset id or code"
// @Success 200 {object} response.Data{data=JobRequest}
// @Router /public-markets/assets/{id}/resume [post]
func (h Handlers) Resume(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	j, err := h.S.Resume(c.Param("id"), by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Reconciliation requested; the asset resumes if it matches", j)
}

// PriceRequest is a manual price.
type PriceRequest struct {
	Price string `json:"price"`
}

// @Summary Set a manual reference price
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Asset id or code"
// @Param body body PriceRequest true "price (NGN)"
// @Success 200 {object} response.Data{data=Asset}
// @Router /public-markets/assets/{id}/price [post]
func (h Handlers) SetPrice(c *gin.Context) {
	by, okk := h.manager(c)
	var r PriceRequest
	if !okk || !bind(c, &r) {
		return
	}
	a, err := h.S.SetPrice(c.Param("id"), r.Price, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Price recorded", a)
}

// PositionRequest is a Custodian statement's position.
type PositionRequest struct {
	UnitsHeld string `json:"unitsHeld"`
	AsOf      string `json:"asOf"`
	Reference string `json:"reference"`
}

// @Summary Record the Custodian's position
// @Description From the Custodian's statement (MANUAL Custodians, or a correction). Reconciliation uses the latest position.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Asset id or code"
// @Param body body PositionRequest true "unitsHeld, asOf (YYYY-MM-DD), reference"
// @Success 200 {object} response.Data{data=CustodianPosition}
// @Router /public-markets/assets/{id}/position [post]
func (h Handlers) RecordPosition(c *gin.Context) {
	by, okk := h.manager(c)
	var r PositionRequest
	if !okk || !bind(c, &r) {
		return
	}
	p, err := h.S.RecordPosition(c.Param("id"), r.UnitsHeld, r.AsOf, r.Reference, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Position recorded", p)
}

// @Summary An asset's price history
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Asset id or code"
// @Param limit query int false "Snapshots" default(100)
// @Success 200 {object} response.Data{data=[]PriceSnapshot}
// @Router /public-markets/assets/{id}/prices [get]
func (h Handlers) PriceHistory(c *gin.Context) {
	if !h.reader(c) {
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	rows, err := h.S.PriceHistory(c.Param("id"), limit)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Prices fetched", rows)
}

// @Summary Price oracle
// @Description Each asset's reference price, source and freshness, and recent Dealing Member executions against the reference.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data
// @Router /public-markets/prices [get]
func (h Handlers) Prices(c *gin.Context) {
	if h.reader(c) {
		ok(c, "Prices fetched", gin.H{"prices": h.S.Prices(), "executions": h.S.Executions(50)})
	}
}

// ---------------------------------------------------------------- orders

func dateParam(c *gin.Context, name string, add time.Duration) (*time.Time, bool) {
	v := strings.TrimSpace(c.Query(name))
	if v == "" {
		return nil, true
	}
	t, err := time.ParseInLocation("2006-01-02", v, Lagos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": name + " must be YYYY-MM-DD"})
		return nil, false
	}
	t = t.Add(add).UTC()
	return &t, true
}

// @Summary Public Markets orders
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param type query string false "CREATION or REDEMPTION"
// @Param state query string false "A state, or open"
// @Param channel query string false "TROVO_APP or EXCHANGE"
// @Param path query string false "FAST, SLOW or NETTED"
// @Param asset query string false "Asset code"
// @Param batch query string false "Net batch id"
// @Param exchange query string false "Service link id"
// @Param search query string false "Order id, wallet, username or external ref"
// @Param from query string false "YYYY-MM-DD"
// @Param to query string false "YYYY-MM-DD"
// @Success 200 {object} response.Data
// @Router /public-markets/orders [get]
func (h Handlers) Orders(c *gin.Context) {
	if !h.reader(c) {
		return
	}
	from, okk := dateParam(c, "from", 0)
	if !okk {
		return
	}
	to, okk := dateParam(c, "to", 24*time.Hour)
	if !okk {
		return
	}
	rows, total := h.S.Orders(OrderFilters{Type: c.Query("type"), State: c.Query("state"), Channel: c.Query("channel"), Path: c.Query("path"),
		Asset: c.Query("asset"), BatchID: c.Query("batch"), ServiceLinkID: c.Query("exchange"), Search: c.Query("search"), From: from, To: to, Page: page(c)})
	ok(c, "Orders fetched", gin.H{"orders": rows, "total": total})
}

// @Summary An order with its timeline
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Order id"
// @Success 200 {object} response.Data{data=OrderDetail}
// @Router /public-markets/orders/{id} [get]
func (h Handlers) Order(c *gin.Context) {
	if !h.reader(c) {
		return
	}
	o, err := h.S.OrderDetail(c.Param("id"))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Order fetched", o)
}

// @Summary Net batches
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param status query string false "AWAITING_APPROVAL, RELEASED, EXECUTED, SETTLED, PROCESSED, REJECTED, FAILED, INTERNAL"
// @Param asset query string false "Asset code"
// @Success 200 {object} response.Data
// @Router /public-markets/batches [get]
func (h Handlers) Batches(c *gin.Context) {
	if h.reader(c) {
		rows, total := h.S.Batches(c.Query("status"), c.Query("asset"), page(c))
		ok(c, "Batches fetched", gin.H{"batches": rows, "total": total})
	}
}

// @Summary Approve a net batch above the threshold
// @Description Only Net Creation Approvers; the engine releases the batch once enough have approved.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Batch id"
// @Success 200 {object} response.Data{data=BatchView}
// @Router /public-markets/batches/{id}/approve [post]
func (h Handlers) ApproveBatch(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	b, err := h.S.ApproveBatch(c.Param("id"), by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Approval recorded", b)
}

// @Summary Reject a net batch
// @Description The engine refunds its orders.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Batch id"
// @Param body body ReasonRequest true "reason"
// @Success 200 {object} response.Data{data=NetBatch}
// @Router /public-markets/batches/{id}/reject [post]
func (h Handlers) RejectBatch(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	why, okk := reason(c)
	if !okk {
		return
	}
	b, err := h.S.RejectBatch(c.Param("id"), why, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Batch rejected; its orders are refunded", b)
}

// @Summary Roll a batch to the next session
// @Description For a released batch that was not executed (e.g. the market was closed): its orders wait for the next session's batch.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Batch id"
// @Param body body ReasonRequest true "reason"
// @Success 200 {object} response.Data{data=NetBatch}
// @Router /public-markets/batches/{id}/roll [post]
func (h Handlers) RollBatch(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	why, okk := reason(c)
	if !okk {
		return
	}
	b, err := h.S.RollBatch(c.Param("id"), why, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Batch rolled to the next session", b)
}

// @Summary Outbound instructions
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param status query string false "PENDING, ACCEPTED, EXECUTED, SETTLED, ESCALATED, HANDLED, REJECTED"
// @Param asset query string false "Asset code"
// @Param batch query string false "Batch id"
// @Success 200 {object} response.Data
// @Router /public-markets/instructions [get]
func (h Handlers) Instructions(c *gin.Context) {
	if h.reader(c) {
		rows, total := h.S.Instructions(c.Query("status"), c.Query("asset"), c.Query("batch"), page(c))
		ok(c, "Instructions fetched", gin.H{"instructions": rows, "total": total})
	}
}

// @Summary Retry an escalated instruction
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Instruction id"
// @Success 200 {object} response.Data{data=Instruction}
// @Router /public-markets/instructions/{id}/retry [post]
func (h Handlers) RetryInstruction(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	in, err := h.S.RetryInstruction(c.Param("id"), by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Instruction re-queued", in)
}

// @Summary Mark an instruction handled
// @Description Closes it (no more retries or alerts). Record the partner's outcome so its batch moves on.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Instruction id"
// @Param body body ReasonRequest true "reason: how it was handled"
// @Success 200 {object} response.Data{data=Instruction}
// @Router /public-markets/instructions/{id}/handled [post]
func (h Handlers) MarkHandled(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	why, okk := reason(c)
	if !okk {
		return
	}
	in, err := h.S.MarkHandled(c.Param("id"), why, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Instruction marked handled", in)
}

// @Summary Record a partner's outcome
// @Description For a MANUAL partner (or a webhook that never arrived): a Dealing Member's fill or a Custodian's settlement, applied by the engine as if the partner had sent it.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Instruction id"
// @Param body body OutcomeRequest true "status and the fill / settlement details"
// @Success 200 {object} response.Data{data=PartnerEvent}
// @Router /public-markets/instructions/{id}/outcome [post]
func (h Handlers) RecordOutcome(c *gin.Context) {
	by, okk := h.manager(c)
	var r OutcomeRequest
	if !okk || !bind(c, &r) {
		return
	}
	ev, err := h.S.RecordOutcome(c.Param("id"), r, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Outcome recorded; the engine applies it within seconds", ev)
}

// ---------------------------------------------------------------- reconciliation

// @Summary Reconciliation
// @Description Each asset's latest run, with likely causes of a drift.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data{data=Reconciliation}
// @Router /public-markets/reconciliation [get]
func (h Handlers) Reconciliation(c *gin.Context) {
	if h.reader(c) {
		ok(c, "Reconciliation fetched", h.S.Reconciliation())
	}
}

// @Summary Reconciliation history
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param asset query string false "Asset code"
// @Success 200 {object} response.Data
// @Router /public-markets/reconciliation/runs [get]
func (h Handlers) ReconciliationRuns(c *gin.Context) {
	if h.reader(c) {
		rows, total := h.S.ReconciliationHistory(c.Query("asset"), page(c))
		ok(c, "Runs fetched", gin.H{"runs": rows, "total": total})
	}
}

// TargetRequest names an asset (or ALL) or a Custodian code.
type TargetRequest struct {
	Target string `json:"target"`
}

// @Summary Run reconciliation now
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param body body TargetRequest false "target: an asset code, or ALL"
// @Success 200 {object} response.Data{data=JobRequest}
// @Router /public-markets/reconciliation/run [post]
func (h Handlers) RunReconciliation(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	var r TargetRequest
	if c.Request.ContentLength != 0 && !bind(c, &r) {
		return
	}
	j, err := h.S.RequestReconciliation(r.Target, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Reconciliation requested", j)
}

// @Summary Request a mock Custodian's position feed
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param body body TargetRequest true "target: the Custodian's code"
// @Success 200 {object} response.Data{data=JobRequest}
// @Router /public-markets/position-feed [post]
func (h Handlers) PositionFeed(c *gin.Context) {
	by, okk := h.manager(c)
	var r TargetRequest
	if !okk || !bind(c, &r) {
		return
	}
	j, err := h.S.RequestPositionFeed(r.Target, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Position feed requested", j)
}

// @Summary Job requests
// @Description What was asked of the engine from Trovo Manager, and the results.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data{data=[]JobRequest}
// @Router /public-markets/jobs [get]
func (h Handlers) Jobs(c *gin.Context) {
	if h.reader(c) {
		ok(c, "Job requests fetched", h.S.JobRequests(100))
	}
}

// ---------------------------------------------------------------- corporate actions

// @Summary Corporate actions
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param status query string false "ANNOUNCED, SNAPSHOTTED, APPROVED, PAYING, DISTRIBUTED, NEEDS_MANUAL, CANCELLED"
// @Param asset query string false "Asset code"
// @Success 200 {object} response.Data
// @Router /public-markets/corporate-actions [get]
func (h Handlers) CorporateActions(c *gin.Context) {
	if h.reader(c) {
		rows, total := h.S.CorporateActions(c.Query("status"), c.Query("asset"), page(c))
		ok(c, "Corporate actions fetched", gin.H{"actions": rows, "total": total})
	}
}

// @Summary A corporate action with its distribution
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Corporate action id"
// @Param status query string false "Entitlement status filter"
// @Success 200 {object} response.Data{data=ActionDetail}
// @Router /public-markets/corporate-actions/{id} [get]
func (h Handlers) CorporateAction(c *gin.Context) {
	if !h.reader(c) {
		return
	}
	a, err := h.S.ActionDetail(c.Param("id"), c.Query("status"), page(c))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Corporate action fetched", a)
}

// @Summary Declare a corporate action
// @Description Recorded for the engine, which declares it as it would a Custodian's notice. Dividends and coupons are distributed automatically; bonus, rights and split need manual handling.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param body body DeclareRequest true "assetCode, eventType, recordDate, payDate, amountPerUnit"
// @Success 200 {object} response.Data{data=PartnerEvent}
// @Router /public-markets/corporate-actions [post]
func (h Handlers) Declare(c *gin.Context) {
	by, okk := h.manager(c)
	var r DeclareRequest
	if !okk || !bind(c, &r) {
		return
	}
	ev, err := h.S.Declare(r, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Corporate action recorded; it appears within seconds", ev)
}

// @Summary Approve a distribution
// @Description Only Dividend Approvers; the approval is on the snapshot's checksum.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Corporate action id"
// @Success 200 {object} response.Data{data=ActionDetail}
// @Router /public-markets/corporate-actions/{id}/approve [post]
func (h Handlers) ApproveAction(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	a, err := h.S.ApproveAction(c.Param("id"), by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Approval recorded", a)
}

// @Summary Cancel a corporate action
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Corporate action id"
// @Param body body ReasonRequest true "reason"
// @Success 200 {object} response.Data{data=CorporateAction}
// @Router /public-markets/corporate-actions/{id}/cancel [post]
func (h Handlers) CancelAction(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	why, okk := reason(c)
	if !okk {
		return
	}
	a, err := h.S.CancelAction(c.Param("id"), why, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Corporate action cancelled", a)
}

// @Summary Exchange confirmations of dividend.paid
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data{data=[]ConfirmationRow}
// @Router /public-markets/confirmations [get]
func (h Handlers) Confirmations(c *gin.Context) {
	if h.reader(c) {
		ok(c, "Confirmations fetched", h.S.Confirmations())
	}
}

// ---------------------------------------------------------------- exchanges

// @Summary Exchange partners
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param status query string false "active or suspended"
// @Param search query string false "Name"
// @Success 200 {object} response.Data
// @Router /public-markets/exchanges [get]
func (h Handlers) Exchanges(c *gin.Context) {
	if h.reader(c) {
		ok(c, "Exchanges fetched", gin.H{"exchanges": h.S.Exchanges(c.Query("status"), c.Query("search")), "candidates": h.S.ExchangeCandidates()})
	}
}

// @Summary An exchange partner
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Service link id"
// @Param deliveries query string false "Delivery status filter (PENDING, DELIVERED, DEAD_LETTER)"
// @Success 200 {object} response.Data{data=ExchangeDetail}
// @Router /public-markets/exchanges/{id} [get]
func (h Handlers) Exchange(c *gin.Context) {
	if !h.reader(c) {
		return
	}
	e, err := h.S.ExchangeDetail(c.Param("id"), c.Query("deliveries"))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Exchange fetched", e)
}

// @Summary Onboard an exchange
// @Description Makes a verified service link an exchange partner. The signing secret is in the response once; hand it to the exchange securely.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param body body ExchangeRequest true "serviceLinkId, callbackUrl, fundingAddress, environment, tiers"
// @Success 200 {object} response.Data
// @Router /public-markets/exchanges [post]
func (h Handlers) Onboard(c *gin.Context) {
	by, okk := h.manager(c)
	var r ExchangeRequest
	if !okk || !bind(c, &r) {
		return
	}
	e, secret, err := h.S.Onboard(r, by)
	if err != nil {
		fail(c, err)
		return
	}
	c.Header("Cache-Control", "no-store")
	ok(c, "Exchange onboarded", gin.H{"exchange": e, "signingSecret": secret})
}

// @Summary Edit an exchange partner
// @Description A rate-limit tier change updates the service link's rate limit.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Service link id"
// @Param body body ExchangeRequest true "Fields to change"
// @Success 200 {object} response.Data{data=ExchangeRow}
// @Router /public-markets/exchanges/{id} [put]
func (h Handlers) UpdateExchange(c *gin.Context) {
	by, okk := h.manager(c)
	var r ExchangeRequest
	if !okk || !bind(c, &r) {
		return
	}
	e, err := h.S.UpdateExchange(c.Param("id"), r, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Exchange updated", e)
}

// @Summary Rotate an exchange's signing secret
// @Description The new secret is in the response once; the previous one keeps working for 24 hours.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Service link id"
// @Success 200 {object} response.Data
// @Router /public-markets/exchanges/{id}/rotate-secret [post]
func (h Handlers) RotateSecret(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	secret, until, err := h.S.RotateSecret(c.Param("id"), by)
	if err != nil {
		fail(c, err)
		return
	}
	c.Header("Cache-Control", "no-store")
	ok(c, "Signing secret rotated", gin.H{"signingSecret": secret, "previousValidUntil": until})
}

// @Summary Suspend an exchange
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Service link id"
// @Success 200 {object} response.Data{data=ExchangeRow}
// @Router /public-markets/exchanges/{id}/suspend [post]
func (h Handlers) SuspendExchange(c *gin.Context) { h.exchangeStatus(c, false) }

// @Summary Reactivate an exchange
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Service link id"
// @Success 200 {object} response.Data{data=ExchangeRow}
// @Router /public-markets/exchanges/{id}/activate [post]
func (h Handlers) ActivateExchange(c *gin.Context) { h.exchangeStatus(c, true) }

func (h Handlers) exchangeStatus(c *gin.Context, active bool) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	e, err := h.S.SetExchangeStatus(c.Param("id"), active, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Exchange "+e.Status, e)
}

// AmountRequest carries an amount.
type AmountRequest struct {
	Amount string `json:"amount"`
}

// @Summary Pay an exchange's balance back
// @Description Asks the engine to transfer the amount from the treasury to the exchange's funding wallet.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Service link id"
// @Param body body AmountRequest true "amount"
// @Success 200 {object} response.Data{data=JobRequest}
// @Router /public-markets/exchanges/{id}/withdrawals [post]
func (h Handlers) Withdraw(c *gin.Context) {
	by, okk := h.manager(c)
	var r AmountRequest
	if !okk || !bind(c, &r) {
		return
	}
	j, err := h.S.RequestWithdrawal(c.Param("id"), r.Amount, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Withdrawal requested", j)
}

// @Summary Replay an exchange's dead-lettered webhooks
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Service link id"
// @Success 200 {object} response.Data
// @Router /public-markets/exchanges/{id}/replay-dead-letters [post]
func (h Handlers) ReplayDeadLetters(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	n, err := h.S.ReplayDeadLetters(c.Param("id"), by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Deliveries replayed", gin.H{"replayed": n})
}

// @Summary Replay a dead-lettered webhook
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Delivery id (evt_...)"
// @Success 200 {object} response.Data{data=WebhookDelivery}
// @Router /public-markets/webhooks/{id}/replay [post]
func (h Handlers) ReplayWebhook(c *gin.Context) {
	by, okk := h.manager(c)
	if !okk {
		return
	}
	w, err := h.S.ReplayWebhook(c.Param("id"), by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Delivery replayed", w)
}

// @Summary Wallet provisioning
// @Description Wallets exchanges opened for their customers. Legal names and tax IDs are masked unless the admin holds VIEW_PUBLIC_MARKETS_PII.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Param exchange query string false "Service link id"
// @Param status query string false "active or rejected"
// @Param search query string false "External user ref, wallet id or address"
// @Success 200 {object} response.Data
// @Router /public-markets/wallets [get]
func (h Handlers) Wallets(c *gin.Context) {
	if !h.reader(c) {
		return
	}
	reveal := h.Can != nil && h.Can(c, PermPII)
	rows, total, stats := h.S.Wallets(c.Query("exchange"), c.Query("status"), c.Query("search"), reveal, page(c))
	ok(c, "Wallets fetched", gin.H{"wallets": rows, "total": total, "stats": stats, "personalDataVisible": reveal})
}

// ---------------------------------------------------------------- settings

// @Summary Public Markets settings
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data{data=Settings}
// @Router /public-markets/settings [get]
func (h Handlers) Settings(c *gin.Context) {
	if h.reader(c) {
		ok(c, "Settings fetched", h.S.LoadSettings())
	}
}

// @Summary Change Public Markets settings
// @Description Thresholds, approvers, fees, withholding tax, market hours and holidays, rate-limit tiers. Needs MANAGE_SETTINGS.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param body body SettingsRequest true "Fields to change"
// @Success 200 {object} response.Data{data=Settings}
// @Router /public-markets/settings [put]
func (h Handlers) UpdateSettings(c *gin.Context) {
	by, okk := h.admin(c, PermSettings)
	var r SettingsRequest
	if !okk || !bind(c, &r) {
		return
	}
	st, err := h.S.UpdateSettings(r, by)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Settings saved", st)
}

// @Summary Dealing Members
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data{data=[]DealingMemberRow}
// @Router /public-markets/dealing-members [get]
func (h Handlers) DealingMembers(c *gin.Context) {
	if h.reader(c) {
		ok(c, "Dealing Members fetched", h.S.DealingMembers())
	}
}

// @Summary Add or edit a Dealing Member
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path int false "Dealing Member id (PUT)"
// @Param body body PartnerRequest true "name, code, country, CSCS member code, fees, integration"
// @Success 200 {object} response.Data{data=DealingMember}
// @Router /public-markets/dealing-members [post]
// @Router /public-markets/dealing-members/{id} [put]
func (h Handlers) SaveDealingMember(c *gin.Context) {
	_, okk := h.admin(c, PermSettings)
	var r PartnerRequest
	if !okk || !bind(c, &r) {
		return
	}
	var id uint64
	if v := c.Param("id"); v != "" {
		var err error
		if id, err = strconv.ParseUint(v, 10, 64); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
	}
	dm, err := h.S.SaveDealingMember(id, r)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Dealing Member saved", dm)
}

// @Summary Custodians
// @Description Approved Asset Custodians with their Public Markets integration.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data{data=[]CustodianRow}
// @Router /public-markets/custodians [get]
func (h Handlers) Custodians(c *gin.Context) {
	if h.reader(c) {
		ok(c, "Custodians fetched", h.S.Custodians())
	}
}

// @Summary Configure a Custodian for Public Markets
// @Description Its partner code (X-Partner-Code), nominee, mode (MOCK, REST, MANUAL), REST base URL, auth scheme and credentials reference.
// @Tags Public Markets
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path int true "Approved Asset Custodian id"
// @Param body body PartnerRequest true "Integration settings"
// @Success 200 {object} response.Data{data=Custodian}
// @Router /public-markets/custodians/{id} [put]
func (h Handlers) ConfigureCustodian(c *gin.Context) {
	_, okk := h.admin(c, PermSettings)
	var r PartnerRequest
	if !okk || !bind(c, &r) {
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	cu, err := h.S.ConfigureCustodian(id, r)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, "Custodian saved", cu)
}

// Me reports what the calling admin may do in Public Markets (the pages
// hide what they cannot).
//
// @Summary My Public Markets permissions
// @Tags Public Markets
// @Security JwtTokenAuth
// @Produce json
// @Success 200 {object} response.Data
// @Router /public-markets/me [get]
func (h Handlers) Me(c *gin.Context) {
	email, okk := h.admin(c, "")
	if !okk {
		return
	}
	st := h.S.LoadSettings()
	can := func(p string) bool { return h.Can != nil && h.Can(c, p) }
	ok(c, "Permissions fetched", gin.H{"email": email, "manage": can(PermManage), "settings": can(PermSettings), "personalData": can(PermPII),
		"netCreationApprover": contains(CSV(st.NetCreationApprovers), email), "dividendApprover": contains(CSV(st.DividendApprovers), email)})
}
