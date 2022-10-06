package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	joke, err := getDadJoke()
	if err != nil {
		panic(err)
	}

	pirateSpeak, err := translateToPirate(joke)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Original joke: %s\nPirate speak: %s\n", joke, pirateSpeak)
}

type joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

func getDadJoke() (string, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com", nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// decode response into joke
	var j joke
	err = json.NewDecoder(res.Body).Decode(&j)
	if err != nil {
		return "", err
	}

	return j.Joke, nil
}

type pirateSpeak struct {
	Success struct {
		Total int `json:"total"`
	} `json:"success"`
	Contents struct {
		Translated  string `json:"translated"`
		Text        string `json:"text"`
		Translation string `json:"translation"`
	} `json:"contents"`
}

func translateToPirate(input string) (string, error) {
	params := url.Values{}
	params.Add("text", input)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://api.funtranslations.com/translate/pirate.json", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Funtranslations-Api-Secret", "<api_key>")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	// decode response into pirateSpeak
	var p pirateSpeak
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return "", err
	}

	return p.Contents.Translated, nil
}
