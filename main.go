package main

import (
	"bufio"
	"fmt"
	"log"
	captcha "murakami/cap"
	"murakami/email"
	"murakami/site"
	"murakami/wallet"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var prod string = "dev"
var verified int = 0
var created int = 0
var content []byte

func main() {

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go create()
		time.Sleep(time.Millisecond * 105)
	}

	wg.Wait()
}
func VerifyCreated(emaill string) {
	a, link := email.CheckEmail(emaill)
	if a != -1 {
		verify(emaill, link)
	}
	time.Sleep(time.Second * 30)
	wg.Done()
}

func create() {
	em_addy := email.GenerateEmail()
	link := createAcc(em_addy)
	// createAcc(em_addy)

	verify(em_addy, link)
	time.Sleep(time.Second * 10)
	wg.Done()
}

func createAcc(emaill string) string {
	token := site.GetAuthToken()
	capId := captcha.PostCaptcha()
	capKey := captcha.RecursiveCaptchaCheck(capId)
	created++
	fmt.Println("created", emaill, created)
	success := site.CreateLogin(emaill, token, capKey)

	if success {
		Write("unverified.txt", emaill)
	}
	link := email.RecursiveCheck(emaill)
	return link
}

func verify(email, urll string) {
	// if checkDone(email) {
	// 	return
	// }
	token := site.GetRegToken(urll)
	capId := captcha.PostCaptcha()
	capKey := captcha.RecursiveCaptchaCheck(capId)
	addy, key := wallet.GenerateWallet()
	t, u := parseURL(urll)
	success := site.SubmitRegister(email, addy, t, u, token, capKey)
	if success {
		verified++
		fmt.Println("verified", email, verified)
		Write("verified.txt", email+";"+addy+";"+key)
	}
}

func parseURL(urllink string) (t, u string) {
	link, _ := url.Parse(urllink)
	v, _ := url.ParseQuery(link.RawQuery)
	return v.Get("t"), v.Get("u")
}

func Write(filePath, text string) {
	var workingDir string
	if prod == "dev" {
		workingDir, _ = os.Getwd()
	} else {
		exeDir, _ := os.Executable()
		workingDir = filepath.Dir(exeDir)
	}

	file, err := os.OpenFile(workingDir+`/`+filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	if _, err := file.Write([]byte(text + "\n")); err != nil {
		log.Fatal(err)
	}

	defer file.Close()
}

func checkDone(email string) bool {
	return strings.Contains(string(content), email)
}

func ReadtoArray(filePath string) []string {
	var workingDir string
	if prod == "dev" {
		workingDir, _ = os.Getwd()
	} else {
		exeDir, _ := os.Executable()
		workingDir = filepath.Dir(exeDir)
	}

	file, err := os.Open(filepath.Join(workingDir, filePath))
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var content []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	return content
}

func Chunks(xs []string, chunkSize int) [][]string {
	if len(xs) == 0 {
		return nil
	}
	divided := make([][]string, (len(xs)+chunkSize-1)/chunkSize)
	prev := 0
	i := 0
	till := len(xs) - chunkSize
	for prev < till {
		next := prev + chunkSize
		divided[i] = xs[prev:next]
		prev = next
		i++
	}
	divided[i] = xs[prev:]
	return divided
}
