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
	"encoding/json"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
)

const (
	OpenWeatherMapAPI = "http://api.openweathermap.org/data/2.5/weather"
	SecretFile        = "client_secret.json"
)

type Secret struct {
	OpenWeatherMapKey string `json:"openweathermap_key"`
}

type Weather struct {
	CurTemp       float32
	Pressure      int
	Humidity      int
	Cloudness     int
	CurWeather    string
	WindSpeed     float32
	WindDirection string
}

// TODO: create OpenWeatherMap package
type OpenWeatherMapResult struct {
	Weather []OpenWeatherMapWeather `json:"weather"`
	Main    OpenWeatherMapMain      `json:"main"`
	Clouds  struct {
		All int `json:"all"`
	} `json:"clouds"`
	Wind struct {
		Speed  float32 `json:"speed"`
		Degree int     `json:"deg"`
	} `json:"wind"`
}

type OpenWeatherMapWeather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type OpenWeatherMapMain struct {
	Temp     float32 `json:"temp"`
	Pressure int     `json:"pressure"`
	Humidity int     `json:"humidity"`
	TempMin  float32 `json:"temp_min"`
	TempMax  float32 `json:"temp_max"`
}

func GetWeather(id string) (Weather, error) {
	f, err := os.Open(SecretFile)
	if err != nil {
		return Weather{}, err
	}
	var s Secret
	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		return Weather{}, err
	}
	v := url.Values{}
	v.Add("id", id)
	v.Add("appid", s.OpenWeatherMapKey)
	v.Add("units", "metric")
	url := OpenWeatherMapAPI + "?" + v.Encode()
	log.Println(url)
	res, err := http.Get(url)
	if err != nil {
		return Weather{}, err
	}
	var owm OpenWeatherMapResult
	err = json.NewDecoder(res.Body).Decode(&owm)
	if err != nil {
		return Weather{}, err
	}

	return Weather{
		CurTemp:       owm.Main.Temp,
		Humidity:      owm.Main.Humidity,
		Pressure:      owm.Main.Pressure,
		Cloudness:     owm.Clouds.All,
		CurWeather:    JapaneseWeather(owm.Weather[0].Main),
		WindSpeed:     owm.Wind.Speed,
		WindDirection: WindDirection(owm.Wind.Degree),
	}, nil
}

var WeatherMapEnJa = map[string]string{
	"Clear":  "快晴",
	"Sunny":  "晴れ",
	"Cloudy": "曇り",
	"Shower": "にわか雨",
	"Rain":   "雨",
	"Snow":   "雪",
}

func JapaneseWeather(en string) string {
	ja, ok := WeatherMapEnJa[en]
	if !ok {
		return en
	}
	return ja
}

var Direction = map[int]string{
	0:  "北",
	1:  "北北東",
	2:  "北東",
	3:  "東北東",
	4:  "東",
	5:  "東南東",
	6:  "南東",
	7:  "南南東",
	8:  "南",
	9:  "南南西",
	10: "南西",
	11: "西南西",
	12: "西",
	13: "西北西",
	14: "北西",
	15: "北北西",
}

func WindDirection(degree int) string {
	index := int(math.Floor(float64(degree)/22.5 + 0.5))
	return Direction[index]
}
