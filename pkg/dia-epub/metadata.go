package diaEpub

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type MetaData struct {
	EpubID         any    `json:"epubId"`
	Component      any    `json:"component"`
	Version        any    `json:"version"`
	ImageURL       string `json:"imageUrl"`
	AuthorHomePage string `json:"authorHomePage"`
	MetaData       struct {
		Author             string `json:"author"`
		BookTitle          string `json:"bookTitle"`
		Publisher          string `json:"publisher"`
		PublishDate        string `json:"publishDate"`
		SourcePublishPlace string `json:"sourcePublishPlace"`
		SourcePublishDate  string `json:"sourcePublishDate"`
		SourcePublisher    string `json:"sourcePublisher"`
	} `json:"metaData"`
	OtherBooksWrapper struct {
		AuthorKey    string `json:"authorKey"`
		RecordNumber int    `json:"recordNumber"`
		OtherBooks   any    `json:"otherBooks"`
	} `json:"otherBooksWrapper"`
}

const metaDataPath = "/rest/epub-reader/metadata"

func getMetadata(epubId string) (MetaData, error) {
	urlString, err := url.JoinPath(rootUrl, metaDataPath)

	if err != nil {
		return MetaData{}, err
	}

	url, _ := url.Parse(urlString)
	q := url.Query()

	q.Set("epubId", epubId)
	url.RawQuery = q.Encode()

	// Send the request
	resp, err := http.Get(url.String())
	if err != nil {
		return MetaData{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return MetaData{}, err
	}

	metadata := MetaData{}
	err = json.Unmarshal(body, &metadata)
	if err != nil {
		return MetaData{}, err
	}

	return metadata, nil
}
