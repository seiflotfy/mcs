package speech

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	tokenURL  = "https://api.cognitive.microsoft.com/sts/v1.0/issueToken"
	speechURL = "https://speech.platform.bing.com/recognize"

	version = "3.0"
	appID   = "D4D52672-91D7-4C74-8AD8-42B1D98141A5"
	format  = "json"
)

func GetToken() (string, error) {
	apiKey := os.Getenv("MICROSOFT_SPEECH_API_KEY")
	req, err := http.NewRequest("POST", tokenURL, nil)
	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type SpeechRequest struct {
	version   string
	appID     string
	format    string
	token     string
	requestID string

	Locale     string
	DeviceOS   string
	Scenarios  string
	InstanceID string

	Audio io.ReadCloser
}

type SpeechResults struct {
	Results []SpeechResult
}

type SpeechResult struct {
	Scenario   string
	Name       string
	Lexical    string
	Confidence string
	Properties map[string]string
}

func NewSpeechRequest(locale, deviceOS, scenario, instanceID, token string, speech io.ReadCloser) (*SpeechRequest, error) {
	reqID, err := NewUUID()
	if err != nil {
		return nil, err
	}
	reco := &SpeechRequest{
		version:   version,
		appID:     appID,
		format:    format,
		requestID: reqID,
		token:     fmt.Sprintf("Bearer %s", token),

		Locale:     locale,
		DeviceOS:   deviceOS,
		Scenarios:  scenario,
		InstanceID: instanceID,
		Audio:      speech,
	}
	return reco, nil
}

func (reco *SpeechRequest) toMap() map[string]string {
	return map[string]string{
		"version":   reco.version,
		"appID":     reco.appID,
		"format":    reco.format,
		"requestID": reco.requestID,

		"locale":     reco.Locale,
		"device.os":  reco.DeviceOS,
		"scenarios":  reco.Scenarios,
		"instanceID": reco.InstanceID,
	}
}

func Recognize(reco *SpeechRequest) (*SpeechResults, error) {
	req, err := http.NewRequest("POST", speechURL, nil)
	if err != nil {
		return nil, err
	}

	// Set Headers
	req.Header.Set("Authorization", reco.token)
	req.Header.Set("Content-Type", "audio/wav; samplerate=16000")

	// Set Query Keys
	query := req.URL.Query()
	for k, v := range reco.toMap() {
		query.Set(k, v)
	}
	req.URL.RawQuery = query.Encode()

	req.Body = reco.Audio

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results SpeechResults
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	return &results, nil
}
