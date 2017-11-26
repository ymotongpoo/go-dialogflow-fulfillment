//    Copyright 2017 Yoshi Yamaguchi
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

const (
	WelcomeIntent = "input.welcome"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeInput(r)
	if err != nil {
		log.Println(err)
		return
	}

	var res *Response
	intent := req.GetIntent()

	switch intent {
	case WelcomeIntent:
		res, err = welcomeIntent(req)
	}
	if err != nil {
		log.Println(err)
	}

	if err = EncodeOutput(w, res); err != nil {
		log.Println(err)
	}
}

// DecodeInput
func DecodeInput(r *http.Request) (*Request, error) {
	var req Request
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	defer r.Body.Close()
	err := json.NewDecoder(tee).Decode(&req)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v\n", err)
		b, err := ioutil.ReadAll(&buf)
		if err != nil {
			return nil, fmt.Errorf("ioutil error: %v\n", err)
		}
		log.Printf("%s\n", b)
	}
	return &req, nil
}

func EncodeOutput(w http.ResponseWriter, res *Response) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("encode error: %v\n", err)
	}
	return nil
}

func welcomeIntent(r *Request) (*Response, error) {
	w, err := GetWeather("1850147")
	if err != nil {
		return nil, err
	}
	now := time.Now().Format("15時04分")
	template := `こんにちは。時刻は%sです。現在の天気は%s、%.1f度です。` +
		`気圧は%dヘクトパスカル、曇り度数は%dです。湿度は%dパーセントです。風速は秒速%.1fメートル、風光は%sです。`
	voice := fmt.Sprintf(template, now, w.CurWeather, w.CurTemp, w.Pressure, w.Cloudness, w.Humidity, w.WindSpeed, w.WindDirection)
	return NewResponse(voice).SetDisplayText(voice), nil
}
