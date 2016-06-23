package http

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type Request struct {
	Method         string
	Url            string
	Params         map[string]string
	CurrentCookies []*http.Cookie
}

func (r *Request) Request() (data string, err error) {
	startTime := time.Now().UnixNano()
	fmt.Println("===============================================")
	fmt.Println("Cookies:", r.CurrentCookies)
	fmt.Println("Method:", r.Method)
	switch r.Method {
	case "POST":
		data, err = r.httpPost()
		fmt.Println("Url:", r.Url)
		fmt.Println("Params:", r.Params)
	case "GET":
		data, err = r.httpGet()
		fmt.Println("Url:", r.Url)
	}
	fmt.Println("Return:", string(data))
	endTime := time.Now().UnixNano()
	fmt.Println("Time(ms):", (endTime-startTime)/int64(time.Millisecond))
	fmt.Println("---------------------------------------------------------------------------------")
	return data, err
}

func (r *Request) httpPost() (string, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	content_type := w.FormDataContentType()
	//form
	for k, v := range r.Params {
		w.WriteField(k, v)
	}
	w.Close()
	//request
	req, _ := http.NewRequest("POST", r.Url, body)
	req.Header.Set("Content-Type", content_type)
	// add cookie
	for _, v := range r.CurrentCookies {
		req.AddCookie(v)
	}
	resp, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("POST request failed !")
	}
	// set cookie
	r.CurrentCookies = resp.Cookies()
	return string(data), nil
}

func (r *Request) httpGet() (string, error) {
	body := new(bytes.Buffer)
	content_type := "application/x-www-form-urlencoded; param=value"
	//request
	req, _ := http.NewRequest("GET", r.Url, body)
	req.Header.Set("Content-Type", content_type)
	// add cookie
	for _, v := range r.CurrentCookies {
		req.AddCookie(v)
	}
	resp, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("GET request failed !")
	}
	return string(data), nil
}
