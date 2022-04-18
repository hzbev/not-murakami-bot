package email

import (
	"math/rand"
	randgen "murakami/src"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"mvdan.cc/xurls/v2"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var AllEmails []string = []string{"block521.com", "greencafe24.com", "crepeau12.com", "appzily.com", "coffeetimer24.com", "popcornfarm7.com", "bestparadize.com", "cashflow35.com", "crossmailjet.com", "kobrandly.com", "blondemorkin.com", "block521.com", "popcornfly.com"}

type EmailRes struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type GenMailStruct []string

type CheckMailStruct []struct {
	ID          string        `json:"id"`
	From        string        `json:"from"`
	To          string        `json:"to"`
	Cc          interface{}   `json:"cc"`
	Subject     string        `json:"subject"`
	BodyText    string        `json:"body_text"`
	BodyHTML    string        `json:"body_html"`
	CreatedAt   time.Time     `json:"created_at"`
	Attachments []interface{} `json:"attachments"`
}

type RecEmail struct {
	ID          int           `json:"id"`
	From        string        `json:"from"`
	Subject     string        `json:"subject"`
	Date        string        `json:"date"`
	Attachments []interface{} `json:"attachments"`
	Body        string        `json:"body"`
	TextBody    string        `json:"textBody"`
	HTMLBody    string        `json:"htmlBody"`
}

var client = req.C()

func GenerateEmail() string {
	domain := randgen.RandStringRunes(8)

	res_body := EmailRes{}
	client.R().
		SetResult(&res_body).
		SetBodyJsonBytes([]byte(`{"domain":"` + AllEmails[rand.Intn(len(AllEmails))] + `","name":"` + domain + `"}`)).
		Post("https://api.internal.temp-mail.io/api/v3/email/new")

	return res_body.Email
}

func CheckEmail(email string) (int, string) {
	res_body := CheckMailStruct{}

	client.R().
		SetResult(&res_body).
		Get(`https://api.internal.temp-mail.io/api/v3/email/` + email + `/messages`)

	if len(res_body) == 0 {
		return -1, ""
	} else {
		urll := strings.Index(res_body[0].BodyText, `"Murakami.Flowers" lottery registration`)
		endd := strings.Index(res_body[0].BodyText, `<Note>`)
		return 1, xurls.Relaxed().FindString(res_body[0].BodyText[urll+30 : endd+5])
	}
}

func RecursiveCheck(email string) string {
	for {
		time.Sleep(5 * time.Second)
		i, res_body := CheckEmail(email)
		if i != -1 {
			return res_body
		}
	}
}
