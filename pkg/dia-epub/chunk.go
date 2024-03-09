package diaEpub

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// TODO make this a struct
func getChunk(path string, cookie *http.Cookie) (string, string, error) {
	urlString, err := url.JoinPath(rootUrl, path)
	if err != nil {
		fmt.Println("Error joining URL:", err)
		return "", "", err
	}

	// Create a new request
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", "", err
	}

	// Set the cookie to the request
	req.AddCookie(cookie)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error fetching URL:", resp.Status)
		return "", "", err
	}

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error reading response body:", err)
	// 	return "", err
	// }

	if err != nil {
		fmt.Println("Error getting HTML body:", err)
		return "", "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return "", "", err
	}

	body, err := doc.Find("body").Html()

	if err != nil {
		return "", "", err
	}

	chunkTitle := doc.Find(".cim").First().Text()

	return string(body), chunkTitle, nil
}
