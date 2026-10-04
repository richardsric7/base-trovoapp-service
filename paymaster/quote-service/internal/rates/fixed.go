package rates

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

// fixed: {"type": "fixed", "value": "1"}
//
// A constant rate - a stablecoin's peg to its currency, or a manual
// override while a market source is unavailable.
type fixedSource struct{ v decimal.Decimal }

func (f fixedSource) Fetch(context.Context) (decimal.Decimal, error) { return f.v, nil }

func init() {
	Register("fixed", func(def json.RawMessage, _ Env) (Source, error) {
		var c struct {
			Value string `json:"value"`
		}
		if err := json.Unmarshal(def, &c); err != nil {
			return nil, err
		}
		v, err := decimal.NewFromString(c.Value)
		if err != nil || !v.IsPositive() {
			return nil, fmt.Errorf(`"value" must be a positive decimal`)
		}
		return fixedSource{v}, nil
	})
}
