package main

import (
	"fmt"
	"net/http"
)

const rootUrl = "https://reader.dia.hu"

func main() {
	// URL of the HTML file you want to fetch
	const url = "https://reader.dia.hu/document/Csengey_Denes-Mezitlabas_szabadsag-36045"

	// Fetch HTML content
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}

	token := resp.Request.URL.Query().Get("token")

	defer resp.Body.Close()

	cookie := http.Cookie{Name: "token", Value: token}

	initSettings, err := getInitSettings(&cookie)

	if err != nil {
		fmt.Println("Error getting init settings:", err)
		return
	}

	for _, component := range initSettings.View.Components {
		_, err = getChunk(component, &cookie)
		if err != nil {
			fmt.Println("Error getting chunk:", err)
			return
		}
	}

	metadata, err := getMetadata(initSettings.EpubID)
	if err != nil {
		fmt.Println("Error getting metadata:", err)
		return
	}

	fmt.Println(metadata.MetaData.Author, metadata.MetaData.BookTitle)

}
