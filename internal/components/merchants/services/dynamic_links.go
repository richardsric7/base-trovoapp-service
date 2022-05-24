package merchants

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
	"trovo-wallet-api/internal/cache"

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

type PayWithBantupayData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
}
type RefferalLinkData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
}

//FBDL is model for sending the POST request to the DL proxy
type FBDL struct {
	Link   string `json:"link"`
	APIKey string `json:"apiKey"`
}

type FBDLResponse struct {
	DynamicLink string `json:"dynamicLink"`
}
type LoginWithBantupayData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
	LoginID     string `json:"loginId"`
}

type BantupayAuthorizationData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
	AuthID      string `json:"authId"`
}

func GenerateDynamicLinkWithStaticService(link string, dynamicLinkServiceUrl string, redisCache *cache.RedisCache) (dynamicLink string, err error) {

	cacheKey := link
	{

		// search cache for link

		ok, _, response := redisCache.CachedHttpResponse(cacheKey)

		if ok {
			log.Printf("[%v], served from cache\n", cacheKey)
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

	//apiKey holds the dlink service authentication  API key
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
	body, err := ioutil.ReadAll(resp.Body)
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
	redisCache.CacheHttpResponse(cacheKey, 200, sr.DynamicLink, (525960 * 3 * 60))

	return sr.DynamicLink, nil

}

func GenerateDynamicLink(link string, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (dynamicLink string, err error) {

	cacheKey := link
	{

		// search cache for link

		ok, _, response := redisCache.CachedHttpResponse(cacheKey)

		if ok {
			log.Printf("[%v], served from cache\n", cacheKey)
			dynamicLink = response.(string)
			return
		}

	}

	baseUrl := <-dynamicLinkServiceUrlChan
	defer func() {
		time.Sleep(200 * time.Millisecond) // wait for 200ms before sending next request. enough time to achieve 5 requests per ip
		//return the link to waiting list
		dynamicLinkServiceUrlChan <- baseUrl
	}()
	if len(link) == 0 {
		err = errors.New("no link submitted for QRCode")
		return
	}
	// p := new(ShortLinkResponse)
	// e := new(ErrorResponse)

	//apiKey holds the dlink service authentication  API key
	apiKey := os.Getenv("FBDL_SERVICE_API_KEY")

	// Get service URL to use for dynamic links
	rBody := FBDL{
		Link:   link,
		APIKey: apiKey,
	}

	jbody, err := json.Marshal(rBody)

	if err != nil {
		log.Printf("[GenerateDynamicLink] Unable to marshal json string Error: %v\n", err)
		return
	}
	jb := bytes.NewBuffer(jbody)
	resp, err := http.Post(baseUrl, "application/json", jb)
	if err != nil {
		log.Printf("[GenerateDynamicLink] Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[GenerateDynamicLink]Reading Dynamics Links response Body failed with", err)
		return
	}

	var sr FBDLResponse
	err = json.Unmarshal(body, &sr)
	if err != nil {
		log.Printf("[GenerateDynamicLink] Error: %v\n", err)
		return "", err
	}
	// _, err = sling.New().Base(baseUrl).Post("v1/shortLinks?key="+apiKey).BodyJSON(shortLinkBody).Receive(p, e)
	// if err != nil {
	// 	log.Printf("[GenerateDynamicLink] Error: %v\n", err)
	// 	return "", err
	// }

	if len(sr.DynamicLink) == 0 {
		return "", errors.New("no short link generated")
	}
	// cache the link
	redisCache.CacheHttpResponse(cacheKey, 200, sr.DynamicLink, (525960 * 3 * 60))

	return sr.DynamicLink, nil

	// if len(p.ShortLink) == 0 {
	// 	return "", errors.New("no short link generated")
	// }
	// return p.ShortLink, nil
}

//GenerateLoginData generates Login Data
func GenerateLoginData(merchant, merchantShortName, targetUser, loginID, deviceInfo string, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (p LoginWithBantupayData, err error) {

	var dynamicLink, pngDataURI string
	params := url.Values{}
	params.Add("action", "login")
	params.Add("merchant", merchant)
	params.Add("merchantShortName", merchantShortName)
	params.Add("targetUser", targetUser)
	params.Add("deviceInfo", deviceInfo)
	params.Add("loginId", loginID)
	link := fmt.Sprintf("%v?%v", os.Getenv("DYNAMIC_LINKS_FALLBACK_BASE_URL"), params.Encode())
	// log.Println("[GenerateLoginData]link=", link)

	dynamicLink, err = GenerateDynamicLink(link, dynamicLinkServiceUrlChan, redisCache)

	if err != nil {
		log.Printf("[GenerateLoginData]could not generate dynamiclink for [%v]. error: %v\n", link, err)
		return
	}
	if len(dynamicLink) == 0 {
		log.Println("[GenerateLoginData] unable to generate dynamic link=", dynamicLink)
		return
	}

	pngDataURI, err = GenerateQRCode(dynamicLink)
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

//GenerateAuthorizationData generates authorization Data
func GenerateAuthorizationData(merchant, merchantShortName, description, targetUser, deviceInfo, authID string, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (p BantupayAuthorizationData, err error) {

	var dynamicLink, pngDataURI string
	params := url.Values{}
	params.Add("action", "authorize")
	params.Add("merchant", merchant)
	params.Add("merchantShortName", merchantShortName)
	params.Add("targetUser", targetUser)
	params.Add("deviceInfo", deviceInfo)
	params.Add("description", description)
	params.Add("authId", authID)
	link := fmt.Sprintf("https://bantupay.org?%v", params.Encode())
	// log.Println("[GenerateAuthorizationData]link=", link)

	dynamicLink, err = GenerateDynamicLink(link, dynamicLinkServiceUrlChan, redisCache)

	if err != nil {
		log.Printf("[GenerateAuthorizationData]could not generate dynamiclink for [%v]. error: %v\n", link, err)
		return
	}
	// log.Println("[GenerateAuthorizationData] generated dynamic link=", dynamicLink)
	if len(dynamicLink) == 0 {
		log.Println("[GenerateAuthorizationData] unable to generate dynamic link=", dynamicLink)
		return
	}
	pngDataURI, err = GenerateQRCode(dynamicLink)
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

//GeneratePaymentData generates payment Data
func GeneratePaymentData(paymentDestination, assetCode, assetIssuer, amount, memo string, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (p PayWithBantupayData, err error) {
	if len(paymentDestination) == 0 {
		err = errors.New("no payment destination")
		return
	}
	if len(memo) > 28 {
		err = errors.New("memo cannot be more than 28 bytes in length")
		return
	}

	if assetCode == "" || assetCode == "XBN" {
		assetCode = "XBN"
		assetIssuer = ""
	}

	if len(assetIssuer) > 0 && len(assetIssuer) != 56 {
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
	link := fmt.Sprintf("https://bantupay.org?%v", params.Encode())
	// log.Println("[GeneratePaymentData]link=", link)

	dynamicLink, err = GenerateDynamicLink(link, dynamicLinkServiceUrlChan, redisCache)

	if err != nil {
		log.Printf("[GeneratePaymentData]could not generate dynamiclink for [%v]. error: %v\n", link, err)
		return
	}
	// log.Println("[GenerateLoginData] generated dynamic link=", dynamicLink)
	if len(dynamicLink) == 0 {
		log.Println("[GeneratePaymentData] unable to generate dynamic link=", dynamicLink)
		return
	}
	pngDataURI, err = GenerateQRCode(dynamicLink)
	if err != nil {
		log.Printf("[GeneratePaymentData] could not generate QRCode for [%v]. error: %v\n", dynamicLink, err)
		return
	}
	p.DynamicLink = dynamicLink
	p.QRCode = pngDataURI
	// log.Printf("[GeneratePaymentData] App Data Link:[%+v]\n", p)
	return p, nil
}

//GenerateReferralLinkWithStaticURL generates payment Data
func GenerateReferralLinkWithStaticURL(username string, dynamicLinkServiceUrl string, redisCache *cache.RedisCache) (p RefferalLinkData, err error) {
	if len(username) == 0 {
		err = errors.New("no username")
		return
	}

	var dynamicLink, pngDataURI string
	params := url.Values{}
	params.Add("action", "register")
	params.Add("referrer", username)

	link := fmt.Sprintf("https://bantupay.org?%v", params.Encode())

	dynamicLink, err = GenerateDynamicLinkWithStaticService(link, dynamicLinkServiceUrl, redisCache)

	if err != nil {
		log.Printf("[GenerateReferralLink]could not generate dynamiclink for [%v]. error: %v\n", username, err)
		return
	}
	// log.Println("[GenerateLoginData] generated dynamic link=", dynamicLink)
	if len(dynamicLink) == 0 {
		log.Println("[GenerateReferralLink] unable to generate dynamic link=", dynamicLink, "for username=", username)
		return
	}
	pngDataURI, err = GenerateQRCode(dynamicLink)
	if err != nil {
		log.Printf("[GenerateReferralLink] could not generate QRCode for [%v]. error: %v\n", dynamicLink, err)
		return
	}
	p.DynamicLink = dynamicLink
	p.QRCode = pngDataURI
	// log.Printf("[GenerateReferralLink] App Data Link:[%+v]\n", p)
	return p, nil
}

//GenerateReferralLink generates payment Data
func GenerateReferralLink(username string, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (p RefferalLinkData, err error) {
	if len(username) == 0 {
		err = errors.New("no username")
		return
	}

	var dynamicLink, pngDataURI string
	params := url.Values{}
	params.Add("action", "register")
	params.Add("referrer", username)

	link := fmt.Sprintf("https://bantupay.org?%v", params.Encode())

	dynamicLink, err = GenerateDynamicLink(link, dynamicLinkServiceUrlChan, redisCache)

	if err != nil {
		log.Printf("[GenerateReferralLink]could not generate dynamiclink for [%v]. error: %v\n", username, err)
		return
	}
	// log.Println("[GenerateLoginData] generated dynamic link=", dynamicLink)
	if len(dynamicLink) == 0 {
		log.Println("[GenerateReferralLink] unable to generate dynamic link=", dynamicLink, "for username=", username)
		return
	}
	pngDataURI, err = GenerateQRCode(dynamicLink)
	if err != nil {
		log.Printf("[GenerateReferralLink] could not generate QRCode for [%v]. error: %v\n", dynamicLink, err)
		return
	}
	p.DynamicLink = dynamicLink
	p.QRCode = pngDataURI
	// log.Printf("[GenerateReferralLink] App Data Link:[%+v]\n", p)
	return p, nil
}
