package diaEpub

import (
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Chunk struct {
	Body  string
	Title string
}

func getChunk(path string, cookie *http.Cookie) (Chunk, error) {
	urlString, err := url.JoinPath(rootUrl, path)
	if err != nil {
		log.Println("Error joining URL:", err)
		return Chunk{}, err
	}

	// Create a new request
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return Chunk{}, err
	}

	// Set the cookie to the request
	req.AddCookie(cookie)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return Chunk{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error fetching URL:", resp.Status)
		return Chunk{}, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return Chunk{}, err
	}

	body, err := doc.Find("body").Html()

	if err != nil {
		return Chunk{}, err
	}

	chunkTitle := doc.Find(".cim").First().Text()

	return Chunk{Body: body, Title: chunkTitle}, nil
}
