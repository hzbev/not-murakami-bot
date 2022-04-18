package captcha

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

var client = resty.New()

type PostRes struct {
	ErrorID          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	TaskID           int    `json:"taskId"`
}

type GetTaskRes struct {
	ErrorID          int         `json:"errorId"`
	ErrorCode        interface{} `json:"errorCode"`
	ErrorDescription interface{} `json:"errorDescription"`
	Solution         struct {
		GRecaptchaResponse string `json:"gRecaptchaResponse"`
	} `json:"solution"`
	Status string `json:"status"`
}

type CheckCapStru struct {
	ClientKey string `json:"clientKey"`
	TaskID    int    `json:"taskId"`
}

func PostCaptcha() int {
	res_body := PostRes{}
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"clientKey":"CAPMONSTERKEYHERE","task":{"type":"NoCaptchaTaskProxyless","websiteURL":"https://murakamiflowers.kaikaikiki.com/register/new","websiteKey":"6LeoiQ4eAAAAAH6gsk9b7Xh7puuGd0KyfrI-RHQY"}}`).
		Post("https://api.capmonster.cloud/createTask")
	json.Unmarshal(resp.Body(), &res_body)
	return res_body.TaskID
}

func CheckCapRes(id int) GetTaskRes {
	res_body := GetTaskRes{}
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"clientKey":"CAPMONSTERKEYHERE", "taskId":` + strconv.Itoa(id) + `}`).
		Post("https://api.capmonster.cloud/getTaskResult")
	json.Unmarshal(resp.Body(), &res_body)
	return res_body
}

func RecursiveCaptchaCheck(id int) string {
	for {
		time.Sleep(5 * time.Second)
		i := CheckCapRes(id)
		if i.Status == "ready" {
			return i.Solution.GRecaptchaResponse
		}
	}
}
