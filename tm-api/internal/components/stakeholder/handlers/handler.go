package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	portalModels "admin-panel-dashboard/internal/components/stakeholder/models"
	"admin-panel-dashboard/internal/components/stakeholder/services"
	"admin-panel-dashboard/internal/middleware"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db                    *gorm.DB
	Assets                *services.AssetService
	Profile               *services.ProfileService
	Authorizations        *services.AuthorizationService
	Notifications         *services.NotificationService
	FundReleases          *services.FundReleaseService
	DueDiligence          *services.DueDiligenceService
	Financial             *services.FinancialService
	Custodian             *services.CustodianService
	Documents             *services.DocumentService
	Operations            *services.OperationsService
	Dashboards            *services.DashboardService
	Structuring           *services.StructuringService
	FundManagement        *services.FundManagementService
	ComplianceRequirement *services.ComplianceRequirementService
}

func NewHandler(
	db *gorm.DB,
	assets *services.AssetService,
	profile *services.ProfileService,
	authorizations *services.AuthorizationService,
	notifications *services.NotificationService,
	fundReleases *services.FundReleaseService,
	dueDiligence *services.DueDiligenceService,
	financial *services.FinancialService,
	custodian *services.CustodianService,
	documents *services.DocumentService,
	operations *services.OperationsService,
	dashboards *services.DashboardService,
	structuring *services.StructuringService,
	fundManagement *services.FundManagementService,
	complianceRequirement *services.ComplianceRequirementService,
) *Handler {
	return &Handler{
		db:                    db,
		Assets:                assets,
		Profile:               profile,
		Authorizations:        authorizations,
		Notifications:         notifications,
		FundReleases:          fundReleases,
		DueDiligence:          dueDiligence,
		Financial:             financial,
		Custodian:             custodian,
		Documents:             documents,
		Operations:            operations,
		Dashboards:            dashboards,
		Structuring:           structuring,
		FundManagement:        fundManagement,
		ComplianceRequirement: complianceRequirement,
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "stakeholder_portal"})
}

func (h *Handler) GetTrusteeFundManagementSummary(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	summary, err := h.FundManagement.TrusteeSummary(c.Request.Context(), auth)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": summary})
}

func (h *Handler) CreateAssetAssignment(c *gin.Context) {
	var req portalModels.CreateAssetAssignmentRequest
	if !bind(c, &req) {
		return
	}
	assignment, err := h.Assets.CreateAssignment(c.Request.Context(), services.AuthContext{}, "trovo_admin", req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assignment})
}

func (h *Handler) CreateCompliance(c *gin.Context) {
	var req portalModels.CreateComplianceItemRequest
	if !bind(c, &req) {
		return
	}
	actor := services.AuthContext{
		MemberID:      c.GetString(middleware.AdminContextUserID),
		DashboardRole: c.GetString(middleware.AdminContextRole),
	}
	item, err := h.Custodian.CreateCompliance(c.Request.Context(), actor, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *Handler) CreateComplianceRequirement(c *gin.Context) {
	var req portalModels.CreateComplianceRequirementRequest
	if !bind(c, &req) {
		return
	}
	item, err := h.ComplianceRequirement.CreateAdHocRequirement(c.Request.Context(), adminActorFromContext(c), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *Handler) ListAdminCompliance(c *gin.Context) {
	p := services.ParsePagination(c)
	filters := services.AdminComplianceListFilters{
		CustodianOrgID: c.Query("custodian_org_id"),
		AssetID:        c.Query("asset_id"),
		Status:         c.Query("status"),
	}
	records, total, err := h.Custodian.ListComplianceForAdmin(c.Request.Context(), p.Page, p.Limit, filters)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{
		Records: records,
		Meta:    services.PaginationMeta(p.Page, p.Limit, total),
	}})
}

func adminActorFromContext(c *gin.Context) services.AuthContext {
	return services.AuthContext{
		MemberID:      c.GetString(middleware.AdminContextUserID),
		DashboardRole: c.GetString(middleware.AdminContextRole),
	}
}

func (h *Handler) ListComplianceTemplates(c *gin.Context) {
	p := services.ParsePagination(c)
	orgType := c.Query("org_type")
	var level *int
	if raw := c.Query("level"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, coreModels.ErrorResponse{Error: "level must be a positive integer"})
			return
		}
		level = &parsed
	}
	records, total, err := h.ComplianceRequirement.ListTemplates(c.Request.Context(), p.Page, p.Limit, orgType, level)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) CreateComplianceTemplate(c *gin.Context) {
	var req portalModels.CreateComplianceTemplateRequest
	if !bind(c, &req) {
		return
	}
	template, err := h.ComplianceRequirement.CreateTemplate(c.Request.Context(), adminActorFromContext(c), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": template})
}

func (h *Handler) UpdateComplianceTemplate(c *gin.Context) {
	var req portalModels.UpdateComplianceTemplateRequest
	if !bind(c, &req) {
		return
	}
	template, err := h.ComplianceRequirement.UpdateTemplate(c.Request.Context(), adminActorFromContext(c), c.Param("id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": template})
}

func (h *Handler) ListComplianceRequirements(c *gin.Context) {
	p := services.ParsePagination(c)
	var level *int
	if raw := c.Query("level"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, coreModels.ErrorResponse{Error: "level must be a positive integer"})
			return
		}
		level = &parsed
	}
	records, total, err := h.ComplianceRequirement.ListRequirements(c.Request.Context(), p.Page, p.Limit, c.Query("org_id"), c.Query("category"), level, c.Query("status"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) GetComplianceRequirement(c *gin.Context) {
	instance, err := h.ComplianceRequirement.GetRequirementDetail(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": instance})
}

func (h *Handler) StartComplianceRequirementReview(c *gin.Context) {
	instance, err := h.ComplianceRequirement.StartReview(c.Request.Context(), adminActorFromContext(c), c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": instance})
}

func (h *Handler) ReviewComplianceRequirement(c *gin.Context) {
	var req portalModels.ReviewComplianceRequirementRequest
	if !bind(c, &req) {
		return
	}
	instance, err := h.ComplianceRequirement.ReviewRequirement(c.Request.Context(), adminActorFromContext(c), c.Param("id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": instance})
}

func (h *Handler) ListOrgCompliance(c *gin.Context) {
	auth, ok := organizationAuthFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	records, total, err := h.ComplianceRequirement.ListOrgRequirements(c.Request.Context(), auth, p.Page, p.Limit, c.Query("status"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) SubmitOrgCompliance(c *gin.Context) {
	auth, ok := organizationAuthFromContext(c)
	if !ok {
		return
	}
	var req portalModels.SubmitComplianceRequirementRequest
	if !bind(c, &req) {
		return
	}
	instance, err := h.ComplianceRequirement.SubmitOrgRequirement(c.Request.Context(), auth, c.Param("item_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": instance})
}

func (h *Handler) UpdateOrgCompliance(c *gin.Context) {
	auth, ok := organizationAuthFromContext(c)
	if !ok {
		return
	}
	var req portalModels.UpdateOrgComplianceRequirementRequest
	if !bind(c, &req) {
		return
	}
	instance, err := h.ComplianceRequirement.UpdateOrgRequirementStatus(c.Request.Context(), auth, c.Param("item_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": instance})
}

func (h *Handler) GetProfile(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	profile, err := h.Profile.GetProfile(c.Request.Context(), auth)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": profile})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.UpdateProfileRequest
	if !bind(c, &req) {
		return
	}
	profile, err := h.Profile.UpdateProfile(c.Request.Context(), auth, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": profile})
}

func (h *Handler) UpdateNotificationPreferences(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.UpdateNotificationPreferencesRequest
	if !bind(c, &req) {
		return
	}
	preferences, err := h.Profile.UpdateNotificationPreferences(c.Request.Context(), auth, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": preferences})
}

func (h *Handler) ListAssets(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	resp, err := h.Assets.ListAssets(c.Request.Context(), auth, portalModels.AssetListFilters{
		Page:      p.Page,
		Limit:     p.Limit,
		Status:    c.Query("status"),
		AssetType: c.Query("asset_type"),
		Sector:    c.Query("asset_sector"),
		Search:    c.Query("search"),
		DateFrom:  c.Query("date_from"),
		DateTo:    c.Query("date_to"),
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *Handler) GetAsset(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	asset, err := h.Assets.GetAssetDetails(c.Request.Context(), auth, c.Param("asset_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": asset})
}

func (h *Handler) ListAssetActivities(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	activities, total, err := h.Assets.ListAssetActivities(c.Request.Context(), auth, c.Param("asset_id"), p.Page, p.Limit)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{
		Records: activities,
		Meta:    services.PaginationMeta(p.Page, p.Limit, total),
	}})
}

func (h *Handler) CreateAuthorization(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.CreateAuthorizationRequest
	if !bind(c, &req) {
		return
	}
	challenge, err := h.Authorizations.CreateChallenge(c.Request.Context(), auth, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": challenge})
}

func (h *Handler) VerifyAuthorization(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.VerifyAuthorizationRequest
	_ = c.ShouldBindJSON(&req)
	challenge, err := h.Authorizations.VerifyChallenge(c.Request.Context(), auth, c.Param("challenge_id"), req.AuthID)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": challenge})
}

func (h *Handler) ListNotifications(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	records, total, err := h.Notifications.List(c.Request.Context(), auth, p.Page, p.Limit)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) MarkNotificationRead(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	notification, err := h.Notifications.MarkRead(c.Request.Context(), auth, c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": notification})
}

func (h *Handler) ListAuditTrail(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	query := h.db.WithContext(c.Request.Context()).Model(&portalModels.StakeholderAuditLog{}).Where("actor_org_id = ?", auth.OrganizationID)
	if action := c.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}
	if role := c.Query("role"); role != "" {
		query = query.Where("actor_role = ?", role)
	}
	if entityType := c.Query("entity_type"); entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		writeError(c, err)
		return
	}
	var logs []portalModels.StakeholderAuditLog
	if err := query.Order("created_at desc").Offset((p.Page - 1) * p.Limit).Limit(p.Limit).Find(&logs).Error; err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: logs, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) CreateDocument(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.CreateDocumentRequest
	if !bind(c, &req) {
		return
	}
	doc, err := h.Documents.Create(c.Request.Context(), auth, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": doc})
}

func (h *Handler) UploadDocument(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	h.uploadDocument(c, auth, "")
}

func (h *Handler) UploadOrganizationComplianceDocument(c *gin.Context) {
	auth, ok := organizationAuthFromContext(c)
	if !ok {
		return
	}
	h.uploadDocument(c, auth, portalModels.DocumentCategoryComplianceRequirement)
}

func (h *Handler) uploadDocument(c *gin.Context, auth services.AuthContext, forcedCategory string) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, (10<<20)+(1<<20))
	file, header, err := c.Request.FormFile("document_file")
	if errors.Is(err, http.ErrMissingFile) {
		file, header, err = c.Request.FormFile("documentFile")
	}
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "document exceeds 10MB"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "document_file is required"})
		return
	}
	defer file.Close()
	content, err := io.ReadAll(io.LimitReader(file, (10<<20)+1))
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "document exceeds 10MB"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read document"})
		return
	}
	roles := make([]string, 0)
	for _, value := range c.PostFormArray("access_roles") {
		roles = append(roles, strings.Split(value, ",")...)
	}
	category := c.PostForm("category")
	assetID := c.PostForm("asset_id")
	if forcedCategory != "" {
		category = forcedCategory
		assetID = ""
		roles = nil
	}
	doc, err := h.Documents.Upload(c.Request.Context(), auth, portalModels.UploadDocumentRequest{
		AssetID: assetID, Category: category,
		Title: c.PostForm("title"), AccessRoles: roles,
		Filename: header.Filename, Content: content,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": doc})
}

func (h *Handler) DocumentCategories(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": h.Documents.Categories()})
}

func (h *Handler) ListDocuments(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	docs, total, err := h.Documents.List(c.Request.Context(), auth, p.Page, p.Limit, c.Query("asset_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: docs, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) GetDocument(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	doc, err := h.Documents.Get(c.Request.Context(), auth, c.Param("doc_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": doc})
}

func (h *Handler) DownloadDocument(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	download, err := h.Documents.Download(c.Request.Context(), auth, c.Param("doc_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	writeDocumentDownload(c, download)
}

func (h *Handler) DownloadOrganizationComplianceDocument(c *gin.Context) {
	auth, ok := organizationAuthFromContext(c)
	if !ok {
		return
	}
	download, err := h.Documents.DownloadOwnedByOrganization(c.Request.Context(), auth.OrganizationID, c.Param("doc_id"), portalModels.DocumentCategoryComplianceRequirement)
	if err != nil {
		writeError(c, err)
		return
	}
	writeDocumentDownload(c, download)
}

func (h *Handler) DownloadComplianceRequirementDocumentForAdmin(c *gin.Context) {
	download, err := h.ComplianceRequirement.DownloadSubmissionDocument(c.Request.Context(), c.Param("id"), c.Param("doc_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	writeDocumentDownload(c, download)
}

func writeDocumentDownload(c *gin.Context, download *services.DocumentDownload) {
	defer download.Body.Close()
	filename := strings.NewReplacer(`"`, "", "\r", "", "\n", "").Replace(download.OriginalFilename)
	if filename == "" {
		filename = "document"
	}
	mimeType := download.MimeType
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	c.DataFromReader(http.StatusOK, download.SizeBytes, mimeType, download.Body, map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, filename),
		"Cache-Control":       "private, no-store",
	})
}

func (h *Handler) Dashboard(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	dashboard, err := h.Dashboards.Dashboard(c.Request.Context(), auth)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": dashboard})
}

func (h *Handler) GetDueDiligence(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	checklist, err := h.DueDiligence.GetChecklist(c.Request.Context(), auth, c.Param("asset_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": checklist})
}

func (h *Handler) UpdateDueDiligenceItem(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.UpdateDueDiligenceItemRequest
	if !bind(c, &req) {
		return
	}
	item, err := h.DueDiligence.UpdateItem(c.Request.Context(), auth, c.Param("asset_id"), c.Param("item_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *Handler) UpdateDueDiligenceCategory(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.UpdateDueDiligenceCategoryRequest
	if !bind(c, &req) {
		return
	}
	category, err := h.DueDiligence.UpdateCategory(c.Request.Context(), auth, c.Param("asset_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (h *Handler) ApproveDueDiligence(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	checklist, err := h.DueDiligence.Approve(c.Request.Context(), auth, c.Param("asset_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": checklist})
}

func (h *Handler) RejectDueDiligence(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.RejectRequest
	if !bind(c, &req) {
		return
	}
	checklist, err := h.DueDiligence.Reject(c.Request.Context(), auth, c.Param("asset_id"), req.Reason)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": checklist})
}

func (h *Handler) CreateFundRelease(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.CreateFundReleaseRequest
	if !bind(c, &req) {
		return
	}
	record, err := h.FundReleases.Create(c.Request.Context(), auth, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": record})
}

func (h *Handler) ListFundReleases(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	records, total, err := h.FundReleases.List(c.Request.Context(), auth, p.Page, p.Limit, c.Query("status"), c.Query("asset_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) GetFundRelease(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	record, err := h.FundReleases.GetDetails(c.Request.Context(), auth, c.Param("request_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": record})
}

func (h *Handler) GetAssetManagerFundSummary(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	summary, err := h.FundReleases.AssetManagerSummary(c.Request.Context(), auth, c.Query("asset_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": summary})
}

func (h *Handler) ApproveFundRelease(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.ChallengeActionRequest
	if !bind(c, &req) {
		return
	}
	record, err := h.FundReleases.Approve(c.Request.Context(), auth, c.Param("request_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": record})
}

func (h *Handler) RejectFundRelease(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.RejectRequest
	if !bind(c, &req) {
		return
	}
	record, err := h.FundReleases.Reject(c.Request.Context(), auth, c.Param("request_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": record})
}

func (h *Handler) ExecuteFundRelease(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.ExecuteFundReleaseRequest
	if !bind(c, &req) {
		return
	}
	record, err := h.FundReleases.Execute(c.Request.Context(), auth, c.Param("request_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": record})
}

func (h *Handler) UpdateFundReleaseStatus(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.UpdateFundReleaseStatusRequest
	if !bind(c, &req) {
		return
	}
	record, err := h.FundReleases.UpdateExecutionStatus(c.Request.Context(), auth, c.Param("request_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": record})
}

func (h *Handler) ListRevenue(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	records, total, err := h.Financial.ListRevenue(c.Request.Context(), auth, p.Page, p.Limit)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) CreateRevenue(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.CreateRevenueRequest
	if !bind(c, &req) {
		return
	}
	record, err := h.Financial.CreateRevenue(c.Request.Context(), auth, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": record})
}

func (h *Handler) SubmitDistribution(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.SubmitDistributionRequest
	if !bind(c, &req) {
		return
	}
	distribution, err := h.Financial.SubmitDistribution(c.Request.Context(), auth, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": distribution})
}

func (h *Handler) ListDistributions(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	records, total, err := h.Financial.ListDistributions(c.Request.Context(), auth, p.Page, p.Limit, false)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) ListDistributionHistory(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	records, total, err := h.Financial.ListDistributions(c.Request.Context(), auth, p.Page, p.Limit, true)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) GetDistribution(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	detail, err := h.Financial.GetDistribution(c.Request.Context(), auth, c.Param("dist_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": detail})
}

func (h *Handler) DownloadDistributionPayouts(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	distributionID := strings.TrimSpace(c.Param("dist_id"))
	content, err := h.Financial.DistributionPayoutCSV(c.Request.Context(), auth, distributionID)
	if err != nil {
		writeError(c, err)
		return
	}
	safeID := strings.NewReplacer(`"`, "", "\r", "", "\n", "").Replace(distributionID)
	if safeID == "" {
		safeID = "distribution"
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="distribution-payouts-%s.csv"`, safeID))
	c.Header("Cache-Control", "private, no-store")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", content)
}

func (h *Handler) AuthorizeDistribution(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.ChallengeActionRequest
	if !bind(c, &req) {
		return
	}
	distribution, err := h.Financial.AuthorizeDistribution(c.Request.Context(), auth, c.Param("dist_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": distribution})
}

func (h *Handler) RejectDistribution(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.RejectRequest
	if !bind(c, &req) {
		return
	}
	distribution, err := h.Financial.RejectDistribution(c.Request.Context(), auth, c.Param("dist_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": distribution})
}

func (h *Handler) ListValuations(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	records, total, err := h.Financial.ListValuations(c.Request.Context(), auth, c.Param("asset_id"), p.Page, p.Limit)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) CreateValuation(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.CreateValuationRequest
	if !bind(c, &req) {
		return
	}
	valuation, err := h.Financial.CreateValuation(c.Request.Context(), auth, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": valuation})
}

func (h *Handler) RequestIndependentValuation(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	valuation, err := h.Financial.RequestIndependentValuation(c.Request.Context(), auth, c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": valuation})
}

func (h *Handler) ListAccounts(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	records, total, err := h.Custodian.ListAccounts(c.Request.Context(), auth, p.Page, p.Limit)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) CreateAccount(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.CreateSegregatedAccountRequest
	if !bind(c, &req) {
		return
	}
	account, err := h.Custodian.CreateAccount(c.Request.Context(), auth, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": account})
}

func (h *Handler) GetAccount(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	account, err := h.Custodian.GetAccount(c.Request.Context(), auth, c.Param("account_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": account})
}

func (h *Handler) UpdateAccount(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.UpdateSegregatedAccountRequest
	if !bind(c, &req) {
		return
	}
	account, err := h.Custodian.UpdateAccount(c.Request.Context(), auth, c.Param("account_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": account})
}

func (h *Handler) ListCompliance(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	records, total, err := h.Custodian.ListCompliance(c.Request.Context(), auth, p.Page, p.Limit)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) UpdateCompliance(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.UpdateComplianceItemRequest
	if !bind(c, &req) {
		return
	}
	item, err := h.Custodian.UpdateCompliance(c.Request.Context(), auth, c.Param("item_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *Handler) UpdateAssetOperation(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.UpdateAssetOperationRequest
	if !bind(c, &req) {
		return
	}
	operation, err := h.Operations.UpdateAssetOperation(c.Request.Context(), auth, c.Param("asset_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": operation})
}

func (h *Handler) CreateReport(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.CreateReportRequest
	if !bind(c, &req) {
		return
	}
	report, err := h.Operations.CreateReport(c.Request.Context(), auth, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": report})
}

func (h *Handler) ListReports(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	p := services.ParsePagination(c)
	records, total, err := h.Operations.ListReports(c.Request.Context(), auth, p.Page, p.Limit)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": portalModels.PaginatedResponse{Records: records, Meta: services.PaginationMeta(p.Page, p.Limit, total)}})
}

func (h *Handler) ListReportTypes(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	reportTypes, err := h.Operations.ReportTypes(auth)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reportTypes})
}

func (h *Handler) SubmitReport(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	report, err := h.Operations.SubmitReport(c.Request.Context(), auth, c.Param("report_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": report})
}

func (h *Handler) ConfirmStructuring(c *gin.Context) {
	auth, ok := authFromContext(c)
	if !ok {
		return
	}
	var req portalModels.ConfirmStructuringRequest
	if !bind(c, &req) {
		return
	}
	status, err := h.Structuring.ConfirmComplete(c.Request.Context(), auth, c.Param("asset_id"), req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": status})
}

func bind(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, coreModels.ErrorResponse{Error: err.Error()})
		return false
	}
	return true
}

func authFromContext(c *gin.Context) (services.AuthContext, bool) {
	auth, err := services.AuthContextFromGin(c)
	if err != nil {
		writeError(c, err)
		return services.AuthContext{}, false
	}
	return auth, true
}

func organizationAuthFromContext(c *gin.Context) (services.AuthContext, bool) {
	auth, err := services.OrganizationAuthContextFromGin(c)
	if err != nil {
		writeError(c, err)
		return services.AuthContext{}, false
	}
	return auth, true
}

func writeError(c *gin.Context, err error) {
	c.JSON(services.ErrorStatus(err), coreModels.ErrorResponse{Error: services.ErrorMessage(err)})
}
