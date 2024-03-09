package diaEpub

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-shiori/go-epub"
)

const rootUrl = "https://reader.dia.hu"

func UrlToEpub(url string, dest io.Writer) error {
	// urlParts := strings.Split(url, "/")

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return nil
	}

	token := resp.Request.URL.Query().Get("token")

	defer resp.Body.Close()

	cookie := http.Cookie{Name: "token", Value: token}

	initSettings, err := getInitSettings(&cookie)

	if err != nil {
		fmt.Println("Error getting init settings:", err)
		return err
	}

	e, err := epub.NewEpub(initSettings.MetaData.Title)

	if err != nil {
		fmt.Println("Error creating new EPUB:", err)
		return err
	}

	e.SetAuthor(initSettings.MetaData.Author)

	// served as static css
	cssPath, err := e.AddCSS("https://reader.dia.hu/online-reader/resources/epub-reader/extension/monocle.extension.css", "")

	if err != nil {
		fmt.Println("Error adding CSS:", err)
		return err
	}

	for _, component := range initSettings.View.Components {
		chunk, chunkTitle, err := getChunk(component, &cookie)
		if err != nil {
			fmt.Println("Error getting chunk:", err)
			return err
		}
		_, err = e.AddSection(chunk, chunkTitle, "", cssPath)
		if err != nil {
			fmt.Println("Error adding section:", err)
			return err
		}
	}

	// FIXME Do we need this?
	// metadata, err := getMetadata(initSettings.EpubID)
	// if err != nil {
	// 	fmt.Println("Error getting metadata:", err)
	// 	return
	// }

	// fmt.Println(metadata.MetaData.Author, metadata.MetaData.BookTitle)

	// err = e.Write("temp/" + urlParts[len(urlParts)-1] + ".epub")

	_, err = e.WriteTo(dest)

	if err != nil {
		fmt.Println("Error writing EPUB:", err)
		return err
	}

	return nil
}
