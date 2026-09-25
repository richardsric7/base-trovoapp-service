package usermetrics

import (
	"fmt"
	"strings"

	"admin-panel-dashboard/internal/models"

	"gorm.io/gorm"
)

const fiatPaymentsBaseQuery = `
	SELECT
		id,
		service_provider,
		username,
		transaction_id,
		amount,
		payment_type,
		created_at,
		status,
		refunded,
		record_type
	FROM (
		SELECT
			CAST(id AS TEXT) AS id,
			service_provider,
			username,
			transaction_id,
			amount,
			payment_type,
			created_at,
			NULL::text AS status,
			NULL::bigint AS refunded,
			'$PAYMENT$' AS record_type
		FROM fiat_payments
		UNION ALL
		SELECT
			id,
			service_provider,
			username,
			NULL::text AS transaction_id,
			amount,
			payment_type,
			created_at,
			status,
			refunded,
			'$INVOICE$' AS record_type
		FROM fiat_payment_invoices
	) AS combined
`

// GetPaginatedFiatPaymentRecords fetches fiat payment data from both fiat_payments and fiat_payment_invoices tables.
func GetPaginatedFiatPaymentRecords(db *gorm.DB, filter models.FiatPaymentListFilter) (*models.PaginatedFiatPaymentResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	queryTemplate := strings.ReplaceAll(fiatPaymentsBaseQuery, "$PAYMENT$", models.FiatRecordTypePayment)
	queryTemplate = strings.ReplaceAll(queryTemplate, "$INVOICE$", models.FiatRecordTypeInvoice)

	var conditions []string
	var args []interface{}

	if filter.RecordType != "" {
		conditions = append(conditions, "record_type = ?")
		args = append(args, filter.RecordType)
	}
	if filter.ServiceProvider != "" {
		conditions = append(conditions, "service_provider ILIKE ?")
		args = append(args, "%"+filter.ServiceProvider+"%")
	}
	if filter.Username != "" {
		conditions = append(conditions, "username ILIKE ?")
		args = append(args, "%"+filter.Username+"%")
	}
	if filter.PaymentType != "" {
		conditions = append(conditions, "payment_type = ?")
		args = append(args, filter.PaymentType)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(1) FROM (%s) AS aggregated %s", queryTemplate, whereClause)
	var total int64
	if err := db.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	dataQuery := fmt.Sprintf("%s %s ORDER BY created_at DESC OFFSET ? LIMIT ?", queryTemplate, whereClause)
	dataArgs := append(args, offset, filter.PageSize)

	var records []models.FiatPaymentRecord
	if err := db.Raw(dataQuery, dataArgs...).Scan(&records).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize))

	response := &models.PaginatedFiatPaymentResponse{
		Data:       records,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}

	return response, nil
}
