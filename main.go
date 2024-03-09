package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-shiori/go-epub"
)

const rootUrl = "https://reader.dia.hu"

func main() {
	url := os.Args[1]

	urlParts := strings.Split(url, "/")

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

	e, err := epub.NewEpub(initSettings.MetaData.Title)

	if err != nil {
		fmt.Println("Error creating new EPUB:", err)
		return
	}

	e.SetAuthor(initSettings.MetaData.Author)

	// served as static css
	cssPath, err := e.AddCSS("monocle.extension.css", "")

	if err != nil {
		fmt.Println("Error adding CSS:", err)
		return
	}

	for _, component := range initSettings.View.Components {
		chunk, chunkTitle, err := getChunk(component, &cookie)
		if err != nil {
			fmt.Println("Error getting chunk:", err)
			return
		}
		_, err = e.AddSection(chunk, chunkTitle, "", cssPath)
		if err != nil {
			fmt.Println("Error adding section:", err)
			return
		}
	}

	// FIXME Do we need this?
	// metadata, err := getMetadata(initSettings.EpubID)
	// if err != nil {
	// 	fmt.Println("Error getting metadata:", err)
	// 	return
	// }

	// fmt.Println(metadata.MetaData.Author, metadata.MetaData.BookTitle)

	err = e.Write("temp/" + urlParts[len(urlParts)-1] + ".epub")
	if err != nil {
		fmt.Println("Error writing EPUB:", err)
		return
	}

}
