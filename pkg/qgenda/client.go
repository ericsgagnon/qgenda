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
	CompanyKey     string `yaml:"companyKey"`
	Password       string
	CacheConfig    *CacheConfig
}

// DefaultClientConfig returns a ClientConfig pointer
// with reasonable defaults. By default, it will recommend
// environment variables in the form ${VARIABLE_NAME}
func DefaultClientConfig() *ClientConfig {

	// cc, err := NewCacheConfig("qgenda")
	// if err != nil {
	// 	panic(err)
	// }

	return &ClientConfig{
		URL:            "https://api.qgenda.com/v2",
		ClientTimeout:  time.Second * 30,
		RequestTimeout: time.Second * 30,
		Email:          "${QGENDA_EMAIL}",
		CompanyKey:     "${QGENDA_COMPANY_KEY}",
		Password:       "${QGENDA_PASSWORD}",
		// CacheConfig:    cc,
	}

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

	if cc == nil {
		cc = DefaultClientConfig()
	}

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

	// yreq, err := yaml.Marshal(req)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(yreq))
	// fmt.Println(req.Method)
	// header := req.Header
	// for k, v := range header {
	// 	fmt.Printf("Request: %25s %s\n", k, v)
	// }
	// fmt.Println(req.Header.Get("Authorization"))
	// fmt.Println(req.Header)
	// fmt.Println(req.Proto)
	// fmt.Println(req.URL.String())
	// req.Header.Set("Authorization", "bearer eyJhbGciOiJBMjU2S1ciLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwidHlwIjoiSldUIn0.Q_BPwbDlA7fCQInwhvoNoYTLpk0zkrBi2FxQw1l3XbHGLrpTO8Z9lPsA7ToAtUKyMd6cPx4Gda4syO8jIFAQNvN2XQydPqbh.7ba2Gz1TESWEZL7F1qWXsg.v6na2rp2dEesCq7vS_lT8qzCOrawNpoGRJVsBr5jFzXq_srlOo-mA4_JUO3RRRG8AfUJyQsykYgCp7ZihhMGVE-iH6K5wNTqkBQykGLMEmG0sWXLI3Znhy8clXoYw5FMv438pjyGE-VmXs4IjwU47nXbWu4qv2S5WmQZxYHUbgbULl8rqSjvijJylMySP7-nM4ypxLMEPRU25AiIR3IKhvFFnnZai1oED1VaI3Pq8wUstuJVelfe-uon0UjZp4HZquY3FVmfxMIz7HoaGdDnSpFuYXX0_7EfSsMtOgO_8aIAT4sd3Uhg5y0FoI8xjvRArf6AjnWjKHxlhmUJzOMc7fgKgJM3b1PCkbcXkqEyPejy9QZzw0GXBFwCQ4tiHCIm8n2wryb4kkW0Nvjfsmft2q1WxgtWwrEmHkyIj_kpYSpP8Xvk4NcR_4hct_U0-iIGeUFCs4-_Y-9Eyq0E7jRWI11JobPORa41Td5G8q-lGj-vlutlreP8IagI_oh6VYsFelNNmw-4G7-KNrbcnGsalDKXHt0E2bwXW6XKby2R5bgVUonY4BV0pRCg1qQBhgH7yeu1i42s_RxJe1BllSYzKOpAaLUpCpipUka9KycjvZl31Siool3ybE30Vk4BKlKlDu1rcBGOs53vYLIRjY3-QP0MpMnu-NCjBcrGqZWMR9BeS2qeEIn0yfX6Z6QB3U2uVFtIJ67nZLPoLl4k9__gwA.zL9gWN-Hyt6e_AMY4ALsNZiNvjoWNv7jXa_PrN6fFT0")
	// fmt.Println(req.Header.Get("Authorization"))
	// fmt.Println()

	return c.Client.Do(req)
	// return &http.Response{}, nil
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
