package site

import (
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
)

var client = req.C()

type NameStru []string

func GetName() string {
	res_body := NameStru{}
	client.R().
		SetResult(&res_body).
		SetFormData(map[string]string{
			"type":   "firstname",
			"number": "1",
		}).
		Post("https://randommer.io/Name")
	return res_body[0]
}

func GetAuthToken() string {
	resp, _ := client.R().
		Get("https://murakamiflowers.kaikaikiki.com/register/new")
	if resp.StatusCode == 502 {
		fmt.Println(resp.StatusCode)
		return GetAuthToken()
	}
	res := resp.String()
	firstInd := strings.Index(res, `authenticity_token" value="`)
	lastInd := strings.Index(res, `autocomplete="off" /><input name="t" value="new"`)
	return res[firstInd+27 : lastInd-2]
}

func GetRegToken(url string) string {
	resp, _ := client.R().
		Get(url)
	if resp.StatusCode == 502 {
		fmt.Println(resp.StatusCode, url)
		return GetRegToken(url)
	}
	res := resp.String()
	firstInd := strings.Index(res, `authenticity_token" value="`)
	lastInd := strings.Index(res, `autocomplete="off" /><div class="p-accountFormOuter">`)
	return res[firstInd+27 : lastInd-2]
}

func CreateLogin(email, token, recap string) bool {
	resp, _ := client.R().
		SetFormData(map[string]string{
			"authenticity_token":   token,
			"t":                    "new",
			"email":                email,
			"g-recaptcha-response": recap,
			"commit":               "SEND+REGISTRATION+MAIL",
		}).
		Post("https://murakamiflowers.kaikaikiki.com/register/new_account")
	res := resp.String()
	contain := strings.Contains(res, `class="p-accountForm__txt">We will send you the URL for registration soon.</p><p class="p-accountForm__txt">Please check your mailbox.</p><div class="p-accountForm__btnOuter"><a class="c-`)
	return contain
}

func SubmitRegister(email, wallet, t, u, token, recap string) bool {
	newName := GetName()
	resp, _ := client.R().
		SetFormData(map[string]string{
			"authenticity_token":            token,
			"user[email]":                   email,
			"user[name]":                    newName,
			"user[metamask_wallet_address]": wallet,
			"g-recaptcha-response":          recap,
			"user[password]":                "YUJGyuwejdgh62343",
			"user[password_confirmation]":   "YUJGyuwejdgh62343",
			"t":                             t,
			"u":                             u,
			"commit":                        "Confirm",
		}).
		Post("https://murakamiflowers.kaikaikiki.com/register/register")
	res := resp.String()
	contain := strings.Contains(res, `Thank you for your registering.`)
	return contain
}
