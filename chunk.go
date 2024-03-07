package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getChunk(path string, cookie *http.Cookie) (string, error) {
	urlString, err := url.JoinPath(rootUrl, path)
	if err != nil {
		fmt.Println("Error joining URL:", err)
		return "", err
	}

	// Create a new request
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Set the cookie to the request
	req.AddCookie(cookie)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error fetching URL:", resp.Status)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	pathSegments := strings.Split(path, "/")

	err = os.WriteFile(pathSegments[len(pathSegments)-1], body, 0644)

	if err != nil {
		fmt.Println("Error writing file:", err)
		return "", err
	}

	return string(body), nil
}