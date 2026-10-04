package users

import (
	"net/http"
	"time"
)

// outboundHTTP is the client for calls to outside services (IP geolocation,
// ...): with a timeout, unlike http.DefaultClient.
var outboundHTTP = &http.Client{Timeout: 15 * time.Second}
