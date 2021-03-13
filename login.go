package zacclient

import (
	"fmt"
	"github.com/pkg/errors"
	"io/fs"
	"io/ioutil"
	"net/url"
)

func (z *zacClient) Login() error {
	// Get login form data
	formData, err := z.getSecureLogonPage()
	if err != nil {
		return err
	}

	formData, err = z.secureLogon(formData)
	if err != nil {
		return err
	}

	//if err := z.secureUserLogon(formData); err != nil {
	//	return err
	//}

	if err := z.secureUserCheck(); err != nil {
		return err
	}

	if err := z.secureInnerCheck(); err != nil {
		return err
	}

	if err := z.secureLoginSS(); err != nil {
		return err
	}

	return nil
}

func (z *zacClient) getSecureLogonPage() (url.Values, error) {
	// https://secure.zac.ai/{tenantCode}/Logon.aspx
	uri := fmt.Sprintf(urlLoginPage, z.config.TenantCode)

	formKeys := []string{
		"__VIEWSTATE",
		"__VIEWSTATEGENERATOR",
		"__VIEWSTATEENCRYPTED",
		"__EVENTVALIDATION",
	}

	doc, err := z.get(uri, nil)
	if err != nil {
		return nil, err
	}

	formData := make(url.Values, len(formKeys))
	for _, key := range formKeys {
		if val, ok := doc.Find(fmt.Sprintf("input[name=%s]", key)).Attr("value"); ok {
			formData.Add(key, val)
		} else {
			return nil, errors.Errorf("%s not found", key)
		}
	}
	return formData, nil
}

func (z *zacClient) secureLogon(body url.Values) (url.Values, error) {
	// https://secure.zac.ai/{tenantCode}/Logon.aspx
	uri := fmt.Sprintf(urlLoginPage, z.config.TenantCode)

	body.Add("Login1$UserName", z.config.ID)
	body.Add("Login1$Password", z.config.Password)

	doc, err := z.post(uri, nil, body)
	if err != nil {
		return nil, err
	}

	formKeys := []string{
		"__VIEWSTATE",
		"__VIEWSTATEGENERATOR",
		"__VIEWSTATEENCRYPTED",
		"__EVENTVALIDATION",
	}
	formData := make(url.Values, len(formKeys))
	for _, key := range formKeys {
		if val, ok := doc.Find(fmt.Sprintf("input[name=%s]", key)).Attr("value"); ok {
			formData.Add(key, val)
		} else {
			return nil, errors.Errorf("%s not found", key)
		}
	}
	return formData, nil
}

func (z *zacClient) secureUserCheck() error {
	// https://secure.zac.ai/{tenantCode}/User/user_check.asp
	uri := fmt.Sprintf(urlUserCheckLogin, z.config.TenantCode)

	body := make(url.Values, 2)
	body.Add("user_name", z.config.ID)
	body.Add("password", z.config.Password)

	headers := make(map[string]string)
	headers["origin"] = "https://secure.zac.ai"
	headers["referer"] = "https://secure.zac.ai/beex/User/user_logon.asp"
	headers["dnt"] = "1"
	headers["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	headers["accept-encoding"] = "gzip, deflate, br"
	headers["accept-language"] = "ja-JP,ja;q=0.9"
	headers["cache-control"] = "max-age=0"
	headers["sec-ch-ua-mobile"] = "?0"
	headers["sec-fetch-dest"] = "document"
	headers["sec-fetch-mode"] = "navigate"
	headers["sec-fetch-site"] = "same-origin"
	headers["sec-fetch-user"] = "?1"
	headers["upgrade-insecure-requests"] = "1"

	doc, err := z.post(uri, headers, body)
	if err != nil {
		return err
	}
	if html, err := doc.Html(); err == nil {
		if err := ioutil.WriteFile("secureUserCheck.html", []byte(html), fs.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (z *zacClient) secureInnerCheck() error {
	// https://secure.zac.ai/{tenantCode}/User/inter_check.asp
	uri := fmt.Sprintf(urlUserInnerCheckLogin, z.config.TenantCode)

	headers := make(map[string]string)
	headers["origin"] = "https://secure.zac.ai"
	headers["referer"] = "https://secure.zac.ai/beex/User/user_logon.asp"
	headers["dnt"] = "1"
	headers["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	headers["accept-encoding"] = "gzip, deflate, br"
	headers["accept-language"] = "ja-JP,ja;q=0.9"
	headers["cache-control"] = "max-age=0"
	headers["sec-ch-ua-mobile"] = "?0"
	headers["sec-fetch-dest"] = "document"
	headers["sec-fetch-mode"] = "navigate"
	headers["sec-fetch-site"] = "same-origin"
	headers["sec-fetch-user"] = "?1"
	headers["upgrade-insecure-requests"] = "1"

	doc, err := z.get(uri, headers)
	if err != nil {
		return err
	}
	if html, err := doc.Html(); err == nil {
		if err := ioutil.WriteFile("secureInnerCheck.html", []byte(html), fs.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (z *zacClient) secureLoginSS() error {
	// https://secure.zac.ai/{tenantCode}/b/Api/Account/LoginSS
	uri := fmt.Sprintf(urlAccountLoginSS, z.config.TenantCode)

	headers := make(map[string]string)
	headers["origin"] = "https://secure.zac.ai"
	headers["referer"] = "https://secure.zac.ai/beex/User/user_logon.asp"
	headers["dnt"] = "1"
	headers["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	headers["accept-encoding"] = "gzip, deflate, br"
	headers["accept-language"] = "ja-JP,ja;q=0.9"
	headers["cache-control"] = "max-age=0"
	headers["sec-ch-ua-mobile"] = "?0"
	headers["sec-fetch-dest"] = "script"
	headers["sec-fetch-mode"] = "no-cors"
	headers["sec-fetch-site"] = "same-origin"

	doc, err := z.get(uri, headers)
	if err != nil {
		return err
	}
	if html, err := doc.Html(); err == nil {
		if err := ioutil.WriteFile("secureLoginSS.html", []byte(html), fs.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (z *zacClient) IsLoggedIn() bool {
	// https://secure.zac.ai/{tenantCode}/b/top
	uri := fmt.Sprintf(urlZacTopPage, z.config.TenantCode)

	headers := make(map[string]string)
	headers["origin"] = "https://secure.zac.ai"
	headers["referer"] = "https://secure.zac.ai/beex/User/user_logon.asp"
	headers["dnt"] = "1"
	headers["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	headers["accept-encoding"] = "gzip, deflate, br"
	headers["accept-language"] = "ja-JP,ja;q=0.9"
	headers["cache-control"] = "max-age=0"
	headers["sec-ch-ua-mobile"] = "?0"
	headers["sec-fetch-dest"] = "document"
	headers["sec-fetch-mode"] = "navigate"
	headers["sec-fetch-site"] = "same-origin"
	headers["upgrade-insecure-requests"] = "1"

	if doc, err := z.get(uri, headers); err != nil {
		return false
	} else {
		if html, err := doc.Html(); err == nil {
			if err := ioutil.WriteFile("top.html", []byte(html), fs.ModePerm); err != nil {
				return false
			}
		}
	}
	return true
}
