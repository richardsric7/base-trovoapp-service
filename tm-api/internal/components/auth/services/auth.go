package auth

import (
	"admin-panel-dashboard/internal/components/accesslog"
	userServices "admin-panel-dashboard/internal/components/users/services"
	conDB "admin-panel-dashboard/internal/db"
	p2pErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	"admin-panel-dashboard/internal/observe"
	serverModels "admin-panel-dashboard/internal/server/models"

	serverResponse "admin-panel-dashboard/internal/server/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ecnepsnai/discord"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Handles login callback requests.
// @Description Handles login callback requests.
// @ID LoginCallbackHandler
// @Tags Callbacks
// @Produce json
// @Param serviceName path string true "Service Name"
// @Success 200 {object} models.LoginCallbackInput
// @Failure 400,401,500 {object} object
// @Router /v1/callbacks/login/{serviceName} [post]
// func LoginCallbackHandler(db *gorm.DB, gc *models.GlobalConfig) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var err error
//
//		var callbakInput models.LoginCallbackInput
//		data, _ := io.ReadAll(c.Request.Body)
//
//		err = json.Unmarshal(data, &callbakInput)
//
//		var invalidJSON p2pErrors.ErrorInvalidJSON
//
//		if err != nil {
//			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
//			return
//		}
//		sn := c.Param("serviceName")
//		conDB.PrintDBStats("POST /v1/callbacks/login/"+sn, db)
//
//		targetUser := callbakInput.TargetUser
//		if !strings.Contains(targetUser, "@"+sn) {
//			targetUser = fmt.Sprintf("%v@%v", callbakInput.TargetUser, sn)
//		}
//
//		if len(callbakInput.LoginID) == 0 {
//			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
//			return
//		}
//		{
//			//invalidate cache
//			cacheKey := fmt.Sprintf("[GET] /v1/users/verify/%v/%v", targetUser, callbakInput.LoginID)
//			s.GC.Cache.InvalidateCachedHttpResponse(cacheKey)
//		}
//		{
//			//check login requests/processes for the id
//			log.Printf("[LoginCallback] request received from [%v], Payload: [%+v]\n", sn, callbakInput)
//
//			pl := struct {
//				Username string `json:"username"`
//			}{
//				Username: targetUser,
//			}
//			gc.BroadcastToLoginID(callbakInput.LoginID, "loginNotice", pl)
//
//		}
//
//		c.JSON(http.StatusOK, callbakInput)
//	}
//}

func LoginCallback(c *gin.Context, s *serverModels.Server) {
	// gc := models.GlobalConfig{} // check to insert correctly

	// router.POST("/v1/callbacks/login/:serviceName", func(c *gin.Context) {
	var err error

	var callbakInput models.LoginCallbackInput
	data, _ := io.ReadAll(c.Request.Body)

	err = json.Unmarshal(data, &callbakInput)

	var invalidJSON p2pErrors.ErrorInvalidJSON
	if err != nil {
		c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
		return
	}
	sn := c.Param("serviceName")
	conDB.PrintDBStats("POST /v1/callbacks/login/"+sn, s.TrovoWalletDB)

	targetUser := callbakInput.TargetUser
	if !strings.Contains(targetUser, "@"+sn) {
		targetUser = fmt.Sprintf("%v@%v", callbakInput.TargetUser, sn)
	}

	if len(callbakInput.LoginID) == 0 {
		c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
		return
	}
	{
		// invalidate cache
		cacheKey := fmt.Sprintf("[GET] /v1/users/verify/%v/%v", targetUser, callbakInput.LoginID)
		s.GC.Cache.InvalidateCachedHttpResponse(cacheKey)
	}
	{
		// check login requests/processes for the id
		log.Printf("[LoginCallback] request received from [%v], Payload: [%+v]\n", sn, callbakInput)

		pl := struct {
			Username string `json:"username"`
		}{
			Username: targetUser,
		}
		s.GC.BroadcastToLoginID(callbakInput.LoginID, "loginNotice", pl)

	}

	c.JSON(http.StatusOK, callbakInput)
}

func AuthCallback(c *gin.Context, s *serverModels.Server) {
	// gc := models.GlobalConfig{} // check to insert correctly

	// router.POST("/v1/callbacks/login/:serviceName", func(c *gin.Context) {
	var err error

	// perform login request action. on success, it will return loginID
	var callbakInput models.AuthorizationCallbackInput
	data, _ := io.ReadAll(c.Request.Body)

	err = json.Unmarshal(data, &callbakInput)

	var invalidJSON p2pErrors.ErrorInvalidJSON

	if err != nil {
		c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
		return
	}
	sn := c.Param("serviceName")
	conDB.PrintDBStats("POST /v1/callbacks/auth/"+sn, s.TrovoWalletDB)

	targetUser := callbakInput.TargetUser
	if !strings.Contains(targetUser, sn) {
		targetUser = fmt.Sprintf("%v@%v", callbakInput.TargetUser, sn)
	}
	if len(callbakInput.AuthID) == 0 {
		c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
		return
	}
	// get authid model from the table
	var pendingAuthorization models.PendingAuthorization
	e := s.TrovoWalletDB.Where("auth_id = ?", callbakInput.AuthID).First(&pendingAuthorization).Error
	if e != nil {
		log.Printf("unable to find authorization object amongst pending authorizations due to err:[%v], callback:[%+v]", e, callbakInput)

		c.JSON(http.StatusBadRequest, gin.H{"error": "could not process request", "message": "could not process request"})

		return
	}

	if pendingAuthorization.TargetUser != targetUser {
		log.Printf("target user from db is not same as target user from request. callback:[%+v]", callbakInput)

		c.JSON(http.StatusBadRequest, gin.H{"error": "could not process request", "message": "could not process request"})

		return
	}
	{
		// check authorizable processes for the id, eg. authorizing beneficiary, authorizing token release
		log.Printf("[AuthorizationCallback] request received:[%+v]\nServiceName:[%v]", callbakInput, sn)
		// notify listeners
		pl := struct {
			Username string `json:"username"`
		}{
			Username: targetUser,
		}
		s.GC.BroadcastToAuthID(callbakInput.AuthID, "authorizationNotice", pl)

	}

	c.JSON(http.StatusOK, callbakInput)
}

// Login handles the login process.
// @Summary User login
// @Description Performs user login based on provided credentials.
// @ID UserLogin
// @Accept json
// @Produce json
// @Tags Authentication
// @Param body body models.LoginInput true "User credentials for login"
// @Success 200 {object} trovosdk.LoginWithTrovoWalletData "Successful login response"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Router /login [post]
func Login(c *gin.Context, s *serverModels.Server) {

	// router.POST("/v1/callbacks/login/:serviceName", func(c *gin.Context) {
	var err error

	conDB.PrintDBStats("POST /v1/login/", s.TrovoWalletDB)

	// perform login request action. on success, it will return loginID
	var loginInput models.LoginInput
	data, _ := io.ReadAll(c.Request.Body)

	err = json.Unmarshal(data, &loginInput)

	var invalidJSON p2pErrors.ErrorInvalidJSON

	if err != nil {
		c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
		return
	}

	// Check admin and Trovo user validation
	if !isUserAnAdmin(s, loginInput) || !isUserTrovoUser(s, loginInput) {
		msg := fmt.Sprintf("unauthorized user: %v", loginInput.Username)
		accesslog.Record(s.AdminDB, accesslog.Event{
			Event:     models.EventLoginFailure,
			Status:    models.AccessStatusFailed,
			Category:  models.AccessCategorySession,
			Actor:     accessActorForUsername(s, loginInput.Username),
			IPAddress: accessClientIP(c),
			UserAgent: c.Request.UserAgent(),
			Method:    c.Request.Method,
			Path:      c.FullPath(),
			Detail:    "login gate rejected: not an active admin / Trovo user",
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		return
	}

	// if loginInput.LoginType != "bantupay" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "login type is not present"})
	// 	return
	// }
	loginInput.Username = strings.ToLower(strings.TrimSpace(loginInput.Username))
	userLogin := loginInput.Username
	if !strings.Contains(loginInput.Username, "@") {

		userLogin = loginInput.Username + "@trovo"
	}
	sn := strings.Split(userLogin, "@")[1]

	if sn == "trovo" {
		// process trovo login

		deviceInfo := ""
		callbackUrl := fmt.Sprintf("%v/%v", os.Getenv("LOGIN_CALLBACK_URL"), sn)
		loginData, err := s.GC.ServiceLink.SendLoginRequest(strings.Split(userLogin, "@")[0], "", deviceInfo, callbackUrl)
		if err != nil {
			var ex p2pErrors.GenericError
			var ok bool

			ex, ok = err.(p2pErrors.GenericError)
			statusCode := 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error(), "message": err.Error()}
			}

			c.JSON(statusCode, response)
			return
		}

		c.JSON(http.StatusOK, loginData)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Login service is not supported", "message": "Login service is not supported. Only trovo and all trovo compatible wallets are allowed. Send username as username@trovo"})

}

// accessActorForUsername resolves an admin row for a login-related audit event.
// The username may carry an "@trovo" suffix. On any miss it returns an actor with
// just the submitted username, so a failed login by a non-admin is still recorded.
func accessActorForUsername(s *serverModels.Server, username string) accesslog.Actor {
	uname := strings.ToLower(strings.TrimSpace(username))
	uname = strings.TrimSuffix(uname, "@trovo")

	actor := accesslog.Actor{Username: uname}

	var admin models.AdminUser
	if err := s.AdminDB.First(&admin, "username = ?", uname).Error; err == nil {
		id := admin.ID
		actor.AdminID = &id
		actor.Email = admin.Email
		actor.Role = string(admin.Role)
		name := strings.TrimSpace(admin.FirstName + " " + admin.LastName)
		actor.FullName = name
	}
	return actor
}

// accessClientIP mirrors the login flow's localhost substitution so dev logins
// still geo-resolve in the audit log.
func accessClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "127.0.0.1" || ip == "::1" {
		return "41.190.2.194"
	}
	return ip
}

func isUserAnAdmin(s *serverModels.Server, loginInput models.LoginInput) bool {
	var existingAdmin models.AdminUser
	username := loginInput.Username

	// Trim "@trovo" suffix if present
	if strings.HasSuffix(username, "@trovo") {
		username = strings.TrimSuffix(username, "@trovo")
		log.Println("Checking admin username without @trovo: ", username)
	}

	// Query the database for the admin user
	if err := s.AdminDB.First(&existingAdmin, "username = ?", username).Error; err != nil {
		return false // User is not an admin
	}

	// Check if the admin user is suspended
	if existingAdmin.Status == models.SUSPENDED {
		return false
	}

	return true // User is an active admin
}

// IsTrovoUser returns true if the given username exists in Trovo Wallet DB and is not suspended.
// Username can be with or without "@trovo" suffix. Used by Login (admin) and by organization wallet-link flow.
func IsTrovoUser(s *serverModels.Server, username string) bool {
	username = strings.TrimSpace(strings.ToLower(username))
	if strings.HasSuffix(username, "@trovo") {
		username = strings.TrimSuffix(username, "@trovo")
		log.Println("Checking Trovo user without @trovo: ", username)
	}

	var existingUser models.User
	if err := s.TrovoWalletDB.First(&existingUser, "username = ?", username).Error; err != nil {
		return false
	}
	if existingUser.Suspended == 1 {
		return false
	}
	return true
}

func isUserTrovoUser(s *serverModels.Server, loginInput models.LoginInput) bool {
	return IsTrovoUser(s, loginInput.Username)
}

// pendingLoginVerifyCacheSeconds bounds how long a "not yet approved" (or transiently failed)
// response from VerifyLoginID's Redis cache can be replayed before the next poll is forced to
// re-check with the external Trovo service. This must stay at or below the frontend's polling
// interval (qr-scan/page.tsx's useVerifyLoginQuery, pollingInterval: 5000ms) — previously this
// was 60s, which meant a poll that cached "pending" could keep replaying that same stale answer
// for up to a minute after the user had actually approved the login on their phone, since every
// later poll hit the cache and never rechecked. A few seconds still dedupes near-simultaneous
// polls (e.g. multiple tabs) without introducing minutes of staleness on a status endpoint whose
// entire purpose is to reflect the current state.
const pendingLoginVerifyCacheSeconds = 3

// VerifyLoginID verifies the login ID of a user.
// @Summary Verify login ID
// @Description Verifies the login ID of a user based on provided parameters.
// @ID VerifyLoginID
// @Tags Authentication
// @Accept json
// @Produce json
// @Param targetUser path string true "Target user's username"
// @Param loginID path string true "Login ID associated with the user"
// @Success 200 {object} object "Successful verification response"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal server error"
// @Router /users/verify/{targetUser}/{loginID} [get]
func VerifyLoginID(c *gin.Context, s *serverModels.Server) {
	// gc := models.GlobalConfig{} // check to insert correctly

	// router.POST("/v1/callbacks/login/:serviceName", func(c *gin.Context) {

	targetUser := strings.ReplaceAll(strings.ToLower(c.Param("targetUser")), " ", "")
	loginID := strings.TrimSpace(strings.ToLower(c.Param("loginID")))
	userLogin := targetUser

	cacheKey := fmt.Sprintf("[GET] /v1/users/verify/%v/%v", targetUser, loginID)
	{
		// search cache

		// cacheKeyParameters := fmt.Sprintf("limit=%v&order=%v&cursor=%v&forTransactionHash=%v&includeHash=%v&temp=%v", limit, orderStr, cursor, forTransactionHash, includeHash, temp)

		ok, status, response := s.GC.Cache.CachedHttpResponse(cacheKey)

		if ok {
			log.Printf("[%v], served from cache\n", cacheKey)
			c.JSON(status, response)
			return
		}
	}

	if !strings.Contains(targetUser, "@") {
		log.Println("target user does not have @domain attached to it: ", targetUser)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user does not have any @domain attached to it", "message": "user does not have any @domain attached to it"})
		return

	}

	if len(strings.Split(targetUser, "@")[0]) < 2 {

		c.JSON(http.StatusBadRequest, gin.H{"error": "valid user is not supplied", "message": "valid user is not supplied"})
		return

	}

	sn := strings.Split(userLogin, "@")[1]
	conDB.PrintDBStats(fmt.Sprintf("GET /v1/users/verify/%v/%v", userLogin, loginID), s.TrovoWalletDB)

	if sn == "trovo" {
		// process bantupay request
		responseData, err := s.GC.ServiceLink.VerifyLoginRequest(strings.Split(userLogin, "@")[0], loginID)
		if err != nil {
			log.Printf("[TrovoLoginVerifyFail] could not verify login data for user [%v] due to [%v]\n", strings.Split(userLogin, "@")[0], err)
			// NOTE: do NOT record a login.failure here. This is the poll endpoint:
			// the frontend calls it repeatedly (~1/s) while the user is still approving
			// the push/QR request, and every not-yet-approved poll lands in this same
			// error branch (it is cached only for pendingLoginVerifyCacheSeconds and then
			// retried). Recording a failure here produced a burst of false "Login failed"
			// audit rows for a single successful login. The genuine failure signals are
			// captured once, at login initiation: the gate-reject in Login (not an admin /
			// suspended / not a Trovo user) and the suspended-at-exchange check below.
			var ex p2pErrors.GenericError
			var ok bool

			ex, ok = err.(p2pErrors.GenericError)
			statusCode := 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error(), "message": err.Error()}
			}

			c.JSON(statusCode, response)
			s.GC.Cache.CacheHttpResponse(cacheKey, statusCode, response, pendingLoginVerifyCacheSeconds)
			return
		}

		// if len(userData.Mobile) == 0 {
		// 	log.Println("[Could not create User with no mobile phone]", err)
		// 	var serverError p2pErrors.ErrorTemporaryServerError
		// 	c.JSON(http.StatusInternalServerError, serverError.JSONError())
		// 	cacheDurationInSeconds := 60 //1 minute

		// 	s.GC.Cache.CacheHttpResponse(cacheKey, http.StatusInternalServerError, serverError.JSONError(), cacheDurationInSeconds)
		// 	return
		// }

		// check if user exists, then update the records with the latest
		var user models.User
		username := strings.Split(userLogin, "@")[0]
		err = s.TrovoWalletDB.Where("username = ?", username).First(&user).Error

		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("[Error Checking If User Exists] user [%v], error [%v]\n", userLogin, err)
				var serverError p2pErrors.ErrorTemporaryServerError
				c.JSON(http.StatusInternalServerError, serverError.JSONError())
				s.GC.Cache.CacheHttpResponse(cacheKey, http.StatusInternalServerError, serverError.JSONError(), pendingLoginVerifyCacheSeconds)

				return
			}
			// user does not exist, create it
			log.Printf("[FirstTime Login Detected] user [%v], error [%v]\n", userLogin, err)
			clientIP := c.ClientIP()
			if clientIP == "127.0.0.1" {
				clientIP = "41.190.2.194"
			}
			userData := responseData.UserInfo.UserData
			user = models.User{
				ID:       userData.ID,
				Username: userLogin,

				FirstName: userData.FirstName,
				Email:     userData.Email,
				PublicKey: userData.PublicKey,
				// Verified:  userData.Verified,
				PublicIP: clientIP,
			}
			if userData.LastName != nil {
				user.LastName = *userData.LastName
			}

			// user.ImageThumbnail = userData.ImageThumbnailURL

			if userData.Mobile != nil {
				user.Mobile = *userData.Mobile
			}

			lastLoginData := user.AppendGeoInfo()
			// err = db.Create(&user).Error
			//if err != nil {
			//	log.Printf("[Could not create User] user [%v] error [%v]\n", userLogin, err)
			//	LogDiscordError(fmt.Sprintf("[Could not create User] user [%v] error [%v]\n", userLogin, err))
			//	var serverError p2pErrors.ErrorTemporaryServerError
			//	c.JSON(http.StatusInternalServerError, serverError.JSONError())
			//	cacheDurationInSeconds := 60 //1 minute
			//
			//	s.GC.Cache.CacheHttpResponse(cacheKey, http.StatusInternalServerError, serverError.JSONError(), cacheDurationInSeconds)
			//
			//	return
			//}
			////log LastLogin
			user.LogLastLoginData(s.TrovoWalletDB, lastLoginData)
		} else {
			userData := responseData.UserInfo.UserData
			// check if user is already suspended
			if user.Suspended == 1 {
				accesslog.Record(s.AdminDB, accesslog.Event{
					Event:       models.EventLoginSuspended,
					Status:      models.AccessStatusFailed,
					Category:    models.AccessCategorySession,
					Actor:       accessActorForUsername(s, userLogin),
					IPAddress:   accessClientIP(c),
					UserAgent:   c.Request.UserAgent(),
					Method:      c.Request.Method,
					Path:        c.FullPath(),
					LoginMethod: "push_approval",
					Detail:      "login blocked: account suspended",
				})
				c.JSON(http.StatusBadRequest, gin.H{"error": "account has been suspended", "message": "account has been suspended"})
				return
			}
			// user exists. update it

			clientIP := c.ClientIP()
			if clientIP == "127.0.0.1" {
				clientIP = "41.190.2.194"
			}
			user.PublicIP = clientIP
			if userData.LastName != nil {
				user.LastName = *userData.LastName
			}

			user.FirstName = userData.FirstName

			// user.ImageThumbnail = userData.ImageThumbnailURL

			if len(userData.Email) > 0 {
				user.Email = userData.Email
			}

			user.PublicKey = userData.PublicKey

			lastLoginData := user.AppendGeoInfo()
			// err = db.Save(&user).Error
			//if err != nil {
			//	log.Printf("[Could not update User] user [%v] error [%v]\n", userLogin, err)
			//	LogDiscordError(fmt.Sprintf("[Could not update User] user [%v] error [%v]\n", userLogin, err))
			//	var serverError p2pErrors.ErrorTemporaryServerError
			//	c.JSON(http.StatusInternalServerError, serverError.JSONError())
			//	cacheDurationInSeconds := 30 //1 minute
			//
			//	s.GC.Cache.CacheHttpResponse(cacheKey, http.StatusInternalServerError, serverError.JSONError(), cacheDurationInSeconds)
			//
			//	return
			//}
			//log LastLogin
			user.LogLastLoginData(s.TrovoWalletDB, lastLoginData)
		}
		userInfo, _ := userServices.GetUserDetails(user.ID, s.TrovoWalletDB)
		log.Println("tom: ", userInfo)
		tokenResponse := struct {
			AccessToken  string      `json:"accessToken"`
			RefreshToken string      `json:"refreshToken"`
			UserInfo     interface{} `json:"userInfo"`
		}{
			AccessToken:  responseData.AccessToken,
			RefreshToken: responseData.RefreshToken,
			UserInfo:     userInfo,
		}
		log.Println("tom: ", tokenResponse)
		{
			// perform cache ops

			// {
			// 	//invalidate cache
			// 	cacheKey := fmt.Sprintf("[GET] /v1/users/verify/%v/%v", targetUser, loginID)
			// 	s.GC.Cache.InvalidateCachedHttpResponse(cacheKey)
			// }

			cacheDurationInSeconds := 2 * 60 // 2 minutes

			// The cached value must match what the live response below actually sends, both
			// in status code (http.StatusCreated, not the previously-cached 200 OK) and in
			// body shape: the live path wraps tokenResponse in serverResponse.JSON's
			// {message, data, timestamp, status} envelope, but the cached value used to be
			// the raw, unwrapped tokenResponse struct — a cache hit would have sent the
			// frontend a body without the `data` wrapper it expects (loginDetails.data.accessToken
			// would have been undefined).
			cachedBody := serverResponse.Data{
				Message:   "verify data fetched successfully",
				Data:      tokenResponse,
				Status:    http.StatusText(http.StatusCreated),
				Timestamp: time.Now().Format(time.RFC850),
			}
			s.GC.Cache.CacheHttpResponse(cacheKey, http.StatusCreated, cachedBody, cacheDurationInSeconds)

		}
		accesslog.Record(s.AdminDB, accesslog.Event{
			Event:       models.EventLoginSuccess,
			Status:      models.AccessStatusSuccessful,
			Category:    models.AccessCategorySession,
			Actor:       accessActorForUsername(s, userLogin),
			IPAddress:   accessClientIP(c),
			UserAgent:   c.Request.UserAgent(),
			Method:      c.Request.Method,
			Path:        c.FullPath(),
			LoginMethod: "push_approval",
		})
		serverResponse.JSON(c, http.StatusCreated, "verify data fetched successfully", tokenResponse, nil)
	}

}

// loginStreamHeartbeatInterval keeps an otherwise-idle SSE connection alive through
// intermediary proxies that would close it after a period of no bytes sent.
const loginStreamHeartbeatInterval = 20 * time.Second

// loginStreamMaxDuration caps how long a single stream connection is held open, so an
// abandoned login attempt (tab left open, never approved) doesn't hold a server-side
// connection forever. The frontend's existing polling keeps working as a fallback once this
// closes.
const loginStreamMaxDuration = 5 * time.Minute

// LoginNotificationStream is the subscribe side of the login push-notification flow whose
// publish side (GlobalConfig.LoginUsers / BroadcastToLoginID) already exists and already fires
// on every real approval, from LoginCallback — this handler is what was missing to actually
// deliver that notification to a browser. It streams Server-Sent Events rather than using a
// WebSocket: this is a one-directional, server-to-client notification (there is nothing for the
// client to send back), so SSE needs no new dependency and no protocol upgrade, just a
// long-lived text/event-stream HTTP response.
//
// Unauthenticated by design, exactly like GET /users/verify/:targetUser/:loginID which it
// complements — at this point in the flow the caller has no JWT yet, and EventSource cannot send
// custom headers regardless. loginID (an unguessable, single-use identifier minted when the
// login attempt starts) is the same de facto capability boundary the polling endpoint already
// relies on.
//
// @Summary Stream login approval notifications
// @Description Server-Sent Events stream that pushes a "loginNotice" event the moment BroadcastToLoginID fires for this loginID (i.e. the moment the wallet approval callback arrives). The event payload only signals that a check is worth doing now — the client still calls GET /users/verify/:targetUser/:loginID to fetch the actual tokens.
// @Tags Authentication
// @Produce text/event-stream
// @Param loginID path string true "Login ID associated with the login attempt"
// @Success 200 {string} string "text/event-stream"
// @Failure 400 {object} object
// @Router /login/stream/{loginID} [get]
func LoginNotificationStream(c *gin.Context, s *serverModels.Server) {
	loginID := strings.TrimSpace(strings.ToLower(c.Param("loginID")))
	if len(loginID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loginID is required", "message": "loginID is required"})
		return
	}

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming unsupported", "message": "streaming unsupported"})
		return
	}

	// Buffered so BroadcastToLoginID's send (which holds GlobalConfig.Mutex while it sends)
	// never blocks waiting on this handler's select loop below — a login notice fires at most
	// once per attempt, so a capacity of 1 is always enough.
	ch := make(chan map[string]interface{}, 1)

	s.GC.Mutex.Lock()
	s.GC.LoginUsers[loginID] = ch
	s.GC.Mutex.Unlock()

	defer func() {
		s.GC.Mutex.Lock()
		// Only remove our own registration — a reconnect for the same loginID may already
		// have replaced it with a newer channel by the time this deferred cleanup runs.
		if existing, ok := s.GC.LoginUsers[loginID]; ok && existing == ch {
			delete(s.GC.LoginUsers, loginID)
		}
		s.GC.Mutex.Unlock()
	}()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // disable nginx response buffering for this stream, if present

	fmt.Fprint(c.Writer, ": connected\n\n")
	flusher.Flush()

	heartbeat := time.NewTicker(loginStreamHeartbeatInterval)
	defer heartbeat.Stop()

	deadline := time.NewTimer(loginStreamMaxDuration)
	defer deadline.Stop()

	for {
		select {
		case msg := <-ch:
			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("[LoginNotificationStream] failed to marshal message for loginID [%v]: %v\n", loginID, err)
				return
			}
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			flusher.Flush()
			return // one-shot: a login notice fires once per attempt, nothing more to stream
		case <-heartbeat.C:
			fmt.Fprint(c.Writer, ": heartbeat\n\n")
			flusher.Flush()
		case <-deadline.C:
			return
		case <-c.Request.Context().Done():
			// client disconnected (tab closed, navigated away, network dropped)
			return
		}
	}
}

// RefreshToken is a function to handle refreshing of access token
// @Summary Refreshes access token
// @Description Refreshes the access token using the provided refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// // @Param gc body models.GlobalConfig true "Global configuration"
// @Param Authorization header string true "Authorization token"
// @Param refreshToken body string true "Refresh token"

// @Failure 401 {object} map[string]string "Unauthorized: refresh token expired"
// @Failure 422 {object} map[string]string "Unprocessable Entity: unable to retrieve refresh token"
// @Router /v1/login/token/refresh [post]
func RefreshToken(c *gin.Context, s *serverModels.Server) {

	// gc := models.GlobalConfig{} // check to insert correctly

	// router.POST("/v1/callbacks/login/:serviceName", func(c *gin.Context) {
	var err error

	conDB.PrintDBStats("POST /v1/login/token/refresh", s.TrovoWalletDB)
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		log.Printf("[HandleRefreshToken] error could not bind json body: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "unable to retrieve refresh token", "message": "unable to retrieve refresh token"})
		return
	}
	refreshToken := mapToken["refreshToken"]

	// verify the token
	tokenR, err := s.GC.ServiceLink.JwtTokenRefresh(refreshToken)
	// if there is an error, the token must have expired
	if err != nil {
		log.Printf("[HandleRefreshToken] err token parse: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired", "message": "refresh token expired"})
		return
	}

	// Since token is valid, get the uuid:

	userID := tokenR.UserInfo.UserData.ID

	// tokens := map[string]string{
	// 	"accessToken":  ts.AccessToken,
	// 	"refreshToken": ts.RefreshToken,
	// }
	userInfo, _ := userServices.GetUser(userID, s.TrovoWalletDB)
	tokenResponse := struct {
		AccessToken  string      `json:"accessToken"`
		RefreshToken string      `json:"refreshToken"`
		UserInfo     interface{} `json:"userInfo"`
	}{
		AccessToken:  tokenR.AccessToken,
		RefreshToken: tokenR.RefreshToken,
		UserInfo:     userInfo,
	}
	// Resolve the admin actor from the wallet user id carried by the refreshed token.
	refreshActor := accesslog.Actor{}
	{
		var admin models.AdminUser
		if err := s.AdminDB.First(&admin, "wallet_user_id = ?", userID).Error; err == nil {
			id := admin.ID
			refreshActor = accesslog.Actor{
				AdminID:  &id,
				Username: admin.Username,
				FullName: strings.TrimSpace(admin.FirstName + " " + admin.LastName),
				Email:    admin.Email,
				Role:     string(admin.Role),
			}
		}
	}
	accesslog.Record(s.AdminDB, accesslog.Event{
		Event:     models.EventTokenRefresh,
		Status:    models.AccessStatusSuccessful,
		Category:  models.AccessCategorySession,
		Actor:     refreshActor,
		IPAddress: accessClientIP(c),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.FullPath(),
	})
	c.JSON(http.StatusCreated, tokenResponse)

}

// Logout handles user logout.
// @Summary User logout
// @Description Logs out the currently authenticated user.
// @ID UserLogout
// @Tags Authentication
// @Produce json
// @Security JwtTokenAuth
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} object "Successfully logged out"
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal server error"
// @Router /logout [post]
func Logout(c *gin.Context, s *serverModels.Server) {

	conDB.PrintDBStats("/v1/logout", s.TrovoWalletDB)

	delErr := middleware.DeleteAuth(c.Request)
	if delErr != nil { // if any goes wrong

		var ex p2pErrors.GenericError
		var ok bool

		ex, ok = delErr.(p2pErrors.GenericError)
		statusCode := 0
		var response interface{}

		if ok {
			statusCode = ex.HTTPCode()
			response = ex.JSONError()
		} else {
			statusCode = http.StatusUnauthorized
			response = gin.H{"error": "unauthorized to perform this action", "message": "unauthorized to perform this action"}
		}

		c.JSON(statusCode, response)
		return

	}
	accesslog.Record(s.AdminDB, accesslog.Event{
		Event:     models.EventLogout,
		Status:    models.AccessStatusSuccessful,
		Category:  models.AccessCategorySession,
		IPAddress: accessClientIP(c),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.FullPath(),
	})
	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})

}

func LogDiscordError(msg string) {
	// Also count it, so the failure shows on the error-rate dashboard
	// instead of only in Discord. Additive: Discord is unchanged.
	observe.RecordHandledFailure(msg)

	discord.WebhookURL = "https://discord.com/api/webhooks/865931042795290636/jObHzZWdnbhX1jomOSZQX8Ip5AXLArh87PI4-ZQ8u6ssnRbZuVdY_iPxz5qoWkHUlZwS"
	if len(os.Getenv("500_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("500_ERROR_WEBHOOK")
	}
	err := discord.Say(msg)
	if err != nil {
		return
	}
}
