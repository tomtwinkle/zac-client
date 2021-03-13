package zacclient

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.190 Safari/537.36"

const (
	urlLoginPage           = "https://secure.zac.ai/%s/Logon.aspx"
	urlUserLogonPage       = "https://secure.zac.ai/%s/User/user_logon.asp"
	urlUserCheckLogin      = "https://secure.zac.ai/%s/User/user_check.asp"
	urlUserInnerCheckLogin = "https://secure.zac.ai/%s/User/inter_check.asp"
	urlAccountLoginSS         = "https://secure.zac.ai/%s/b/Api/Account/LoginSS"
	urlZacTopPage          = "https://secure.zac.ai/%s/b/top"
	urlShinseiNippouPage   = "https://secure.zac.ai/%s/b/asp/Shinsei/Nippou"
)

type ZACClient interface {
	Login() error
	IsLoggedIn() bool
}

type zacClient struct {
	client  *http.Client
	config  *ZACConfig
	lastReq *url.URL
	debug   bool
}

type ZACConfig struct {
	TenantCode string
	ID         string
	Password   string
}

type Options func(*options)

type options struct {
	debug bool
}

func WithDebug() Options {
	return func(ops *options) {
		ops.debug = true
	}
}

func NewClient(config *ZACConfig, opts ...Options) (ZACClient, error) {
	var opt options
	for _, o := range opts {
		o(&opt)
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		fmt.Printf("%#v,%#v,%#v,%#v\n", req.URL, req.Header, req.Cookies(), req.PostForm)
		return nil
	}
	if config == nil {
		return nil, errors.New("invalid argument: config required")
	}
	if config.TenantCode == "" {
		return nil, errors.New("invalid argument: TenantCode required")
	}
	if config.ID == "" {
		return nil, errors.New("invalid argument: ID required")
	}
	if config.Password == "" {
		return nil, errors.New("invalid argument: Password required")
	}
	return &zacClient{client: client, config: config, debug: opt.debug}, nil
}

func (z *zacClient) get(uri string, headers map[string]string) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	defer z.setLastReq(req.URL)

	if z.debug {
		reqDump, err := httputil.DumpRequest(req, false)
		if err != nil {
			return nil, err
		}
		log.Printf("request=%q\n", reqDump)
		log.Printf("request cookies=%v\n", z.client.Jar.Cookies(req.URL))
	}
	res, err := z.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if z.debug {
		resDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}
		log.Printf("response=%q\n", resDump)
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusFound {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, errors.New(res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// nolint:unparam
func (z *zacClient) post(uri string, headers map[string]string, body url.Values) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	if ref := z.refererForURL(z.lastReq); ref != "" {
		req.Header.Set("Referer", ref)
	}
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	defer z.setLastReq(req.URL)

	if z.debug {
		reqDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		log.Printf("request=%q\n", reqDump)
		log.Printf("request cookies=%v\n", z.client.Jar.Cookies(req.URL))
	}
	res, err := z.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if z.debug {
		resDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}
		log.Printf("response=%q\n", resDump)
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusFound {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, errors.New(res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	z.lastReq = req.URL
	return doc, nil
}

func (z *zacClient) refererForURL(lastReq *url.URL) string {
	if lastReq == nil {
		return ""
	}
	return lastReq.String()
}

func (z *zacClient) setLastReq(lastReq *url.URL) {
	z.lastReq = lastReq
}
