package users

import (
	"net/http"
	"time"
)

// outboundHTTP is the client for calls to outside services (1Liquidity,
// Sumsub, ...): with a timeout, so a provider that stops answering cannot
// hold requests, their goroutines and memory forever (http.DefaultClient
// has none).
var outboundHTTP = &http.Client{Timeout: 30 * time.Second}
