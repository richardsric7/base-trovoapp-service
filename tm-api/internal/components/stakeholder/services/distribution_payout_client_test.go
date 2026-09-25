package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"admin-panel-dashboard/internal/components/stakeholder/models"
	"admin-panel-dashboard/internal/trovosdk"
)

func TestWalletDistributionPayoutClientContract(t *testing.T) {
	distributionID := "b1ec38f7-6bfd-46fb-b693-5dad05a905be"
	var registrations int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-TW-SERVICE-LINK-API-KEY") != "service-key" {
			t.Errorf("missing Wallet service-link API key")
		}
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/trovo-api/stakeholder-distributions":
			var input DistributionPayoutRegistration
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				t.Errorf("decode registration: %v", err)
			}
			if input.DistributionID != distributionID || input.Amount != "1250" || input.Currency != "CNGN" {
				t.Errorf("registration = %+v", input)
			}
			registrations++
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": models.DistributionPayoutResponse{
				DistributionID: distributionID, TokenizedAssetID: "asset-1", Amount: "1250",
				Currency: "CNGN", Status: "registered", Payouts: []models.DistributionPayoutRecordResponse{},
			}})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/trovo-api/stakeholder-distributions/"+distributionID:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": models.DistributionPayoutResponse{
				DistributionID: distributionID, TokenizedAssetID: "asset-1", Amount: "1250",
				Currency: "CNGN", Status: "scheduled", Payouts: []models.DistributionPayoutRecordResponse{{ID: "payout-1", Paid: true}},
			}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewWalletDistributionPayoutClient(&trovosdk.ServiceLink{ApiBaseUrl: server.URL, ApiKey: "service-key"})
	registered, err := client.Register(context.Background(), DistributionPayoutRegistration{
		DistributionID: distributionID, TokenizedAssetID: "asset-1", Amount: "1250", Currency: "CNGN",
	})
	if err != nil || registered.Status != "registered" || registrations != 1 {
		t.Fatalf("register response=%+v registrations=%d err=%v", registered, registrations, err)
	}
	payout, err := client.Get(context.Background(), distributionID)
	if err != nil || payout.Status != "scheduled" || len(payout.Payouts) != 1 || !payout.Payouts[0].Paid {
		t.Fatalf("get payout=%+v err=%v", payout, err)
	}
}

func TestWalletDistributionPayoutClientTreatsMissingLegacyPayoutAsEmpty(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	client := NewWalletDistributionPayoutClient(&trovosdk.ServiceLink{ApiBaseUrl: server.URL, ApiKey: "service-key"})
	payout, err := client.Get(context.Background(), "b1ec38f7-6bfd-46fb-b693-5dad05a905be")
	if err != nil || payout != nil {
		t.Fatalf("missing payout=%+v err=%v", payout, err)
	}
}
