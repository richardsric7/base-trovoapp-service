package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"
	"admin-panel-dashboard/internal/trovosdk"

	"github.com/google/uuid"
)

type fakeAdminDocumentStorage struct {
	objects map[string][]byte
	deleted string
}

func (f *fakeAdminDocumentStorage) Upload(_ context.Context, filename, mimeType string, content []byte) (*StoredDocument, error) {
	id := uuid.NewString()
	if f.objects == nil {
		f.objects = make(map[string][]byte)
	}
	f.objects[id] = append([]byte(nil), content...)
	return &StoredDocument{ObjectID: id, OriginalFilename: filename, MimeType: mimeType, SizeBytes: int64(len(content)), SHA256: strings.Repeat("a", 64)}, nil
}

func (f *fakeAdminDocumentStorage) Download(_ context.Context, objectID string) (*DocumentDownload, error) {
	content := f.objects[objectID]
	return &DocumentDownload{Body: io.NopCloser(bytes.NewReader(content)), OriginalFilename: "report.pdf", MimeType: "application/pdf", SizeBytes: int64(len(content))}, nil
}

func (f *fakeAdminDocumentStorage) Delete(_ context.Context, objectID string) error {
	f.deleted = objectID
	delete(f.objects, objectID)
	return nil
}

func TestManagedDocumentsDriveFundReleaseAndValuation(t *testing.T) {
	db := testDB(t, &coreModels.Organization{}, &models.StakeholderDocument{}, &models.StakeholderAssetAssignment{}, &models.FundReleaseRequest{}, &models.AssetValuation{})
	managerID := uint64(77)
	managerOrg, trusteeOrg := "manager-org", "trustee-org"
	asset := stakeholderDB.TokenizedAsset{ID: "asset-1", AssetCode: "AST1", AssetManagerID: managerID}
	assets := NewAssetService(db, fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{asset}}, nil)
	stakeholderType := models.StakeholderTypeAssetManager
	if err := db.Create(&coreModels.Organization{
		ID: managerOrg, Name: "Manager", Email: "manager@example.com", Type: coreModels.OrganizationTypeAssetManager,
		Status: coreModels.OrganizationStatusActive, CreatedBy: "test", StakeholderID: &managerID, StakeholderType: &stakeholderType,
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.StakeholderAssetAssignment{
		ID: uuid.NewString(), AssetID: asset.ID, AssetCode: asset.AssetCode,
		AssetManagerOrgID: &managerOrg, AssetManagerStakeholderID: &managerID,
		TrusteeOrgID: &trusteeOrg, Status: models.AssignmentStatusActive,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}
	storage := &fakeAdminDocumentStorage{}
	documents := NewDocumentService(db, assets, nil, storage)
	auth := AuthContext{OrganizationID: managerOrg, MemberID: "manager-member", StakeholderID: managerID, DashboardRole: models.DashboardRoleAssetManager}
	upload := func(category string) *models.StakeholderDocument {
		doc, err := documents.Upload(context.Background(), auth, models.UploadDocumentRequest{
			AssetID: asset.ID, Category: category, Title: "Evidence", Filename: "evidence.pdf", Content: []byte("%PDF-1.4\nprivate"),
		})
		if err != nil {
			t.Fatalf("upload %s: %v", category, err)
		}
		if doc.ID == "" || doc.StorageObjectID == "" || doc.FileURL != "" || doc.UploadedByOrgID != managerOrg || doc.AssetID != asset.ID {
			t.Fatalf("invalid document metadata: %+v", doc)
		}
		return doc
	}
	fundDoc := upload(models.DocumentCategoryFundReleaseSupporting)
	valuationDoc := upload(models.DocumentCategoryValuationReport)

	fundService := NewFundReleaseService(db, assets, nil, nil, nil, documents)
	fund, err := fundService.Create(context.Background(), auth, models.CreateFundReleaseRequest{
		AssetID: asset.ID, Amount: "2500", Currency: "CNGN", Purpose: "Milestone", SupportingDocumentIDs: []string{fundDoc.ID},
	})
	if err != nil || len(fund.SupportingDocumentIDs) != 1 || fund.SupportingDocumentIDs[0] != fundDoc.ID {
		t.Fatalf("create fund release: fund=%+v err=%v", fund, err)
	}

	financial := NewFinancialService(db, assets, nil, nil, nil, documents)
	if _, err := financial.CreateValuation(context.Background(), auth, models.CreateValuationRequest{
		AssetID: asset.ID, Valuation: "500000", Currency: "CNGN", Methodology: "DCF",
		ValuationDate: "2026-08-20", ReportDocumentID: fundDoc.ID,
	}); ErrorStatus(err) != http.StatusUnprocessableEntity {
		t.Fatalf("wrong-category valuation status=%d err=%v", ErrorStatus(err), err)
	}
	valuation, err := financial.CreateValuation(context.Background(), auth, models.CreateValuationRequest{
		AssetID: asset.ID, Valuation: "500000", Currency: "CNGN", Methodology: "DCF",
		ValuationDate: "2026-08-20", ReportDocumentID: valuationDoc.ID,
	})
	if err != nil || valuation.ReportDocumentID == nil || *valuation.ReportDocumentID != valuationDoc.ID {
		t.Fatalf("create valuation: valuation=%+v err=%v", valuation, err)
	}

	err = documents.ValidateForUse(context.Background(), AuthContext{
		OrganizationID: "other-org", MemberID: "other-member", StakeholderID: managerID, DashboardRole: models.DashboardRoleAssetManager,
	}, asset.ID, models.DocumentCategoryFundReleaseSupporting, []string{fundDoc.ID})
	if ErrorStatus(err) != http.StatusUnprocessableEntity {
		t.Fatalf("cross-org document status=%d err=%v", ErrorStatus(err), err)
	}
	if _, err := fundService.Create(context.Background(), auth, models.CreateFundReleaseRequest{
		AssetID: asset.ID, Amount: "1", Currency: "CNGN", Purpose: "missing docs",
	}); ErrorStatus(err) != http.StatusUnprocessableEntity {
		t.Fatalf("missing document status=%d err=%v", ErrorStatus(err), err)
	}

	download, err := documents.Download(context.Background(), auth, valuationDoc.ID)
	if err != nil {
		t.Fatal(err)
	}
	defer download.Body.Close()
	content, _ := io.ReadAll(download.Body)
	if string(content) != "%PDF-1.4\nprivate" {
		t.Fatalf("download content=%q", content)
	}
}

func TestWalletDocumentStorageClientContract(t *testing.T) {
	objectID := uuid.NewString()
	var uploaded bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-TW-SERVICE-LINK-API-KEY") != "service-key" {
			t.Errorf("missing service link key")
		}
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/trovo-api/stakeholder-documents":
			file, header, err := r.FormFile("document_file")
			if err != nil {
				t.Errorf("multipart file: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer file.Close()
			content, _ := io.ReadAll(file)
			uploaded = header.Filename == "report.pdf" && string(content) == "%PDF-1.4\ncontract"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": StoredDocument{ObjectID: objectID, OriginalFilename: header.Filename, MimeType: "application/pdf", SizeBytes: int64(len(content)), SHA256: strings.Repeat("b", 64)}})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/trovo-api/stakeholder-documents/"+objectID:
			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Disposition", `attachment; filename="report.pdf"`)
			_, _ = w.Write([]byte("download"))
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/trovo-api/stakeholder-documents/"+objectID:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewWalletDocumentStorageClient(&trovosdk.ServiceLink{ApiBaseUrl: server.URL, ApiKey: "service-key"})
	stored, err := client.Upload(context.Background(), "report.pdf", "application/pdf", []byte("%PDF-1.4\ncontract"))
	if err != nil || !uploaded || stored.ObjectID != objectID {
		t.Fatalf("upload stored=%+v uploaded=%v err=%v", stored, uploaded, err)
	}
	download, err := client.Download(context.Background(), objectID)
	if err != nil {
		t.Fatal(err)
	}
	content, _ := io.ReadAll(download.Body)
	_ = download.Body.Close()
	if string(content) != "download" || download.OriginalFilename != "report.pdf" || download.MimeType != "application/pdf" {
		t.Fatalf("download=%+v content=%q", download, content)
	}
	if err := client.Delete(context.Background(), objectID); err != nil {
		t.Fatal(err)
	}
}

func TestDocumentUploadDeletesStoredObjectWhenMetadataWriteFails(t *testing.T) {
	db := testDB(t)
	storage := &fakeAdminDocumentStorage{}
	service := NewDocumentService(db, nil, nil, storage)
	_, err := service.Upload(context.Background(), AuthContext{OrganizationID: "org", MemberID: "member"}, models.UploadDocumentRequest{
		Category: models.DocumentCategoryGeneral, Title: "General", Filename: "general.pdf", Content: []byte("%PDF-1.4\ncleanup"),
	})
	if err == nil || storage.deleted == "" || len(storage.objects) != 0 {
		t.Fatalf("err=%v deleted=%q objects=%d", err, storage.deleted, len(storage.objects))
	}
}
