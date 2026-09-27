package swagger

// Doc-only stub for the "/ping" route registered in Init (main.go) as an
// inline gin.HandlerFunc closure. See users/controllers/swagger_docs.go for
// why this pattern is used.
//
// /swagger/*any (mounted at both / and /api/v1) is deliberately left
// undocumented here: it serves the Swagger UI's own static assets and the
// generated swagger.json/.yaml themselves, not a versioned API operation -
// documenting "GET /swagger/*any" as a single operation would be misleading
// (its response shape is arbitrary static content, and its {any} segment
// isn't a meaningful path parameter for API consumers).

// @Summary Ping
// @Description Unauthenticated liveness smoke check - always returns "pong" if the process is up and routing requests. Distinct from /health, which also reports on DB/cache dependency status.
// @Tags Meta
// @Produce json
// @Success 200 {object} map[string]string "message": "pong"
// @Router /ping [get]
func pingDocs() {}
