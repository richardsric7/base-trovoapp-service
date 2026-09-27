package root

// Doc-only stub for the "/" route registered in Init (main.go) as an inline
// gin.HandlerFunc closure. See users/controllers/swagger_docs.go for why
// this pattern is used.

// @Summary Get service info
// @Description Unauthenticated root endpoint. Returns a static service name and status string, used as a basic "is this the right service" smoke check.
// @Tags Meta
// @Produce json
// @Success 200 {object} map[string]string "service, status"
// @Router / [get]
func rootInfoDocs() {}
