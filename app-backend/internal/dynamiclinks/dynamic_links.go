package dynamiclinks

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
)

type ShortLinkRequest struct {
	DynamicLinkInfo DynamicLinkInfo `json:"dynamicLinkInfo"`
	Suffix          Suffix          `json:"suffix"`
}
type DynamicLinkInfo struct {
	DomainUriPrefix string      `json:"domainUriPrefix"`
	Link            string      `json:"link"`
	AndroidInfo     AndroidInfo `json:"androidInfo"`
	IosInfo         IosInfo     `json:"iosInfo"`
}

type AndroidInfo struct {
	AndroidPackageName string `json:"androidPackageName"`
	// AndroidFallbackLink string `json:"androidFallbackLink"`
}
type IosInfo struct {
	IosBundleId string `json:"iosBundleId"`
	// IosFallbackLink string `json:"iosFallbackLink"`
}

type Suffix struct {
	Option string `json:"option"`
}

type ShortLinkResponse struct {
	ShortLink   string `json:"shortLink"`
	PreviewLink string `json:"previewLink"`
}
type ErrorResponse struct {
	Error DynamicLinkError `json:"error"`
}
type DynamicLinkError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PayWithTrovoWalletData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
}
type ReferralLinkData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
}

type TokenizedAssetDeepLinkData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
}

// FBDL is model for sending the POST request to the DL proxy
type FBDL struct {
	Link   string `json:"link"`
	APIKey string `json:"apiKey"`
}

// DynamicLink is model for saving dybamic links
type DynamicLink struct {
	ID   string `json:"linkId"`
	Link string `json:"link"`
}

type FBDLResponse struct {
	DynamicLink string `json:"dynamicLink"`
}
type LoginWithTrovoWalletData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
	LoginID     string `json:"loginId"`
}

type TrovoWalletAuthorizationData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
	AuthID      string `json:"authId"`
}
type TrovoWalletEventData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
	EventID     string `json:"eventId"`
}

func GenerateDynamicLinkWithStaticService(link string, gc *sharedconfig.GlobalConfig) (dynamicLink string, err error) {

	cacheKey := link
	{

		// search cache for link

		ok, _, response := gc.RedisCache.CachedHttpResponse(cacheKey)

		if ok {
			// log.Printf("[GenerateDynamicLinkWithStaticService][%v], served from cache\n", cacheKey)
			dynamicLink = response.(string)
			return
		}

	}
	return GenerateDynamicLink(link, gc)

}

func GenerateDynamicLinkWithStaticServiceOld(link string, dynamicLinkServiceUrl string, gc *sharedconfig.GlobalConfig) (dynamicLink string, err error) {

	cacheKey := link
	{

		// search cache for link

		ok, _, response := gc.RedisCache.CachedHttpResponse(cacheKey)

		if ok {
			// log.Printf("[GenerateDynamicLinkWithStaticService][%v], served from cache\n", cacheKey)
			dynamicLink = response.(string)
			return
		}

	}

	baseUrl := dynamicLinkServiceUrl

	if len(link) == 0 {
		err = errors.New("no link submitted for QRCode")
		return
	}
	// p := new(ShortLinkResponse)
	// e := new(ErrorResponse)

	//apiKey holds the d-link service authentication  API key
	apiKey := os.Getenv("FBDL_SERVICE_API_KEY")

	// Get service URL to use for dynamic links
	rBody := FBDL{
		Link:   link,
		APIKey: apiKey,
	}

	jbody, err := json.Marshal(rBody)

	if err != nil {
		log.Printf("[GenerateDynamicLinkWithStaticService] Unable to marshal json string Error: %v\n", err)
		return
	}
	jb := bytes.NewBuffer(jbody)
	resp, err := http.Post(baseUrl, "application/json", jb)
	if err != nil {
		log.Printf("[GenerateDynamicLinkWithStaticService] Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	log.Println("[GenerateDynamicLinkWithStaticService]>>>>>>>>>Reading Dynamics Links response Body string", string(body))

	if err != nil {
		log.Println("[GenerateDynamicLinkWithStaticService]Reading Dynamics Links response Body failed with", err)
		return
	}

	var sr FBDLResponse
	err = json.Unmarshal(body, &sr)
	if err != nil {
		log.Printf("[GenerateDynamicLinkWithStaticService] Error: %v\n", err)
		return "", err
	}

	if len(sr.DynamicLink) == 0 {
		return "", errors.New("no short link generated")
	}
	// cache the link
	gc.RedisCache.CacheHttpResponse(cacheKey, 200, sr.DynamicLink, (525960 * 3 * 60))

	return sr.DynamicLink, nil

}

func GenerateDynamicLink(link string, gc *sharedconfig.GlobalConfig) (dynamicLink string, err error) {
	// log.Println("link for: ", link)
	cacheKey := link
	{

		// search cache for link

		ok, response := gc.RedisCache.GetCachedResult(cacheKey)

		if ok {
			// log.Printf("[%v], served from cache\n", cacheKey)
			dynamicLink = response.(string)
			return
		}

	}

	if len(link) == 0 {
		log.Println("[GenerateDynamicLink] No link submitted for qrCode")
		err = errors.New("no link submitted for QRCode")
		return
	}
	shortLink, err := GetShortLink(link, gc)
	if err != nil {
		log.Println("[GenerateDynamicLink] No shortlinks Links generated:", err)

		return "", errors.New("no short link generated")
	}
	// cache the link
	gc.RedisCache.StoreResultToCache(cacheKey, shortLink, (525960 * 3 * 60))

	return shortLink, nil

}

// func GenerateDynamicLinkOld(link string, gc *sharedconfig.GlobalConfig) (dynamicLink string, err error) {
// 	// log.Println("link for: ", link)
// 	cacheKey := link
// 	{

// 		// search cache for link

// 		ok, response := gc.RedisCache.GetCachedResult(cacheKey)

// 		if ok {
// 			// log.Printf("[%v], served from cache\n", cacheKey)
// 			dynamicLink = response.(string)
// 			return
// 		}

// 	}

// 	baseUrl := <-gc.DynamicLinkServiceURLChan
// 	// log.Println("using baseurl:", baseUrl)
// 	defer func() {
// 		time.Sleep(200 * time.Millisecond) // wait for 200ms before sending next request. enough time to achieve 5 requests per ip
// 		//return the link to waiting list
// 		gc.DynamicLinkServiceURLChan <- baseUrl
// 	}()
// 	if len(link) == 0 {
// 		log.Println("[GenerateDynamicLink] No link submitted for qrCode")
// 		err = errors.New("no link submitted for QRCode")
// 		return
// 	}
// 	// p := new(ShortLinkResponse)
// 	// e := new(ErrorResponse)

// 	//apiKey holds the d-link service authentication  API key
// 	apiKey := os.Getenv("FBDL_SERVICE_API_KEY")

// 	// Get service URL to use for dynamic links
// 	rBody := FBDL{
// 		Link:   link,
// 		APIKey: apiKey,
// 	}

// 	jbody, err := json.Marshal(rBody)

// 	if err != nil {
// 		log.Printf("[GenerateDynamicLink] Unable to marshal json string Error: %v\n", err)
// 		return
// 	}
// 	jb := bytes.NewBuffer(jbody)
// 	resp, err := http.Post(baseUrl, "application/json", jb)
// 	if err != nil {
// 		log.Printf("[GenerateDynamicLink] Error: %v\n", err)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println("[GenerateDynamicLink] Reading Dynamics Links response Body failed with", err, "\nbody:", body)
// 		return
// 	}
// 	// log.Printf("[GenerateDynamicLink] Reading Dynamics Links response: %s\n", body)

// 	var sr FBDLResponse
// 	err = json.Unmarshal(body, &sr)
// 	if err != nil {
// 		log.Printf("[GenerateDynamicLink] Error: %v\nbody: %s", err, body)
// 		return "", err
// 	}
// 	// _, err = sling.New().Base(baseUrl).Post("v1/shortLinks?key="+apiKey).BodyJSON(shortLinkBody).Receive(p, e)
// 	// if err != nil {
// 	// 	log.Printf("[GenerateDynamicLink] Error: %v\n", err)
// 	// 	return "", err
// 	// }

// 	if len(sr.DynamicLink) == 0 {
// 		log.Println("[GenerateDynamicLink] No shortlinks Links generated", err, "\nbody:", body)

// 		return "", errors.New("no short link generated")
// 	}
// 	// cache the link
// 	gc.RedisCache.StoreResultToCache(cacheKey, sr.DynamicLink, (525960 * 3 * 60))

// 	return sr.DynamicLink, nil

// }

// GenerateLoginData generates Login Data
func GenerateLoginData(ownerUsername, serviceShortName, targetUser, loginID, deviceInfo, loginDescription string, gc *sharedconfig.GlobalConfig) (p LoginWithTrovoWalletData, err error) {
	if len(loginDescription) == 0 {
		loginDescription = fmt.Sprintf("This is a request to authorize a login session for the Trovo App user account %s on the service %s.", targetUser, strings.ToUpper(serviceShortName))
	}
	var dynamicLink, pngDataURI string
	params := url.Values{}
	params.Add("action", "login")
	params.Add("ownerUsername", ownerUsername)
	params.Add("serviceShortName", serviceShortName)
	params.Add("targetUser", targetUser)
	params.Add("deviceInfo", deviceInfo)
	params.Add("loginId", loginID)
	params.Add("description", loginDescription)
	link := fmt.Sprintf("%v?%v", os.Getenv("DYNAMIC_LINKS_FALLBACK_BASE_URL"), params.Encode())
	// log.Println("[GenerateLoginData]link=", link)

	dynamicLink, err = GenerateDynamicLink(link, gc)

	if err != nil {
		log.Printf("[GenerateLoginData] could not generate dynamicLink for username %s [%v]. error: %v\n", targetUser, link, err)
		return
	}
	if len(dynamicLink) == 0 {
		log.Println("[GenerateLoginData] unable to generate dynamic link=", dynamicLink)
		return
	}

	pngDataURI, err = GenerateQRCode(dynamicLink, gc)
	if err != nil {
		log.Printf("[GenerateLoginData] could not generate QRCode for [%v]. error: %v\n", dynamicLink, err)
		return
	}
	p.DynamicLink = dynamicLink
	p.QRCode = pngDataURI
	p.LoginID = loginID
	// log.Printf("[GenerateLoginData] App Data Link:[%+v]\n", p)
	return p, nil
}

// GenerateAuthorizationData generates authorization Data
func GenerateAuthorizationData(ownerUsername, serviceShortName, description, targetUser, deviceInfo, authID string, gc *sharedconfig.GlobalConfig) (p TrovoWalletAuthorizationData, err error) {
	if len(description) == 0 {
		description = fmt.Sprintf("This is a request to authorize a 2FA action for the TrovoApp user account %s", targetUser)
	}
	description = fmt.Sprintf("%s on the service %s.", description, strings.ToUpper(serviceShortName))
	var dynamicLink, pngDataURI string
	params := url.Values{}
	params.Add("action", "authorize")
	params.Add("ownerUsername", ownerUsername)
	params.Add("serviceShortName", serviceShortName)
	params.Add("targetUser", targetUser)
	params.Add("deviceInfo", deviceInfo)
	params.Add("description", description)
	params.Add("authId", authID)
	link := fmt.Sprintf("%v?%v", os.Getenv("DYNAMIC_LINKS_FALLBACK_BASE_URL"), params.Encode())
	// log.Println("[GenerateAuthorizationData]link=", link)

	dynamicLink, err = GenerateDynamicLink(link, gc)

	if err != nil {
		log.Printf("[GenerateAuthorizationData] could not generate dynamic-link for [%v]. error: %v\n", link, err)
		return
	}
	// log.Println("[GenerateAuthorizationData] generated dynamic link=", dynamicLink)
	if len(dynamicLink) == 0 {
		log.Println("[GenerateAuthorizationData] unable to generate dynamic link=", dynamicLink)
		return
	}
	pngDataURI, err = GenerateQRCode(dynamicLink, gc)
	if err != nil {
		log.Printf("[GenerateAuthorizationData] could not generate QRCode for [%v]. error: %v\n", dynamicLink, err)
		return
	}
	p.DynamicLink = dynamicLink
	p.QRCode = pngDataURI
	p.AuthID = authID
	// log.Printf("[GenerateAuthorizationData] App Data Link:[%+v]\n", p)
	return p, nil
}

// GenerateEventData generates authorization Data
func GenerateEventData(ownerUsername, serviceShortName, description, deviceInfo, eventID string, gc *sharedconfig.GlobalConfig) (p TrovoWalletEventData, err error) {
	if len(description) == 0 {
		description = fmt.Sprintf("This is a request to authorize event registration for your TrovoApp user account on the service %s.", strings.ToUpper(serviceShortName))
	}
	var dynamicLink, pngDataURI string
	params := url.Values{}
	params.Add("action", "event")
	params.Add("ownerUsername", ownerUsername)
	params.Add("serviceShortName", serviceShortName)
	params.Add("deviceInfo", deviceInfo)
	params.Add("description", description)
	params.Add("eventId", eventID)
	link := fmt.Sprintf("%v?%v", os.Getenv("DYNAMIC_LINKS_FALLBACK_BASE_URL"), params.Encode())
	// log.Println("[GenerateAuthorizationData]link=", link)

	dynamicLink, err = GenerateDynamicLink(link, gc)

	if err != nil {
		log.Printf("[GenerateEventData] could not generate dynamic-link for [%v]. error: %v\n", link, err)
		return
	}
	// log.Println("[v] generated dynamic link=", dynamicLink)
	if len(dynamicLink) == 0 {
		log.Println("[GenerateEventData] unable to generate dynamic link=", dynamicLink)
		return
	}
	pngDataURI, err = GenerateQRCode(dynamicLink, gc)
	if err != nil {
		log.Printf("[GenerateEventData] could not generate QRCode for [%v]. error: %v\n", dynamicLink, err)
		return
	}
	p.DynamicLink = dynamicLink
	p.QRCode = pngDataURI
	p.EventID = eventID
	// log.Printf("[GenerateEventData] App Data Link:[%+v]\n", p)
	return p, nil
}

// GeneratePaymentData generates payment Data
func GeneratePaymentData(paymentDestination, assetCode, assetIssuer, amount, memo string, gc *sharedconfig.GlobalConfig) (p PayWithTrovoWalletData, err error) {
	if len(paymentDestination) == 0 {
		err = errors.New("no payment destination")
		return
	}
	if len(memo) > 60 {
		err = errors.New("memo cannot be more than 60 bytes in length")
		return
	}

	if assetCode == "" || assetCode == os.Getenv("NATIVE_ASSET_CODE") {
		assetCode = os.Getenv("NATIVE_ASSET_CODE")
		assetIssuer = ""
	}

	if len(assetIssuer) > 0 && len(assetIssuer) != 42 {
		err = errors.New("invalid asset issuer")
		return
	}
	if len(assetCode) < 3 || len(assetCode) > 12 {
		err = errors.New("invalid asset code")
		return
	}
	if len(amount) > 0 {
		amount = decimal.RequireFromString(amount).Truncate(7).String()

	}
	var dynamicLink, pngDataURI string
	params := url.Values{}
	params.Add("action", "payment")
	params.Add("paymentDestination", paymentDestination)
	params.Add("assetCode", assetCode)
	params.Add("assetIssuer", assetIssuer)
	params.Add("amount", amount)
	params.Add("memo", memo)
	link := fmt.Sprintf("%v?%v", os.Getenv("DYNAMIC_LINKS_FALLBACK_BASE_URL"), params.Encode())
	// log.Println("[GeneratePaymentData]link=", link)

	dynamicLink, err = GenerateDynamicLink(link, gc)

	if err != nil {
		log.Printf("[GeneratePaymentData]could not generate dynamic-link for [%v]. error: %v\n", link, err)
		return
	}
	// log.Println("[GenerateLoginData] generated dynamic link=", dynamicLink)
	if len(dynamicLink) == 0 {
		log.Println("[GeneratePaymentData] unable to generate dynamic link=", dynamicLink)
		return
	}
	pngDataURI, err = GenerateQRCode(dynamicLink, gc)
	if err != nil {
		log.Printf("[GeneratePaymentData] could not generate QRCode for [%v]. error: %v\n", dynamicLink, err)
		return
	}
	p.DynamicLink = dynamicLink
	p.QRCode = pngDataURI
	// log.Printf("[GeneratePaymentData] App Data Link:[%+v]\n", p)
	return p, nil
}

// GenerateTokenizedAssetDeeplink generates deep link for tokenized asset Data
func GenerateTokenizedAssetDeeplink(assetCode, assetIssuer string, gc *sharedconfig.GlobalConfig) (p TokenizedAssetDeepLinkData, err error) {

	if len(assetIssuer) > 0 && len(assetIssuer) != 42 {
		err = errors.New("invalid asset issuer")
		return
	}
	if len(assetCode) < 3 || len(assetCode) > 12 {
		err = errors.New("invalid asset code")
		return
	}

	var dynamicLink, pngDataURI string
	params := url.Values{}
	params.Add("action", "tokenizedAsset")
	params.Add("assetCode", assetCode)
	params.Add("assetIssuer", assetIssuer)

	link := fmt.Sprintf("%v?%v", os.Getenv("DYNAMIC_LINKS_FALLBACK_BASE_URL"), params.Encode())
	// log.Println("[GeneratePaymentData]link=", link)

	dynamicLink, err = GenerateDynamicLink(link, gc)

	if err != nil {
		log.Printf("[GenerateTokenizedAssetDeeplink]could not generate dynamic-link for [%v]. error: %v\n", link, err)
		return
	}
	// log.Println("[GenerateLoginData] generated dynamic link=", dynamicLink)
	if len(dynamicLink) == 0 {
		log.Println("[GenerateTokenizedAssetDeeplink] unable to generate dynamic link=", dynamicLink)
		return
	}
	pngDataURI, err = GenerateQRCode(dynamicLink, gc)
	if err != nil {
		log.Printf("[GenerateTokenizedAssetDeeplink] could not generate QRCode for [%v]. error: %v\n", dynamicLink, err)
		return
	}
	p.DynamicLink = dynamicLink
	p.QRCode = pngDataURI
	// log.Printf("[GeneratePaymentData] App Data Link:[%+v]\n", p)
	return p, nil
}

// GenerateReferralLinkWithStaticURL generates payment Data
func GenerateReferralLinkWithStaticURL(username string, gc *sharedconfig.GlobalConfig) (p ReferralLinkData, err error) {
	if len(username) == 0 {
		err = errors.New("no username")
		return
	}

	var dynamicLink, pngDataURI string
	params := url.Values{}
	params.Add("action", "register")
	params.Add("referrer", username)

	link := fmt.Sprintf("%v?%v", os.Getenv("DYNAMIC_LINKS_FALLBACK_BASE_URL"), params.Encode())

	dynamicLink, err = GenerateDynamicLinkWithStaticService(link, gc)

	if err != nil {
		log.Printf("[GenerateReferralLink]could not generate dynamic-link for [%v]. error: %v\n", username, err)
		return
	}
	// log.Println("[GenerateLoginData] generated dynamic link=", dynamicLink)
	if len(dynamicLink) == 0 {
		log.Println("[GenerateReferralLink] unable to generate dynamic link=", dynamicLink, "for username=", username)
		return
	}
	pngDataURI, err = GenerateQRCode(dynamicLink, gc)
	if err != nil {
		log.Printf("[GenerateReferralLink] could not generate QRCode for [%v]. error: %v\n", dynamicLink, err)
		return
	}
	p.DynamicLink = dynamicLink
	p.QRCode = pngDataURI
	// log.Printf("[GenerateReferralLink] App Data Link:[%+v]\n", p)
	return p, nil
}

// GenerateReferralLink generates payment Data
func GenerateReferralLink(username string, gc *sharedconfig.GlobalConfig) (p ReferralLinkData, err error) {
	if len(username) == 0 {
		err = errors.New("no username")
		return
	}

	var dynamicLink, pngDataURI string
	params := url.Values{}
	params.Add("action", "register")
	params.Add("referrer", username)

	link := fmt.Sprintf("%v?%v", os.Getenv("DYNAMIC_LINKS_FALLBACK_BASE_URL"), params.Encode())

	dynamicLink, err = GenerateDynamicLink(link, gc)

	if err != nil {
		log.Printf("[GenerateReferralLink]could not generate dynamiclink for [%v]. error: %v\n", username, err)
		return
	}
	// log.Println("[GenerateLoginData] generated dynamic link=", dynamicLink)
	if len(dynamicLink) == 0 {
		log.Println("[GenerateReferralLink] unable to generate dynamic link=", dynamicLink, "for username=", username)
		return
	}
	pngDataURI, err = GenerateQRCode(dynamicLink, gc)
	if err != nil {
		log.Printf("[GenerateReferralLink] could not generate QRCode for [%v]. error: %v\n", dynamicLink, err)
		return
	}
	p.DynamicLink = dynamicLink
	p.QRCode = pngDataURI
	// log.Printf("[GenerateReferralLink] App Data Link:[%+v]\n", p)
	return p, nil
}

// GetShortLink saves link to DB and generates shortlink
func GetShortLink(link string, gc *sharedconfig.GlobalConfig) (shortLink string, err error) {
	if len(link) == 0 {
		err = errors.New("no link")
		return
	}
	shortLinkBaseURL := os.Getenv("SHORT_LINKS_BASE_URL")
	if len(shortLinkBaseURL) == 0 {
		shortLinkBaseURL = "https://links.trovo.app"
	}
	linkID := GenerateRandomCode(18)
	dblink := DynamicLink{
		ID:   linkID,
		Link: link,
	}
	e := gc.DB.Create(&dblink).Error
	if e != nil {
		log.Printf("[GetShortLink]1st try: could not save shortlink to db. error: %v\n", e)

		linkID = GenerateRandomCode(16)
		dblink = DynamicLink{
			ID:   linkID,
			Link: link,
		}
		e := gc.DB.Create(&dblink).Error
		if e != nil {
			log.Printf("[GetShortLink]2nd Try: could not save shortlink to db. error: %v\n", e)

			linkID = GenerateRandomCode(16)
			dblink = DynamicLink{
				ID:   linkID,
				Link: link,
			}
			e := gc.DB.Create(&dblink).Error
			if e != nil {
				log.Printf("[GetShortLink]3rd try: could not save shortlink to db. error: %v\n", e)

				return "", errors.New("failed to create shortlink")
			}
		}

	}

	shortLink = fmt.Sprintf("%v/%v", shortLinkBaseURL, linkID)

	return shortLink, nil

}

// GenerateRandomCode generates email verification code
func GenerateRandomCode(codeLength int) string {
	if codeLength == 0 {
		codeLength = 14
	}

	// Create random bytes using crypto/rand
	randomBytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return ""
	}

	// Hash the random bytes with SHA-256
	hash := sha256.Sum256(randomBytes)
	hashStr := hex.EncodeToString(hash[:]) // hex chars are already 0-9 and a-f

	// Allowed characters (a-z, 0-9)
	allowedChars := "abcdefghijklmnopqrstuvwxyz0123456789"

	// Convert the hash to allowed characters
	var filtered strings.Builder
	for _, ch := range hashStr {
		// map a-f to some letters for full coverage
		if ch >= 'a' && ch <= 'f' {
			// shift to random letters from a-z
			filtered.WriteByte(byte('a' + (ch-'a')%26))
		} else {
			filtered.WriteRune(ch)
		}
	}

	// Build the final secure code
	result := make([]byte, codeLength)
	for i := 0; i < codeLength; i++ {
		index, err := rand.Int(rand.Reader, bigInt(len(allowedChars)))
		if err != nil {
			return ""
		}
		result[i] = allowedChars[index.Int64()]
	}

	return string(result)

}

// Helper to convert int to big.Int
func bigInt(n int) *big.Int {
	return new(big.Int).SetInt64(int64(n))
}
