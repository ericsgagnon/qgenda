package qgenda

import (
	// "context"

	// "io/ioutil"
	// "log"

	"context"
	"log"
	"net/http"
	"net/url"

	// "path"
	// "reflect"
	"time"
)

// ClientConfig is used to pass all necessary
// arguments to NewClient
type ClientConfig struct {
	URL            string
	ClientTimeout  time.Duration
	RequestTimeout time.Duration
	Email          string
	CompanyKey     string
	Password       string
	CacheConfig    *CacheConfig
}

// Client is the primary struct for handling client
// interactions with the qgenda api
type Client struct {
	URL          *url.URL
	Client       *http.Client
	Credentials  *url.Values
	AuthToken    *AuthToken
	CacheConfig  *CacheConfig
	ClientConfig *ClientConfig
	Parameters   *Parameters
}

func NewClient(cc *ClientConfig) (*Client, error) {
	// var cfg *ClientConfig
	// cfg = cc
	// fmt.Println(*cfg)

	urlString := "https://api.qgenda.com/v2"
	// if cc.URL != "" {
	// 	urlString = cc.URL
	// }
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	cch, err := NewCacheConfig("qgenda")
	// fmt.Printf("NewClient: %#v\n", cch)
	if err != nil {
		return nil, err
	}
	if cc.CacheConfig != nil {
		cache := *cc.CacheConfig
		cch = &cache
	}

	cl := &http.Client{}
	// provide reasonable default client timeout
	if time.Duration(cc.ClientTimeout) < time.Second*1 {
		cl.Timeout = time.Second * 30
	} else {
		cl.Timeout = cc.ClientTimeout
	}

	cr := &url.Values{}
	cr.Add("email", cc.Email)
	cr.Add("password", cc.Password)
	tkn, err := NewAuthToken(cch)
	// fmt.Println(tkn.Cache)
	if err != nil {
		return nil, err
	}

	cfg := *cc
	p := Parameters{
		"CompanyKey": cc.CompanyKey,
	}

	client := &Client{
		URL:          u,
		Client:       cl,
		Credentials:  cr,
		AuthToken:    tkn,
		CacheConfig:  cch,
		ClientConfig: &cfg,
		Parameters:   &p,
	}
	return client, nil
}

func (c *Client) Auth() error {
	// fmt.Println(c.AuthToken.Cache)
	if c.AuthToken.Valid() {
		log.Println("Client.AuthToken is valid")
		return nil
	}
	tkn, err := AuthTokenFromCacheFile(c.AuthToken.Cache)
	// if err == nil {
	// 	fmt.Println("No issues with the cachefile")
	// }
	if err != nil {
		log.Println("Client.Auth(): No Valid cache - requesting new AuthToken")
		// fmt.Println("CacheFile didn't work so good")
		// log.Printf("(c *Client) Auth(): %s\n", err)
		atreq, err := NewAuthRequest(c.Credentials)
		if err != nil {
			return err
		}
		req := atreq.ToHTTPRequest()
		resp, err := c.Client.Do(req)
		if err != nil {
			return err
		}
		tkn, err = AuthTokenFromResponse(resp)
		tkn.Cache = c.AuthToken.Cache
		// fmt.Printf("Client.Auth: Client.AuthToken.Cache: %s\n", tkn.Cache)
		if err != nil {
			return err
		}
		c.AuthToken = tkn
		err = c.AuthToken.WriteCacheFile()

		if err != nil {
			return err
		}

	}
	c.AuthToken = tkn
	log.Printf("Client.Auth(): AuthToken cache is valid - expires: %s\n", c.AuthToken.Expires)
	return nil
}

func (c *Client) Do(ctx context.Context, r *Request) (*http.Response, error) {

	r.SetCompanyKey(c.ClientConfig.CompanyKey)
	r.SetCompanyKey("8c44c075-d894-4b00-9ae7-3b3842226626")
	req := r.ToHTTPRequest()
	req = AddAuthToken(req, c.AuthToken).WithContext(ctx)
	return c.Client.Do(req)
}

// func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
// 	req = AddAuthToken(req, c.AuthToken).WithContext(ctx)
// 	return c.Client.Do(req)
// }

// func parseRequest(c *Client, r *http.Request) (*http.Request, error) {

// 	return nil, nil
// }

// func ParseParameters(p ...Parameters) (*url.Values, error) {

// }

// func get(c *Client, ctx context.Context, r *http.Request) (*http.Response, error) {

// 	u := *c.BaseURL
// 	// handle authorization
// 	if err := c.Auth(ctx); err != nil {
// 		log.Printf("Error authorizing get request to %v: %v", r.URL, err)
// 		return err
// 	}
// 	r.Query.Add("companyKey", q.Credentials.Get("companyKey"))
// 	// build and send http request
// 	u.RawQuery = r.Query.Encode()
// 	u.Path = path.Join(u.Path, r.Path)
// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
// 	if err != nil {
// 		log.Printf("Error in request to %v: %v", u, err)
// 		return err
// 	}
// 	req.Header = q.Authorization.Token.Clone()
// 	res, err := q.Client.Do(req)
// 	if err != nil {
// 		log.Printf("Error retrieving response from %v: %v", u, err)
// 		return err
// 	}

// 	// handle response
// 	// TODO: improve reading response for larger requests
// 	b, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Printf("Error reading response from %v: %v", u, err)
// 		return err
// 	}
// 	defer res.Body.Close()
// 	resTime, err := http.ParseTime(res.Header.Get("Date"))
// 	if err != nil {
// 		log.Printf("Error parsing date header in response: %v", err)
// 		// accept time.Now as a 'rough' estimate of now
// 		resTime = time.Now()
// 	}

// 	requestConfig := reflect.ValueOf(r.Config)
// 	resourceName := reflect.Indirect(requestConfig).FieldByName("Resource").Interface().(string)
// 	// metadata to capture data heritage
// 	meta := &Metadata{
// 		APIVersion: "v2",
// 		Kind:       "qgenda",
// 		URL:        u.String(),
// 		Name:       resourceName,
// 		Timestamp:  resTime,
// 	}
// 	rr.Response.Metadata = meta
// 	rr.Response.Data = &b

// 	return nil
// }

// // Get handles a *RequestResponse.Request and returns the data and metadata in
// // *RequestResponse.Response
// func (q *Client) Get(ctx context.Context, rr *RequestResponse) error {

// 	u := *q.BaseURL
// 	r := *rr.Request
// 	// handle authorization
// 	if err := q.Auth(ctx); err != nil {
// 		log.Printf("Error authorizing get request to %v: %v", r.Path, err)
// 		return err
// 	}
// 	r.Query.Add("companyKey", q.Credentials.Get("companyKey"))
// 	// build and send http request
// 	u.RawQuery = r.Query.Encode()
// 	u.Path = path.Join(u.Path, r.Path)
// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
// 	if err != nil {
// 		log.Printf("Error in request to %v: %v", u, err)
// 		return err
// 	}
// 	req.Header = q.Authorization.Token.Clone()
// 	res, err := q.Client.Do(req)
// 	if err != nil {
// 		log.Printf("Error retrieving response from %v: %v", u, err)
// 		return err
// 	}

// 	// handle response
// 	// TODO: improve reading response for larger requests
// 	b, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Printf("Error reading response from %v: %v", u, err)
// 		return err
// 	}
// 	defer res.Body.Close()
// 	resTime, err := http.ParseTime(res.Header.Get("Date"))
// 	if err != nil {
// 		log.Printf("Error parsing date header in response: %v", err)
// 		// accept time.Now as a 'rough' estimate of now
// 		resTime = time.Now()
// 	}

// 	requestConfig := reflect.ValueOf(r.Config)
// 	resourceName := reflect.Indirect(requestConfig).FieldByName("Resource").Interface().(string)
// 	// metadata to capture data heritage
// 	meta := &Metadata{
// 		APIVersion: "v2",
// 		Kind:       "qgenda",
// 		URL:        u.String(),
// 		Name:       resourceName,
// 		Timestamp:  resTime,
// 	}
// 	rr.Response.Metadata = meta
// 	rr.Response.Data = &b

// 	return nil
// }
