package qgenda

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"time"
)

// AuthToken holds the authorization: bearer token header and an expire timestamp
type AuthToken struct {
	Token         *http.Header  `json:"token"`
	Timestamp     time.Time     `json:"timeStamp"` // defaults to response's date header
	ValidDuration time.Duration `json:"validDuration"`
	Expires       time.Time     `json:"expires"`
	Cache         *CacheFile    `json:"-"`
}

// NewAuthToken returns an blank AuthToken with a configured CacheFile
func NewAuthToken(cfg *CacheConfig) (*AuthToken, error) {
	cache, err := NewAuthTokenCacheFile(cfg)
	if err != nil {
		return nil, err
	}
	tkn := &AuthToken{
		Token: &http.Header{},
		Cache: cache,
	}
	return tkn, nil
}

// NewAuthRequest creates a Request for authenticating the client with the qgenda
// api. The url.Values should have email and password keys with valid credentials
// see the login section of https://restapi.qgenda.com for more details.
func NewAuthRequest(u *url.Values) (*Request, error) {
	rc := NewRequestConfig("login", nil)
	fmt.Printf("%#v\n", rc)
	r := NewRequest(rc)
	r.Method = http.MethodPost
	r.Body = []byte(u.Encode())
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return r, nil
}

// ReadAuthTokenCacheFile
func ReadAuthTokenCacheFile(cf *CacheFile) (*AuthToken, error) {
	b, err := ReadCacheFile(cf)
	if err != nil {
		return nil, err
	}
	tkn := &AuthToken{
		Token: &http.Header{},
	}
	if err := json.Unmarshal(b, tkn); err != nil {
		log.Printf("Error unmarshalling cached AuthToken: %v", err)
		return nil, err
	}

	cache := *cf
	tkn.Cache = &cache
	return tkn, nil
}

func WriteAuthTokenCacheFile(a *AuthToken, cf *CacheFile) error {
	if !a.Cache.Enable || !a.Valid() {
		return errors.New("cache disabled")
	}

	b, err := json.MarshalIndent(a, "", "\t") // make it pretty
	if err != nil {
		return err
	}
	return WriteCacheFile(a.Cache, b)
}

func AuthTokenFromCacheFile(cf *CacheFile) (*AuthToken, error) {
	tkn, err := ReadAuthTokenCacheFile(cf)
	if err != nil {
		return nil, err
	}
	if !tkn.Valid() {
		return nil, ErrNope
	}
	return tkn, nil
}

// AuthTokenFromResponse is the primary function that returns new valid AuthTokens
// from an http.Response. To minimize scope, it does not setup a cache.
func AuthTokenFromResponse(r *http.Response) (*AuthToken, error) {
	// authorization token is returned in the response body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading AuthToken from response: %v", err)
		return nil, err
	}
	a := &AuthToken{
		Token: &http.Header{},
	}
	var tokenMap map[string]any
	if err := json.Unmarshal(body, &tokenMap); err != nil {
		log.Printf("Error unmarshalling AuthToken from response: %v", err)
		return nil, err
	}
	// fmt.Println(tokenMap)
	// set authorization header in AuthToken
	a.Token.Set(
		http.CanonicalHeaderKey("Authorization"),
		fmt.Sprintf("bearer %v", tokenMap["access_token"]),
	)

	timestamp, err := http.ParseTime(r.Header.Get("date")) // "RFC 7231"??
	if err != nil {
		return nil, err
	}
	a.Timestamp = timestamp

	// use response timestamp + valid duration to set expire time
	validDuration, err := time.ParseDuration(fmt.Sprintf("%v", tokenMap["expires_in"]) + "s")
	if err != nil {
		log.Printf("Error parsing token valid duration: %v", err)
		return nil, err
	}
	a.Expires = a.Timestamp.Add(validDuration)
	log.Printf("AuthToken updated - expiration: %v", a.Expires.UTC().Format(time.RFC3339))

	return a, nil
}

// AddAuthToken replaces any headers in the http.Request with those from AuthToken.Token
// It does not affect any headers in http.Request that aren't in AuthToken.Token
func AddAuthToken(r *http.Request, a *AuthToken) *http.Request {
	for k, v := range *a.Token {
		r.Header[k] = v
	}
	return r
}

// Valid does a simple check if AuthToken has expired
func (a *AuthToken) Valid() bool {
	now := time.Now().UTC()
	switch {
	case now.After(a.Expires):
		return false
	default:
		return true
	}
}

// NewCacheFile creates a cache file at user/cache/dir/appName/authtoken.json
// and writes the AuthToken in json to it
func (a *AuthToken) NewCacheFile(cfg *CacheConfig) error {
	if cfg == nil {
		return errors.New("*CacheConfig is required.")
	}

	cache, err := NewAuthTokenCacheFile(cfg)
	if err != nil {
		return err
	}

	cache.Enable = true
	if err := cache.Create(); err != nil {
		return err
	}
	j, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}
	err = cache.Write(j)
	cache.Expires = a.Expires // cache file can't outlive token validity
	return err
}

// ReadCacheFile reads AuthToken from file and returns it if it is valid or nil if it isn't
func (a *AuthToken) ReadCacheFile() (*AuthToken, error) {
	cache := *a.Cache
	b, err := ReadCacheFile(&cache)
	if err != nil {
		return nil, err
	}
	tkn := &AuthToken{
		Token: &http.Header{},
	}
	if err := json.Unmarshal(b, tkn); err != nil {
		log.Printf("Error unmarshalling cached AuthToken: %v", err)
		return nil, err
	}
	tkn.Cache = &cache
	if tkn.Valid() {
		a = tkn
	}
	return a, nil
}

func (a *AuthToken) WriteCacheFile() error {
	// fmt.Println(a.Valid())
	// fmt.Println(a.Cache)
	if !a.Cache.Enable || !a.Valid() {
		// if !a.Valid() {
		return errors.New("cache disabled")
	}

	b, err := json.MarshalIndent(a, "", "\t") // make it pretty
	if err != nil {
		return err
	}
	return WriteCacheFile(a.Cache, b)
}

func NewAuthTokenCacheFile(cfg *CacheConfig) (*CacheFile, error) {
	cf, err := NewCacheFile("authtoken.json", "auth", cfg)
	// fmt.Printf("NewAuthTokenCacheFile: %s\n", cf)
	if err != nil {
		return nil, err
	}
	// if err := cf.Create(); err != nil {
	// 	return nil, err
	// }
	return cf, nil
}
