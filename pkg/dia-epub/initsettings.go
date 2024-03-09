package diaEpub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type InitSettingsResponse struct {
	EpubID                               string `json:"epubId"`
	ProtectedContent                     bool   `json:"protectedContent"`
	OnlineAccessExpirationTime           int    `json:"onlineAccessExpirationTime"`
	OnlineAccessExpirationTimeExtendTime int    `json:"onlineAccessExpirationTimeExtendTime"`
	View                                 struct {
		Components []string `json:"components"`
		Contents   []struct {
			Src      string `json:"src"`
			Title    string `json:"title"`
			Children any    `json:"children"`
		} `json:"contents"`
	} `json:"view"`
	MetaData struct {
		Author string `json:"author"`
		Title  string `json:"title"`
	} `json:"metaData"`
	InitPosition   any      `json:"initPosition"`
	FulltextFields []string `json:"fulltextFields"`
}

const chunkListUrl = "https://reader.dia.hu/rest/epub-reader/init-setting/"

func getInitSettings(cookie *http.Cookie) (InitSettingsResponse, error) {

	// Create a new request
	req, err := http.NewRequest("GET", chunkListUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return InitSettingsResponse{}, err
	}

	req.AddCookie(cookie)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return InitSettingsResponse{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return InitSettingsResponse{}, err
	}

	initSettings := InitSettingsResponse{}
	err = json.Unmarshal(body, &initSettings)

	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return InitSettingsResponse{}, err
	}

	return initSettings, nil
}
