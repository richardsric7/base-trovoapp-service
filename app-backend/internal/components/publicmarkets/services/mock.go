package publicmarkets

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	pm "trovo-wallet-api/internal/components/publicmarkets/models"
	"trovo-wallet-api/internal/components/publicmarkets/partners"

	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// PartnerFactory builds the client of a Custodian or Dealing Member from
// its configuration. A nil client with no error means MANUAL: Operations
// act on the partner's portal and record the outcome in Trovo Manager.
type PartnerFactory interface {
	Custodian(c pm.Custodian) (partners.Custodian, error)
	DealingMember(dm pm.DealingMember) (partners.DealingMember, error)
}

// DefaultPartners picks mock, REST or manual per partner.
type DefaultPartners struct{ E *Engine }

func (f DefaultPartners) Custodian(c pm.Custodian) (partners.Custodian, error) {
	switch strings.ToUpper(c.Mode) {
	case pm.ModeManual:
		return nil, nil
	case pm.ModeREST:
		rc, err := partners.NewRESTClient(c.BaseURL, c.AuthScheme, c.CredentialsRef)
		if err != nil {
			return nil, err
		}
		return partners.RESTCustodian{RESTClient: rc}, nil
	default:
		return MockCustodian{E: f.E, Code: c.Code}, nil
	}
}

func (f DefaultPartners) DealingMember(dm pm.DealingMember) (partners.DealingMember, error) {
	switch strings.ToUpper(dm.Mode) {
	case pm.ModeManual:
		return nil, nil
	case pm.ModeREST:
		rc, err := partners.NewRESTClient(dm.BaseURL, dm.AuthScheme, dm.CredentialsRef)
		if err != nil {
			return nil, err
		}
		return partners.RESTDealingMember{RESTClient: rc}, nil
	default:
		return MockDealingMember{E: f.E, Code: dm.Code}, nil
	}
}

func envSeconds(key string, def int) time.Duration {
	if v, err := strconv.Atoi(strings.TrimSpace(os.Getenv(key))); err == nil && v >= 0 {
		return time.Duration(v) * time.Second
	}
	return time.Duration(def) * time.Second
}

// MockCustodian stands in for a Custodian until its API is agreed (CSCS is
// only ever reached through the Custodian). It acknowledges instructions,
// keeps its own book of units held and confirms settlement after
// PUBLIC_MARKETS_MOCK_SETTLEMENT_SECONDS (default 60; the real cycle is
// T+2) through the same processing as a real settlement webhook.
type MockCustodian struct {
	E    *Engine
	Code string
}

func (m MockCustodian) SendInstruction(ctx context.Context, side string, in partners.CustodianInstruction) (partners.Ack, error) {
	payload, _ := json.Marshal(map[string]string{
		"instructionId": in.InstructionID, "status": "settlement_final", "settledQuantity": in.Quantity,
		"settlementDate":     m.E.now().Add(envSeconds("PUBLIC_MARKETS_MOCK_SETTLEMENT_SECONDS", 60)).In(Lagos).Format("2006-01-02"),
		"custodianReference": "MOCK-CSD-" + strings.ToUpper(randomHex(4)), "side": side, "assetCode": in.AssetCode,
	})
	ev := pm.MockEvent{DueAt: m.E.now().Add(envSeconds("PUBLIC_MARKETS_MOCK_SETTLEMENT_SECONDS", 60)), Kind: EventSettlement,
		Partner: "CUSTODIAN:" + m.Code, Payload: string(payload)}
	if err := m.E.DB.Create(&ev).Error; err != nil {
		return partners.Ack{}, err
	}
	return partners.Ack{Status: "received", StatusCode: 202}, nil
}

func (m MockCustodian) Position(ctx context.Context, assetCode string) (partners.Position, error) {
	units := m.E.mockBook(assetCode)
	return partners.Position{AssetCode: assetCode, UnitsHeld: units, AsOf: m.E.now()}, nil
}

// mockBook is the mock Custodian's units for an asset, started from the
// latest recorded position.
func (e *Engine) mockBook(assetCode string) decimal.Decimal {
	var b pm.MockCustodianBook
	if e.DB.First(&b, "asset_code = ?", assetCode).Error == nil {
		return d(b.Units)
	}
	start := decimal.Zero
	var a pm.Asset
	if e.DB.First(&a, "asset_code = ?", assetCode).Error == nil {
		if p, ok := e.Position(&a); ok {
			start = d(p.RealUnitsHeld)
		}
	}
	e.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&pm.MockCustodianBook{AssetCode: assetCode, Units: start.String(), UpdatedAt: e.now()})
	return start
}

// MockDealingMember stands in for a broker: it fills every order at the
// reference price (within a few basis points) after
// PUBLIC_MARKETS_MOCK_EXECUTION_SECONDS (default 20).
type MockDealingMember struct {
	E    *Engine
	Code string
}

func (m MockDealingMember) PlaceOrder(ctx context.Context, o partners.DealingOrder) (partners.Ack, error) {
	a, err := m.E.AssetByCode(o.AssetCode)
	if err != nil {
		return partners.Ack{}, &partners.PermanentError{StatusCode: 400, Body: "unknown asset"}
	}
	price := d(a.LastPrice)
	if !price.IsPositive() {
		return partners.Ack{}, &partners.PermanentError{StatusCode: 400, Body: "no market price"}
	}
	bps, _ := rand.Int(rand.Reader, big.NewInt(11)) // -5..+5 bps
	price = price.Mul(decimal.NewFromInt(10000 + bps.Int64() - 5)).Div(decimal.NewFromInt(10000)).Round(2)
	at := m.E.now().Add(envSeconds("PUBLIC_MARKETS_MOCK_EXECUTION_SECONDS", 20))
	payload, _ := json.Marshal(map[string]string{"orderId": o.OrderID, "status": "FILLED", "executedQuantity": o.Quantity,
		"executedPrice": price.String(), "executionTime": at.Format(time.RFC3339)})
	ev := pm.MockEvent{DueAt: at, Kind: EventExecution, Partner: "DEALING_MEMBER:" + m.Code, Payload: string(payload)}
	if err := m.E.DB.Create(&ev).Error; err != nil {
		return partners.Ack{}, err
	}
	return partners.Ack{Status: "received", StatusCode: 202}, nil
}

// DeliverMockEvents hands due mock callbacks to the normal webhook
// processing; a settlement also moves the mock Custodian's book.
func (e *Engine) DeliverMockEvents(ctx context.Context) {
	var due []pm.MockEvent
	e.DB.Where("delivered_at IS NULL AND due_at <= ?", e.now()).Order("due_at").Limit(100).Find(&due)
	for _, m := range due {
		res := e.DB.Model(&pm.MockEvent{}).Where("id = ? AND delivered_at IS NULL", m.ID).Update("delivered_at", e.now())
		if res.RowsAffected == 0 {
			continue
		}
		var body map[string]string
		_ = json.Unmarshal([]byte(m.Payload), &body)
		if m.Kind == EventSettlement && body["status"] == "settlement_final" {
			cur := e.mockBook(body["assetCode"])
			q := d(body["settledQuantity"])
			if strings.EqualFold(body["side"], "SELL") {
				q = q.Neg()
			}
			e.DB.Model(&pm.MockCustodianBook{}).Where("asset_code = ?", body["assetCode"]).Updates(map[string]interface{}{"units": cur.Add(q).String(), "updated_at": e.now()})
		}
		ref := body["instructionId"]
		if m.Kind == EventExecution {
			ref = body["orderId"]
		}
		if _, err := e.ReceivePartnerEvent(m.Partner, m.Kind, ref+":"+body["status"], []byte(m.Payload)); err != nil {
			e.logf("mock %s %s: %v", m.Kind, ref, err)
		}
	}
}

// MockPositionFeed sends each mock Custodian's daily position feed (its
// book), as a real Custodian would at end of day.
func (e *Engine) MockPositionFeed(ctx context.Context, custodianCode string) error {
	var c pm.Custodian
	if err := e.DB.First(&c, "code = ?", custodianCode).Error; err != nil {
		return err
	}
	var assets []pm.Asset
	e.DB.Where("custodian_id = ?", c.CustodianID).Find(&assets)
	type pos struct {
		AssetCode string `json:"assetCode"`
		UnitsHeld string `json:"unitsHeld"`
	}
	feed := struct {
		AsOf      string `json:"asOf"`
		Positions []pos  `json:"positions"`
	}{AsOf: e.now().In(Lagos).Format("2006-01-02")}
	for _, a := range assets {
		feed.Positions = append(feed.Positions, pos{AssetCode: a.AssetCode, UnitsHeld: e.mockBook(a.AssetCode).String()})
	}
	body, _ := json.Marshal(feed)
	_, err := e.ReceivePartnerEvent("CUSTODIAN:"+c.Code, EventPositionFeed, fmt.Sprintf("%s:%d", feed.AsOf, e.now().Unix()), body)
	return err
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// MockPriceFeed moves prices in small steps during market hours and holds
// the last close after hours, until a price vendor is contracted (OI-CONS-4).
type MockPriceFeed struct{ E *Engine }

func (m MockPriceFeed) Quote(ctx context.Context, isin, market string) (partners.Quote, error) {
	var a pm.Asset
	if err := m.E.DB.First(&a, "isin = ?", isin).Error; err != nil {
		return partners.Quote{}, errors.New("unknown ISIN")
	}
	last := d(a.LastPrice)
	if !last.IsPositive() {
		return partners.Quote{}, errors.New("no starting price: record one in Trovo Manager")
	}
	sess := MarketSession(a.Market, m.E.now(), LoadSettings(m.E.DB))
	if !sess.Open {
		return partners.Quote{ISIN: isin, Price: last, AsOf: m.E.now(), MarketOpen: false}, nil
	}
	step, _ := rand.Int(rand.Reader, big.NewInt(61)) // -30..+30 bps
	p := last.Mul(decimal.NewFromInt(10000 + step.Int64() - 30)).Div(decimal.NewFromInt(10000))
	if prev := d(a.PreviousClose); prev.IsPositive() { // stay within NGX's ±10% daily band
		p = decimal.Max(decimal.Min(p, prev.Mul(decimal.NewFromFloat(1.1))), prev.Mul(decimal.NewFromFloat(0.9)))
	}
	return partners.Quote{ISIN: isin, Price: p.Round(2), AsOf: m.E.now(), MarketOpen: true}, nil
}

// PriceFeedFromEnv is the REST vendor when PUBLIC_MARKETS_PRICE_FEED=rest
// (PUBLIC_MARKETS_PRICE_FEED_URL, _AUTH, _CREDENTIALS), else the mock.
func PriceFeedFromEnv(e *Engine) partners.PriceFeed {
	if strings.EqualFold(os.Getenv("PUBLIC_MARKETS_PRICE_FEED"), "rest") {
		rc, err := partners.NewRESTClient(os.Getenv("PUBLIC_MARKETS_PRICE_FEED_URL"), os.Getenv("PUBLIC_MARKETS_PRICE_FEED_AUTH"), os.Getenv("PUBLIC_MARKETS_PRICE_FEED_CREDENTIALS"))
		if err == nil {
			return partners.RESTPriceFeed{RESTClient: rc}
		}
		e.logf("price feed: %v; using the mock feed", err)
	}
	return MockPriceFeed{E: e}
}
