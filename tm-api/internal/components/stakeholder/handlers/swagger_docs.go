package handlers

import models "admin-panel-dashboard/internal/components/stakeholder/models"

var _ = models.CreateAssetAssignmentRequest{}

// stakeholderHealthDocs documents the stakeholder health route.
// @Summary Stakeholder portal health
// @Tags Stakeholder - Shared
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /stakeholder/shared/health [get]
func stakeholderHealthDocs() {}

// stakeholderProfileDocs documents shared profile routes.
// @Summary Get stakeholder profile
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.ProfileDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/profile [get]
func stakeholderProfileDocs() {}

// stakeholderProfileUpdateDocs documents profile update.
// @Summary Update stakeholder profile
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param payload body models.UpdateProfileRequest true "Profile update"
// @Success 200 {object} models.ProfileDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/profile [put]
func stakeholderProfileUpdateDocs() {}

// stakeholderProfileNotificationDocs documents notification preference update.
// @Summary Update stakeholder notification preferences
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param payload body models.UpdateNotificationPreferencesRequest true "Notification preference update"
// @Success 200 {object} models.NotificationPreferenceDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/profile/notifications [put]
func stakeholderProfileNotificationDocs() {}

// stakeholderAssetListDocs documents shared asset list.
// @Summary List visible stakeholder assets with compliance, custody and statistics
// @Description Returns role-scoped assets. Each record includes token-holder, compliance and custody status; summary contains total, active and liquidated counts for the filtered result.
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param status query string false "Tokenization status"
// @Param asset_type query string false "Asset type"
// @Param asset_sector query string false "Asset sector"
// @Param search query string false "Asset name/code search"
// @Success 200 {object} models.AssetPageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/assets [get]
// @Router /stakeholder/asset-manager/assets [get]
func stakeholderAssetListDocs() {}

// stakeholderAssetDetailDocs documents asset detail.
// @Summary Get enriched visible stakeholder asset detail
// @Description Returns tokenizer and stakeholder assignments plus WalletDB-owned profile, location, ownership, risk/compliance, protection, offering/token sale, fee, receiving-account, application-payment, tokenization-document and lifecycle-progression data. risk_assessment_score and risk_level are null until a canonical assessment result is persisted; maximum_risk_score describes the 100-point UI scale. AdminDB stakeholder documents remain in documents; WalletDB tokenization documents are returned separately in tokenization_documents.
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Produce json
// @Param asset_id path string true "Asset ID or code"
// @Success 200 {object} models.AssetDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/assets/{asset_id} [get]
func stakeholderAssetDetailDocs() {}

// stakeholderAssetActivityDocs documents asset-specific activity history.
// @Summary List activity history for a visible stakeholder asset
// @Description Returns portal audit events linked to the specified asset across its assignments, documents, due diligence, fund releases, distributions, valuations, accounts, compliance, reports, and operational updates. Each record includes actor_name and actor_organization when the actor is an organization member. Results are newest first.
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Produce json
// @Param asset_id path string true "Asset ID or code"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Records per page" default(20) maximum(100)
// @Success 200 {object} models.AuditLogPageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/assets/{asset_id}/activities [get]
func stakeholderAssetActivityDocs() {}

// stakeholderAuthorizationDocs documents authorization challenge creation.
// @Summary Create stakeholder step-up authorization challenge
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param payload body models.CreateAuthorizationRequest true "Authorization challenge"
// @Success 201 {object} models.AuthorizationChallengeDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/authorizations [post]
func stakeholderAuthorizationDocs() {}

// stakeholderAuthorizationVerifyDocs documents authorization challenge verification.
// @Summary Verify stakeholder step-up authorization challenge
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param challenge_id path string true "Challenge ID"
// @Param payload body models.VerifyAuthorizationRequest false "Verification payload"
// @Success 200 {object} models.AuthorizationChallengeDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/authorizations/{challenge_id}/verify [post]
func stakeholderAuthorizationVerifyDocs() {}

// stakeholderDocumentCreateDocs documents document creation.
// @Summary Create stakeholder document metadata
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param payload body models.CreateDocumentRequest true "Legacy external document metadata. Prefer POST /stakeholder/shared/documents/upload for managed files."
// @Success 201 {object} models.DocumentDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/documents [post]
func stakeholderDocumentCreateDocs() {}

// stakeholderDocumentUploadDocs documents managed private uploads.
// @Summary Upload stakeholder document
// @Description Uploads a private PDF, JPEG, or PNG (max 10MB). The returned document id is used by valuation and fund-release requests.
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Accept multipart/form-data
// @Produce json
// @Param document_file formData file true "Document file"
// @Param asset_id formData string false "Assigned asset ID (required except for general and custody_compliance)"
// @Param category formData string true "Document category from /stakeholder/shared/document-categories"
// @Param title formData string true "Document title"
// @Param access_roles formData []string false "Additional dashboard roles"
// @Success 201 {object} models.DocumentDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Failure 413 {object} models.StakeholderErrorResponse
// @Failure 503 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/documents/upload [post]
func stakeholderDocumentUploadDocs() {}

// stakeholderDocumentCategoriesDocs documents the category catalogue.
// @Summary List stakeholder document categories
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/document-categories [get]
func stakeholderDocumentCategoriesDocs() {}

// stakeholderDocumentListDocs documents document list.
// @Summary List stakeholder documents
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param asset_id query string false "Asset ID"
// @Success 200 {object} models.DocumentPageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/documents [get]
func stakeholderDocumentListDocs() {}

// stakeholderDocumentGetDocs documents document detail.
// @Summary Get stakeholder document
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Produce json
// @Param doc_id path string true "Document ID"
// @Success 200 {object} models.DocumentDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/documents/{doc_id} [get]
func stakeholderDocumentGetDocs() {}

// stakeholderDocumentDownloadDocs streams an authorized private document.
// @Summary Download stakeholder document
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Produce application/octet-stream
// @Param doc_id path string true "Document ID"
// @Success 200 {file} binary
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Failure 422 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/documents/{doc_id}/download [get]
func stakeholderDocumentDownloadDocs() {}

// stakeholderNotificationDocs documents notification routes.
// @Summary List stakeholder notifications
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.NotificationPageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/notifications [get]
func stakeholderNotificationDocs() {}

// stakeholderNotificationReadDocs documents mark-read route.
// @Summary Mark stakeholder notification read
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} models.NotificationDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/notifications/{id}/read [put]
func stakeholderNotificationReadDocs() {}

// stakeholderAuditDocs documents audit trail.
// @Summary List stakeholder audit trail
// @Tags Stakeholder - Shared
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.AuditLogPageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/shared/audit-trail [get]
func stakeholderAuditDocs() {}

// stakeholderTrusteeDashboardDocs documents the trustee dashboard.
// @Summary Trustee dashboard
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.DashboardDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/dashboard [get]
func stakeholderTrusteeDashboardDocs() {}

// stakeholderTrusteeFundManagementSummaryDocs documents trustee fund management totals.
// @Summary Get trustee fund management summary
// @Description Returns trustee-scoped totals grouped by currency. Processed is completed issuer fund releases plus completed investor distributions; primary sales come from WalletDB subscriptions; milestone balance is active development-fund account balance; fees are configured tokenization fee values including VAT.
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.TrusteeFundManagementSummaryDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/fund-management/summary [get]
func stakeholderTrusteeFundManagementSummaryDocs() {}

// stakeholderCustodianDashboardDocs documents the custodian dashboard.
// @Summary Custodian dashboard
// @Tags Stakeholder - Custodian
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.DashboardDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/custodian/dashboard [get]
func stakeholderCustodianDashboardDocs() {}

// stakeholderManagerDashboardDocs documents the asset manager dashboard.
// @Summary Asset manager dashboard
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.DashboardDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/dashboard [get]
func stakeholderManagerDashboardDocs() {}

// stakeholderDueDiligenceDocs documents due diligence routes.
// @Summary Get trustee due diligence checklist
// @Description Returns the checklist with full active document metadata linked to each checklist item and available_asset_documents sourced read-only from the selected WalletDB tokenized asset.
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Produce json
// @Param asset_id path string true "Asset ID or code"
// @Success 200 {object} models.DueDiligenceChecklistDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/due-diligence/{asset_id} [get]
func stakeholderDueDiligenceDocs() {}

// stakeholderDueDiligenceItemDocs documents due diligence item update.
// @Summary Update trustee due diligence item
// @Description Updates status and notes. When document_ids is supplied, it replaces the item's due-diligence evidence with accessible managed documents for this asset.
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param asset_id path string true "Asset ID or code"
// @Param item_id path string true "Item ID"
// @Param payload body models.UpdateDueDiligenceItemRequest true "Due diligence item update"
// @Success 200 {object} models.DueDiligenceItemDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/due-diligence/{asset_id}/items/{item_id} [put]
func stakeholderDueDiligenceItemDocs() {}

// stakeholderDueDiligenceCategoryDocs documents category-level verification.
// @Summary Verify trustee due diligence category
// @Description Atomically updates every due diligence item in the requested category. Valid statuses are pending, complete, failed, and not_applicable.
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param asset_id path string true "Asset ID or code"
// @Param payload body models.UpdateDueDiligenceCategoryRequest true "Due diligence category update"
// @Success 200 {object} models.DueDiligenceCategoryDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/due-diligence/{asset_id}/categories [put]
func stakeholderDueDiligenceCategoryDocs() {}

// stakeholderDueDiligenceApproveDocs documents due diligence approval.
// @Summary Approve trustee due diligence
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Produce json
// @Param asset_id path string true "Asset ID or code"
// @Success 200 {object} models.DueDiligenceChecklistDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/due-diligence/{asset_id}/approve [post]
func stakeholderDueDiligenceApproveDocs() {}

// stakeholderDueDiligenceRejectDocs documents due diligence rejection.
// @Summary Reject trustee due diligence
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param asset_id path string true "Asset ID or code"
// @Param payload body models.RejectRequest true "Rejection reason"
// @Success 200 {object} models.DueDiligenceChecklistDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/due-diligence/{asset_id}/reject [post]
func stakeholderDueDiligenceRejectDocs() {}

// stakeholderFundReleaseCreateDocs documents manager fund release creation.
// @Summary Create fund release request
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param payload body models.CreateFundReleaseRequest true "Fund release request with managed document IDs and optional receiving-account snapshot"
// @Success 201 {object} models.FundReleaseDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/fund-releases [post]
func stakeholderFundReleaseCreateDocs() {}

// stakeholderTrusteeFundReleaseListDocs documents the trustee fund release list.
// @Summary List trustee fund release requests
// @Description Each record includes resolved requester and reviewer member/organization details while retaining the original IDs.
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Produce json
// @Param status query string false "Filter by status"
// @Param asset_id query string false "Filter by assigned asset ID or code"
// @Success 200 {object} models.FundReleasePageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/fund-releases [get]
func stakeholderTrusteeFundReleaseListDocs() {}

// stakeholderCustodianFundReleaseListDocs documents the custodian fund release list.
// @Summary List custodian fund release requests
// @Description Each record includes resolved requester and reviewer member/organization details while retaining the original IDs.
// @Tags Stakeholder - Custodian
// @Security OrganizationAuth
// @Produce json
// @Param status query string false "Filter by status"
// @Param asset_id query string false "Filter by assigned asset ID or code"
// @Success 200 {object} models.FundReleasePageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/custodian/fund-releases [get]
func stakeholderCustodianFundReleaseListDocs() {}

// stakeholderManagerFundReleaseListDocs documents the asset manager fund release list.
// @Summary List asset manager fund release requests
// @Description Each record includes resolved requester and reviewer member/organization details while retaining the original IDs.
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Produce json
// @Param status query string false "Filter by status"
// @Param asset_id query string false "Filter by assigned asset ID or code"
// @Success 200 {object} models.FundReleasePageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/fund-releases [get]
func stakeholderManagerFundReleaseListDocs() {}

// stakeholderTrusteeFundReleaseGetDocs documents trustee fund release detail.
// @Summary Get enriched trustee fund release request
// @Description Includes requester identity, complete supporting-document metadata and the receiving-account snapshot.
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Produce json
// @Param request_id path string true "Fund release request ID"
// @Success 200 {object} models.FundReleaseDetailDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/fund-releases/{request_id} [get]
func stakeholderTrusteeFundReleaseGetDocs() {}

// stakeholderManagerFundReleaseGetDocs documents asset-manager fund release detail.
// @Summary Get enriched asset manager fund release request
// @Description Includes requester identity, complete supporting-document metadata and the receiving-account snapshot.
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Produce json
// @Param request_id path string true "Fund release request ID"
// @Success 200 {object} models.FundReleaseDetailDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/fund-releases/{request_id} [get]
func stakeholderManagerFundReleaseGetDocs() {}

// stakeholderCustodianFundReleaseGetDocs documents custodian fund release detail.
// @Summary Get enriched custodian fund release request
// @Description Includes requester identity, complete supporting-document metadata and the receiving-account snapshot.
// @Tags Stakeholder - Custodian
// @Security OrganizationAuth
// @Produce json
// @Param request_id path string true "Fund release request ID"
// @Success 200 {object} models.FundReleaseDetailDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/custodian/fund-releases/{request_id} [get]
func stakeholderCustodianFundReleaseGetDocs() {}

// stakeholderManagerFundSummaryDocs documents asset-manager fund-management totals.
// @Summary Get asset manager fund-management summary
// @Description Returns requested, approved, released, pending and rejected totals by currency. Remaining balance is approved minus released.
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Produce json
// @Param asset_id query string false "Assigned asset ID or code"
// @Success 200 {object} models.AssetManagerFundSummaryDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/fund-releases/summary [get]
func stakeholderManagerFundSummaryDocs() {}

// stakeholderFundReleaseApproveDocs documents fund release approval.
// @Summary Approve fund release request
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param request_id path string true "Fund release request ID"
// @Param payload body models.ChallengeActionRequest true "Verified challenge"
// @Success 200 {object} models.FundReleaseDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/fund-releases/{request_id}/approve [post]
func stakeholderFundReleaseApproveDocs() {}

// stakeholderFundReleaseRejectDocs documents fund release rejection.
// @Summary Reject fund release request
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param request_id path string true "Fund release request ID"
// @Param payload body models.RejectRequest true "Rejection reason"
// @Success 200 {object} models.FundReleaseDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/fund-releases/{request_id}/reject [post]
func stakeholderFundReleaseRejectDocs() {}

// stakeholderFundReleaseExecuteDocs documents fund release execution.
// @Summary Execute fund release request
// @Tags Stakeholder - Custodian
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param request_id path string true "Fund release request ID"
// @Param payload body models.ExecuteFundReleaseRequest true "Execution request"
// @Success 200 {object} models.FundReleaseDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/custodian/fund-releases/{request_id}/execute [post]
func stakeholderFundReleaseExecuteDocs() {}

// stakeholderFundReleaseStatusDocs documents execution status update.
// @Summary Update fund release execution status
// @Tags Stakeholder - Custodian
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param request_id path string true "Fund release request ID"
// @Param payload body models.UpdateFundReleaseStatusRequest true "Execution status"
// @Success 200 {object} models.FundReleaseDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/custodian/fund-releases/{request_id}/status [put]
func stakeholderFundReleaseStatusDocs() {}

// stakeholderDistributionDocs documents trustee distribution list.
// @Summary List trustee distributions
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.DistributionPageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/distributions [get]
// @Router /stakeholder/trustee/distributions/history [get]
func stakeholderDistributionDocs() {}

// stakeholderDistributionDetailDocs documents trustee distribution detail.
// @Summary Get trustee distribution details
// @Description Returns the role-scoped distribution and its canonical token-supply breakdown. Distribution per token is calculated as net income divided by total issued tokens.
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Produce json
// @Param dist_id path string true "Distribution ID"
// @Success 200 {object} models.DistributionDetailDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/distributions/{dist_id} [get]
func stakeholderDistributionDetailDocs() {}

// stakeholderDistributionPayoutDownloadDocs documents the complete payout-list export.
// @Summary Download complete distribution payout list
// @Description Downloads every Wallet payout record for the trustee-owned distribution as a CSV file.
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Produce text/csv
// @Param dist_id path string true "Distribution ID"
// @Success 200 {file} binary
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Failure 503 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/distributions/{dist_id}/payouts/download [get]
func stakeholderDistributionPayoutDownloadDocs() {}

// stakeholderDistributionAuthorizeDocs documents distribution authorization.
// @Summary Authorize distribution
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param dist_id path string true "Distribution ID"
// @Param payload body models.ChallengeActionRequest true "Verified challenge"
// @Success 200 {object} models.DistributionDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/distributions/{dist_id}/authorize [post]
func stakeholderDistributionAuthorizeDocs() {}

// stakeholderDistributionRejectDocs documents distribution rejection.
// @Summary Reject distribution
// @Tags Stakeholder - Trustee
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param dist_id path string true "Distribution ID"
// @Param payload body models.RejectRequest true "Rejection reason"
// @Success 200 {object} models.DistributionDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/trustee/distributions/{dist_id}/reject [post]
func stakeholderDistributionRejectDocs() {}

// stakeholderRevenueDocs documents revenue APIs.
// @Summary List asset manager revenue
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.RevenuePageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/revenue [get]
func stakeholderRevenueDocs() {}

// stakeholderRevenueCreateDocs documents revenue creation.
// @Summary Record asset revenue
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param payload body models.CreateRevenueRequest true "Revenue record"
// @Success 201 {object} models.RevenueDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/revenue [post]
func stakeholderRevenueCreateDocs() {}

// stakeholderDistributionSubmitDocs documents distribution submission.
// @Summary Submit revenue for distribution
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param payload body models.SubmitDistributionRequest true "Distribution proposal"
// @Success 201 {object} models.DistributionDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/revenue/submit-distribution [post]
func stakeholderDistributionSubmitDocs() {}

// stakeholderValuationListDocs documents valuation list.
// @Summary List asset valuation history
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Produce json
// @Param asset_id path string true "Asset ID"
// @Success 200 {object} models.ValuationPageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/valuations/{asset_id} [get]
func stakeholderValuationListDocs() {}

// stakeholderValuationCreateDocs documents valuation creation.
// @Summary Submit asset valuation
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param payload body models.CreateValuationRequest true "Asset valuation"
// @Success 201 {object} models.ValuationDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/valuations [post]
func stakeholderValuationCreateDocs() {}

// stakeholderValuationIndependentDocs documents independent valuation request.
// @Summary Request independent valuation
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Produce json
// @Param id path string true "Valuation ID"
// @Success 200 {object} models.ValuationDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/valuations/{id}/request-independent [post]
func stakeholderValuationIndependentDocs() {}

// stakeholderCustodianAccountDocs documents custodian accounts.
// @Summary List custodian segregated accounts
// @Tags Stakeholder - Custodian
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.SegregatedAccountPageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/custodian/accounts [get]
func stakeholderCustodianAccountDocs() {}

// stakeholderCustodianAccountCreateDocs documents custodian account creation.
// @Summary Create custodian segregated account
// @Description Creates an account for an asset assigned to the authenticated custodian organization.
// @Tags Stakeholder - Custodian
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param payload body models.CreateSegregatedAccountRequest true "Segregated account"
// @Success 201 {object} models.SegregatedAccountDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Failure 422 {object} models.StakeholderErrorResponse
// @Router /stakeholder/custodian/accounts [post]
func stakeholderCustodianAccountCreateDocs() {}

// stakeholderCustodianAccountGetDocs documents account detail.
// @Summary Get custodian segregated account
// @Tags Stakeholder - Custodian
// @Security OrganizationAuth
// @Produce json
// @Param account_id path string true "Account ID"
// @Success 200 {object} models.SegregatedAccountDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/custodian/accounts/{account_id} [get]
func stakeholderCustodianAccountGetDocs() {}

// stakeholderCustodianAccountUpdateDocs documents account update.
// @Summary Update custodian segregated account
// @Description Updates an account owned by the authenticated custodian organization.
// @Tags Stakeholder - Custodian
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID"
// @Param payload body models.UpdateSegregatedAccountRequest true "Segregated account update"
// @Success 200 {object} models.SegregatedAccountDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Router /stakeholder/custodian/accounts/{account_id} [put]
func stakeholderCustodianAccountUpdateDocs() {}

// stakeholderComplianceDocs documents compliance list.
// @Summary List custodian compliance items
// @Description Returns compliance requirements and the full active document metadata linked to each item.
// @Tags Stakeholder - Custodian
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.ComplianceItemPageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/custodian/compliance [get]
func stakeholderComplianceDocs() {}

// stakeholderComplianceUpdateDocs documents compliance update.
// @Summary Update custodian compliance item
// @Description Updates status. When document_ids is supplied, it replaces the item's custody-compliance evidence with managed documents uploaded by this custodian for the same asset.
// @Tags Stakeholder - Custodian
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param item_id path string true "Compliance item ID"
// @Param payload body models.UpdateComplianceItemRequest true "Compliance update"
// @Success 200 {object} models.ComplianceItemDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/custodian/compliance/{item_id} [put]
func stakeholderComplianceUpdateDocs() {}

// stakeholderAssetOperationDocs documents manager asset operation update.
// @Summary Update manager asset operational info
// @Description Stores a dashboard-local operational overlay and returns it in asset details as operational_update. This does not change WalletDB assetTokenizationStatus. Allowed statuses: active, inactive, paused, maintenance. metadata accepts source and an RFC3339 effective_at.
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param asset_id path string true "Asset ID"
// @Param payload body models.UpdateAssetOperationRequest true "Operational update"
// @Success 200 {object} models.AssetOperationDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/assets/{asset_id} [put]
func stakeholderAssetOperationDocs() {}

// stakeholderReportDocs documents report creation/list.
// @Summary List asset manager reports
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.ReportPageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/reports [get]
func stakeholderReportDocs() {}

// stakeholderReportCreateDocs documents report creation.
// @Summary Generate asset manager report
// @Description Allowed report_type values are income and operational. Upload a managed asset_report document first and pass its document_id. file_url remains a temporary legacy fallback.
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param payload body models.CreateReportRequest true "Report"
// @Success 201 {object} models.ReportDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/reports [post]
func stakeholderReportCreateDocs() {}

// stakeholderReportTypesDocs documents the backend report-type catalogue.
// @Summary List allowed asset manager report types
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.ReportTypeListDataResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/report-types [get]
func stakeholderReportTypesDocs() {}

// stakeholderReportSubmitDocs documents report submission.
// @Summary Submit asset manager report
// @Tags Stakeholder - Asset Manager
// @Security OrganizationAuth
// @Produce json
// @Param report_id path string true "Report ID"
// @Success 200 {object} models.ReportDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/asset-manager/reports/{report_id}/submit [post]
func stakeholderReportSubmitDocs() {}

// stakeholderAdminAssignmentDocs documents admin assignment creation.
// @Summary Create stakeholder asset assignment
// @Tags Stakeholder - Admin
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param payload body models.CreateAssetAssignmentRequest true "Asset assignment"
// @Success 200 {object} models.AssetAssignmentDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/asset-assignments [post]
func stakeholderAdminAssignmentDocs() {}

// stakeholderAdminComplianceCreateDocs documents admin compliance creation.
// @Summary Create custodian compliance item
// @Description Assigns a pending compliance requirement to an active organization linked to an Asset Custodian stakeholder. asset_id is optional; when supplied, it must be assigned to that custodian.
// @Tags Stakeholder - Admin
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param payload body models.CreateComplianceItemRequest true "Custodian compliance requirement"
// @Success 201 {object} models.ComplianceItemDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Failure 409 {object} models.StakeholderErrorResponse
// @Failure 422 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/custodian-compliance [post]
func stakeholderAdminComplianceCreateDocs() {}

// stakeholderAdminComplianceListDocs documents the admin compliance list.
// @Summary List custodian compliance requirements for Admin
// @Description Returns all compliance requirements created for custodian organizations. Supports pagination and optional custodian, asset, and status filters.
// @Tags Stakeholder - Admin
// @Security JwtTokenAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Records per page" default(20) maximum(100)
// @Param custodian_org_id query string false "Custodian organization ID"
// @Param asset_id query string false "Tokenized asset ID"
// @Param status query string false "Compliance status" Enums(pending,complete,overdue,waived)
// @Success 200 {object} models.ComplianceItemPageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/custodian-compliance [get]
func stakeholderAdminComplianceListDocs() {}

// stakeholderComplianceTemplateListDocs documents compliance template listing.
// @Summary List organization compliance templates
// @Tags Stakeholder - Compliance Admin
// @Security JwtTokenAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Records per page" default(20) maximum(100)
// @Param org_type query string false "Stakeholder organization type" Enums(asset_manager,asset_issuing_house,approved_asset_custodian,legal_and_professionals,rating_agency,trustees,legal_adviser,financial_adviser)
// @Param level query int false "Positive compliance level"
// @Success 200 {object} models.ComplianceRequirementTemplatePageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/compliance-templates [get]
func stakeholderComplianceTemplateListDocs() {}

// stakeholderComplianceTemplateCreateDocs documents compliance template creation.
// @Summary Create an organization compliance template
// @Description Creates the unique template for one stakeholder org_type and positive level. Item input_type must be document_upload, text, or structured_form. required defaults to true when omitted.
// @Tags Stakeholder - Compliance Admin
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param payload body models.CreateComplianceTemplateRequest true "Compliance template"
// @Success 201 {object} models.ComplianceRequirementTemplateDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 409 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/compliance-templates [post]
func stakeholderComplianceTemplateCreateDocs() {}

// stakeholderComplianceTemplateUpdateDocs documents compliance template updates.
// @Summary Update an organization compliance template
// @Tags Stakeholder - Compliance Admin
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Param payload body models.UpdateComplianceTemplateRequest true "Replacement template definition; include each existing item id to preserve its identity when renaming it"
// @Success 200 {object} models.ComplianceRequirementTemplateDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Failure 409 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/compliance-templates/{id} [put]
func stakeholderComplianceTemplateUpdateDocs() {}

// stakeholderComplianceRequirementListDocs documents the admin review queue.
// @Summary List organization compliance requirements
// @Tags Stakeholder - Compliance Admin
// @Security JwtTokenAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Records per page" default(20) maximum(100)
// @Param org_id query string false "Organization ID"
// @Param category query string false "Requirement category"
// @Param level query int false "Positive organization compliance level"
// @Param status query string false "Requirement status" Enums(draft,pending,submitted,under_review,approved,rejected,overdue,waived)
// @Success 200 {object} models.ComplianceRequirementInstancePageDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/compliance-requirements [get]
func stakeholderComplianceRequirementListDocs() {}

// stakeholderComplianceRequirementCreateDocs documents ad-hoc assignment.
// @Summary Create an ad-hoc organization compliance requirement
// @Description Creates an organization-level requirement independently of the asset-specific custodian compliance workflow.
// @Tags Stakeholder - Compliance Admin
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param payload body models.CreateComplianceRequirementRequest true "Ad-hoc requirement"
// @Success 201 {object} models.ComplianceRequirementInstanceDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Failure 409 {object} models.StakeholderErrorResponse
// @Failure 422 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/compliance-requirements [post]
func stakeholderComplianceRequirementCreateDocs() {}

// stakeholderComplianceRequirementDetailDocs documents requirement detail.
// @Summary Get organization compliance requirement details
// @Description Read-only; opening this endpoint does not change the requirement status.
// @Tags Stakeholder - Compliance Admin
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Requirement ID"
// @Success 200 {object} models.ComplianceRequirementInstanceDataResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/compliance-requirements/{id} [get]
func stakeholderComplianceRequirementDetailDocs() {}

// stakeholderComplianceRequirementStartReviewDocs documents explicit review start.
// @Summary Start review of a submitted compliance requirement
// @Description Idempotently changes submitted to under_review and records the admin in the audit trail.
// @Tags Stakeholder - Compliance Admin
// @Security JwtTokenAuth
// @Produce json
// @Param id path string true "Requirement ID"
// @Success 200 {object} models.ComplianceRequirementInstanceDataResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Failure 409 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/compliance-requirements/{id}/review/start [put]
func stakeholderComplianceRequirementStartReviewDocs() {}

// stakeholderComplianceRequirementReviewDocs documents approval/rejection.
// @Summary Approve or reject a compliance requirement
// @Tags Stakeholder - Compliance Admin
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param id path string true "Requirement ID"
// @Param payload body models.ReviewComplianceRequirementRequest true "Review decision"
// @Success 200 {object} models.ComplianceRequirementInstanceDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Failure 409 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/compliance-requirements/{id}/review [put]
func stakeholderComplianceRequirementReviewDocs() {}

// stakeholderComplianceRequirementAdminDocumentDocs documents admin evidence download.
// @Summary Download submitted compliance evidence as Admin
// @Tags Stakeholder - Compliance Admin
// @Security JwtTokenAuth
// @Produce application/octet-stream
// @Param id path string true "Requirement ID"
// @Param doc_id path string true "Submitted document ID"
// @Success 200 {file} file
// @Failure 404 {object} models.StakeholderErrorResponse
// @Router /stakeholder/admin/compliance-requirements/{id}/documents/{doc_id}/download [get]
func stakeholderComplianceRequirementAdminDocumentDocs() {}

// stakeholderOrgComplianceListDocs documents organization requirements.
// @Summary List requirements assigned to the authenticated organization
// @Tags Stakeholder - Organization Compliance
// @Security OrganizationAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Records per page" default(20) maximum(100)
// @Param status query string false "Requirement status" Enums(draft,pending,submitted,under_review,approved,rejected,overdue,waived)
// @Success 200 {object} models.ComplianceRequirementInstancePageDataResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Router /stakeholder/org/compliance [get]
func stakeholderOrgComplianceListDocs() {}

// stakeholderOrgComplianceDraftDocs documents draft transition.
// @Summary Mark a pending organization compliance requirement as draft
// @Tags Stakeholder - Organization Compliance
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param item_id path string true "Requirement ID"
// @Param payload body models.UpdateOrgComplianceRequirementRequest true "Status must be draft"
// @Success 200 {object} models.ComplianceRequirementInstanceDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Failure 409 {object} models.StakeholderErrorResponse
// @Router /stakeholder/org/compliance/{item_id} [put]
func stakeholderOrgComplianceDraftDocs() {}

// stakeholderOrgComplianceSubmitDocs documents typed requirement submission.
// @Summary Submit an organization compliance requirement
// @Description submission_data shape depends on input_type: document_upload uses {"document_id":"uuid"}; text uses {"text":"value"}; structured_form uses {"fields":{"field":"value"}}.
// @Tags Stakeholder - Organization Compliance
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param item_id path string true "Requirement ID"
// @Param payload body models.SubmitComplianceRequirementRequest true "Typed submission"
// @Success 200 {object} models.ComplianceRequirementInstanceDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 404 {object} models.StakeholderErrorResponse
// @Failure 409 {object} models.StakeholderErrorResponse
// @Failure 422 {object} models.StakeholderErrorResponse
// @Router /stakeholder/org/compliance/{item_id}/submit [post]
func stakeholderOrgComplianceSubmitDocs() {}

// stakeholderOrgComplianceUploadDocs documents managed evidence upload.
// @Summary Upload private evidence for an organization compliance requirement
// @Description Uploads a managed private compliance_requirement document. Use the returned id as submission_data.document_id.
// @Tags Stakeholder - Organization Compliance
// @Security OrganizationAuth
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Document title"
// @Param document_file formData file true "PDF, JPEG, or PNG; maximum 10MB"
// @Success 201 {object} models.DocumentDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 413 {object} models.StakeholderErrorResponse
// @Failure 503 {object} models.StakeholderErrorResponse
// @Router /stakeholder/org/compliance/documents/upload [post]
func stakeholderOrgComplianceUploadDocs() {}

// stakeholderOrgComplianceDocumentDownloadDocs documents owner download.
// @Summary Download organization compliance evidence
// @Tags Stakeholder - Organization Compliance
// @Security OrganizationAuth
// @Produce application/octet-stream
// @Param doc_id path string true "Document ID"
// @Success 200 {file} file
// @Failure 404 {object} models.StakeholderErrorResponse
// @Router /stakeholder/org/compliance/documents/{doc_id}/download [get]
func stakeholderOrgComplianceDocumentDownloadDocs() {}

// stakeholderLegalAdviserDashboardDocs documents the legal adviser dashboard.
// @Summary Legal adviser dashboard
// @Description Aggregate view for a Legal Adviser: assigned assets, portfolio summary, and recent activity.
// @Tags Stakeholder - Legal Adviser
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.DashboardDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/legal-adviser/dashboard [get]
func stakeholderLegalAdviserDashboardDocs() {}

// stakeholderFinancialAdviserDashboardDocs documents the financial adviser dashboard.
// @Summary Financial adviser dashboard
// @Description Aggregate view for a Financial Adviser: assigned assets, portfolio summary, and recent activity.
// @Tags Stakeholder - Financial Adviser
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.DashboardDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/financial-adviser/dashboard [get]
func stakeholderFinancialAdviserDashboardDocs() {}

// stakeholderIssuingHouseDashboardDocs documents the issuing house dashboard.
// @Summary Issuing house dashboard (view-only)
// @Description View-only aggregate for an Issuing House: assigned assets, portfolio summary, and recent activity.
// @Tags Stakeholder - Issuing House
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.DashboardDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/issuing-house/dashboard [get]
func stakeholderIssuingHouseDashboardDocs() {}

// stakeholderRatingAgencyDashboardDocs documents the rating agency dashboard.
// @Summary Rating agency dashboard (view-only)
// @Description View-only aggregate for a Rating Agency: assigned assets, portfolio summary, and recent activity.
// @Tags Stakeholder - Rating Agency
// @Security OrganizationAuth
// @Produce json
// @Success 200 {object} models.DashboardDataResponse
// @Failure 400 {object} models.StakeholderErrorResponse
// @Failure 401 {object} models.StakeholderErrorResponse
// @Failure 403 {object} models.StakeholderErrorResponse
// @Router /stakeholder/rating-agency/dashboard [get]
func stakeholderRatingAgencyDashboardDocs() {}

// stakeholderLegalAdviserStructuringDocs documents legal adviser A5 confirmation.
// @Summary Confirm legal adviser structuring complete
// @Description Legal Adviser confirms their A5 legal-documentation workstream is complete for an assigned asset. Optional evidence document. (PRD §5A.3 / OI-14)
// @Tags Stakeholder - Legal Adviser
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param asset_id path string true "Asset ID or code"
// @Param payload body models.ConfirmStructuringRequest true "Structuring confirmation"
// @Success 200 {object} models.StructuringStatusDataResponse
// @Failure 403 {object} models.StakeholderErrorResponse "Not a legal adviser"
// @Failure 404 {object} models.StakeholderErrorResponse "Asset not assigned to this adviser"
// @Failure 409 {object} models.StakeholderErrorResponse "Workstream already confirmed"
// @Router /stakeholder/legal-adviser/structuring/{asset_id}/complete [post]
func stakeholderLegalAdviserStructuringDocs() {}

// stakeholderFinancialAdviserStructuringDocs documents financial adviser A5 confirmation.
// @Summary Confirm financial adviser structuring complete
// @Description Financial Adviser confirms their A5 transaction-structuring workstream is complete for an assigned asset. Optional evidence document. (PRD §5A.3 / OI-14)
// @Tags Stakeholder - Financial Adviser
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param asset_id path string true "Asset ID or code"
// @Param payload body models.ConfirmStructuringRequest true "Structuring confirmation"
// @Success 200 {object} models.StructuringStatusDataResponse
// @Failure 403 {object} models.StakeholderErrorResponse "Not a financial adviser"
// @Failure 404 {object} models.StakeholderErrorResponse "Asset not assigned to this adviser"
// @Failure 409 {object} models.StakeholderErrorResponse "Workstream already confirmed"
// @Router /stakeholder/financial-adviser/structuring/{asset_id}/complete [post]
func stakeholderFinancialAdviserStructuringDocs() {}
