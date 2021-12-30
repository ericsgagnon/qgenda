package qgenda

import (
	// "context"
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io/ioutil"
	// "log"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	// "os"
	// "path"
	"time"

	"golang.org/x/net/context/ctxhttp"
	// "golang.org/x/net/context/ctxhttp"
)

// AuthToken holds the authorization: bearer token header and an expire timestamp
type AuthToken struct {
	Token         *http.Header  `json:"token"`
	Timestamp     time.Time     `json:"timeStamp"` // defaults to response's date header
	ValidDuration time.Duration `json:"validDuration"`
	Expires       time.Time     `json:"expires"`
	cache         *CacheFile    `json:"-"`
}

// NewAuthRequest creates a Request for authenticating the client with the qgenda
// api. The url.Values should have email and password keys with valid credentials
// see the login section of https://restapi.qgenda.com for more details.
func NewAuthRequest(u *url.Values) (*Request, error) {
	r := NewRequest()
	r.AppendPath("/login")
	r.Method = http.MethodPost
	r.Body = []byte(u.Encode())
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r, nil
}

// ProcessAuthResponse is the primary function that returns new valid AuthTokens
// from an http.Response. To minimize scope, it does not setup a cache.
func ProcessAuthResponse(r *http.Response) (*AuthToken, error) {
	// authorization token is returned in the response body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading AuthToken from response: %v", err)
		return nil, err
	}
	a := &AuthToken{
		Token: &http.Header{},
	}
	var tokenMap map[string]string
	if err := json.Unmarshal(body, &tokenMap); err != nil {
		log.Printf("Error unmarshalling AuthToken from response: %v", err)
		return nil, err
	}
	// set authorization header in AuthToken
	a.Token.Set(
		http.CanonicalHeaderKey("Authorization"),
		fmt.Sprintf("bearer %v", tokenMap["access_token"]),
	)

	// "RFC 7231"??
	timestamp, err := http.ParseTime(r.Header.Get("date"))
	if err != nil {
		return nil, err
	}
	a.Timestamp = timestamp

	// use response timestamp + valid duration to set expire time
	validDuration, err := time.ParseDuration(tokenMap["expires_in"] + "s")
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

// CacheFile creates a cache file at user/cache/dir/appName/authtoken.json
// and writes the AuthToken in json to it
func (a *AuthToken) CacheFile(cfg *CacheConfig) error {
	if cfg == nil {
		return errors.New("*CacheConfig is required.")
	}
	a.cache = NewCacheFile("authtoken.json", "auth", cfg)
	a.cache.Created = a.Timestamp
	a.cache.ValidDuration = a.ValidDuration
	a.cache.Expires = a.Expires
	a.cache.Enable = true
	if err := a.cache.Create(); err != nil {
		return err
	}
	j, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}
	return a.cache.Write(j)
}

// ReadCacheFile reads it from file and returns it if it is valid or nil if it isn't
func (a *AuthToken) ReadCacheFile(ctx context.Context) (*AuthToken, error) {
	var b []byte
	cache := *a.cache
	if err := ReadCacheFile(&cache, b); err != nil {
		return nil, err
	}
	tkn := &AuthToken{
		Token: &http.Header{},
	}
	if err := json.Unmarshal(b, tkn); err != nil {
		log.Printf("Error unmarshalling cached AuthToken: %v", err)
		return nil, err
	}
	tkn.cache = &cache
	if tkn.Valid() {
		a = tkn
	}
	return a, nil
}

func (a *AuthToken) WriteCacheFile() error {
	if !a.cache.Enable {
		return errors.New("cache disabled")
	}

	b, err := json.MarshalIndent(a, "", "\t") // make it pretty
	if err != nil {
		return err
	}
	return WriteCacheFile(a.cache, b)
}

// //
// 	//  ill expire more than a minute from now
// 	switch {
// 	case t == nil:
// 		log.Printf("AuthToken not initialized")
// 		return false
// 	case t.Token.Get(http.CanonicalHeaderKey("Authorization")) == "":
// 		log.Printf("AuthToken Authorization header is empty")
// 		fallthrough
// 	case t.Expires.IsZero():
// 		log.Printf("AuthToken expiration time is at zero %v", t.Expires.UTC().String())
// 		return false
// 	case t.Expires.UTC().Before(time.Now().Add(time.Minute).UTC()):
// 		log.Printf("AuthToken expired %v", t.Expires.UTC().String())
// 		return false
// 	case t.Expires.UTC().After(time.Now().Add(time.Minute).UTC()):
// 		log.Printf("AuthToken appears valid until %v", t.Expires.UTC().String())
// 		//TODO: Add jwt format validation
// 		return true
// 	}
// 	return false
// }

// WriteCache writes the AuthToken to a file cache
func (t *AuthToken) WriteCache(ctx context.Context, filename string) error {
	// get the file to write cache to
	filename, err := t.CacheFile(ctx, filename)
	if err != nil {
		return err
	}
	// marshal AuthToken to json
	j, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		log.Printf("Error marshalling AuthToken to json: %v", err)
	}
	log.Printf("Write AuthToken to file %v", filename)
	// write file to file cache
	if err := ioutil.WriteFile(filename, j, 0755); err != nil {
		log.Printf("Error writing file %v to disk: %v", filename, err)
		return err
	}
	// confirm success in logs
	log.Printf("AuthToken written to file %v", filename)
	return nil
}

// ReadCache reads the AuthToken from a file cache
func (t *AuthToken) ReadCache(ctx context.Context, filename string) error {

	log.Printf("Read cached AuthToken from file")
	// get the file to write cache to
	filename, err := t.CacheFile(ctx, filename)
	if err != nil {
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening AuthToken cache file %v: %v\n", filename, err)
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Error reading AuthToken cache file %v: %v\n", filename, err)
		return err
	}

	tkn := &AuthToken{
		Token:   &http.Header{},
		Expires: time.Time{},
	}
	if err := json.Unmarshal(b, tkn); err != nil {
		log.Printf("Error unmarshalling cached AuthToken: %v", err)
	}

	if tkn.Expires.UTC().Before(t.Expires.UTC()) || tkn.Expires.UTC().Before(time.Now().UTC()) {
		m := fmt.Sprintf("Cached AuthToken expired %v", tkn.Expires.UTC().String())
		log.Printf(m)
		return errors.New(m)
	}
	t.Token = tkn.Token
	t.Expires = tkn.Expires
	t = tkn

	log.Printf("Cached AuthToken appears valid until %v", t.Expires.UTC().String())
	return nil

}

// Login submits credentials for authorization bearer token
func (q *Client) Login(ctx context.Context) error {

	// request URL
	route := "/login"
	endpoint := *q.BaseURL
	endpoint.Path = path.Join(endpoint.Path, route)

	// manually set a 'start' time for the AuthToken
	// to be conservative with expiration time, the start timestamp is
	// set just before the token is requested
	startTime := time.Now().UTC()

	// res, err := q.Client.PostForm(url.String(), *q.Credentials)
	res, err := ctxhttp.PostForm(ctx, q.Client, endpoint.String(), *q.Credentials)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Printf("Logged In")

	//response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading AuthToken from response: %v", err)
		return err
	}

	// authorization token is returned in the response body
	var tokenMap map[string]string
	if err := json.Unmarshal(body, &tokenMap); err != nil {
		log.Printf("Error unmarshalling AuthToken from response: %v", err)
		return err
	}

	// Set the Authorization header in the QgendaClient
	q.Authorization.Token.Set(
		http.CanonicalHeaderKey("Authorization"),
		fmt.Sprintf("bearer %v", tokenMap["access_token"]),
	)

	// fmt.Printf("\nin Login:\n\nAuthToken:\n%v\n\n", q.Authorization.Token)
	// use response startTime + valid duration to set expire time
	validDuration, err := time.ParseDuration(tokenMap["expires_in"] + "s")
	if err != nil {
		log.Printf("Error parsing token valid duration: %v", err)
		return err
	}
	// q.Authorization.Expires = resTime.Add(validDuration)
	q.Authorization.Expires = startTime.Add(validDuration)
	log.Printf("AuthToken updated - expiration: %v", q.Authorization.Expires.UTC().Format(time.RFC3339))

	if q.Authorization.Valid(ctx) {
		log.Printf("Token appears valid")
	}

	// // write AuthToken to file
	// if err := q.Authorization.WriteFile(ctx); err != nil {
	// 	log.Printf("Unable to write AuthToken to %v:  %v", credentialsFile, err)
	// 	return err
	// }
	// log.Printf("AuthToken written to %v", credentialsFile)

	// write AuthToken to file
	if err := q.Authorization.WriteCache(ctx, ""); err != nil {
		log.Printf("Troubles using new cache file method %v", err)
		return err
	}

	return nil
}

func (q *Client) xLogin(ctx context.Context) error {

	// request URL
	route := "/login"
	endpoint := *q.BaseURL
	endpoint.Path = path.Join(endpoint.Path, route)

	// manually set a 'start' time for the AuthToken
	// to be conservative with expiration time, the start timestamp is
	// set just before the token is requested
	startTime := time.Now().UTC()

	// res, err := q.Client.PostForm(url.String(), *q.Credentials)
	res, err := ctxhttp.PostForm(ctx, q.Client, endpoint.String(), *q.Credentials)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Printf("Logged In")

	//response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading AuthToken from response: %v", err)
		return err
	}

	// authorization token is returned in the response body
	var tokenMap map[string]string
	if err := json.Unmarshal(body, &tokenMap); err != nil {
		log.Printf("Error unmarshalling AuthToken from response: %v", err)
		return err
	}

	// Set the Authorization header in the QgendaClient
	q.Authorization.Token.Set(
		http.CanonicalHeaderKey("Authorization"),
		fmt.Sprintf("bearer %v", tokenMap["access_token"]),
	)

	// fmt.Printf("\nin Login:\n\nAuthToken:\n%v\n\n", q.Authorization.Token)
	// use response startTime + valid duration to set expire time
	validDuration, err := time.ParseDuration(tokenMap["expires_in"] + "s")
	if err != nil {
		log.Printf("Error parsing token valid duration: %v", err)
		return err
	}
	// q.Authorization.Expires = resTime.Add(validDuration)
	q.Authorization.Expires = startTime.Add(validDuration)
	log.Printf("AuthToken updated - expiration: %v", q.Authorization.Expires.UTC().Format(time.RFC3339))

	if q.Authorization.Valid(ctx) {
		log.Printf("Token appears valid")
	}

	// // write AuthToken to file
	// if err := q.Authorization.WriteFile(ctx); err != nil {
	// 	log.Printf("Unable to write AuthToken to %v:  %v", credentialsFile, err)
	// 	return err
	// }
	// log.Printf("AuthToken written to %v", credentialsFile)

	// write AuthToken to file
	if err := q.Authorization.WriteCache(ctx, ""); err != nil {
		log.Printf("Troubles using new cache file method %v", err)
		return err
	}

	return nil
}

// Auth is the top level authorization method, it checks the validity of an existing token
// reads a cached token from disk into AuthToken, or logs in and retrieves a new token
// as needed
// func (q *Client) Auth(ctx context.Context) error {

// 	log.Printf("----------------------------------------------------------------------------")
// 	if q.Authorization.Valid(ctx) {
// 		return nil
// 	}

// 	log.Printf("----------------------------------------------------------------------------")
// 	err := q.Authorization.ReadCache(ctx, "")
// 	if err == nil {
// 		if q.Authorization.Valid(ctx) {
// 			return nil
// 		}
// 	}
// 	// err = q.Authorization.ReadFile(ctx)
// 	// if err == nil {
// 	// 	if q.Authorization.Valid(ctx) {
// 	// 		return nil
// 	// 	}
// 	// }

// 	log.Printf("----------------------------------------------------------------------------")
// 	log.Printf("Login")
// 	err = q.Login(ctx)
// 	return err
// }
