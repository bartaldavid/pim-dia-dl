package diaEpub

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-shiori/go-epub"
	"golang.org/x/sync/errgroup"
)

const rootUrl = "https://reader.dia.hu"

type EpubResult struct {
	Epub     *epub.Epub
	FileName string
}

func addChildrenRecursively(e *epub.Epub, parentPath string, children *[]Contents, chunks map[string]Chunk, cssPath string) error {
	if children == nil {
		return nil
	}
	for _, subContent := range *children {
		sectionPath, err := e.AddSubSection(parentPath, chunks[subContent.Src].Body, subContent.Title, "", cssPath)
		if err != nil {
			fmt.Println("Error adding subsection:", err)
			return err
		}
		if subContent.Children != nil {
			addChildrenRecursively(e, sectionPath, subContent.Children, chunks, cssPath)
		}
	}
	return nil
}

func UrlToEpub(url string) (EpubResult, error) {
	// urlParts := strings.Split(url, "/")

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return EpubResult{}, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error fetching URL:", resp.Status)
		return EpubResult{}, fmt.Errorf("error fetching URL (%s): %s", url, resp.Status)
	}

	token := resp.Request.URL.Query().Get("token")

	if token == "" {
		fmt.Println("No token found")
		return EpubResult{}, fmt.Errorf("no token found")
	}

	urlParts := strings.Split(url, "/")

	defer resp.Body.Close()

	cookie := http.Cookie{Name: "token", Value: token}

	initSettings, err := getInitSettings(&cookie)

	if err != nil {
		fmt.Println("Error getting init settings:", err)
		return EpubResult{}, err
	}

	e, err := epub.NewEpub(initSettings.MetaData.Title)

	if err != nil {
		fmt.Println("Error creating new EPUB:", err)
		return EpubResult{}, err
	}

	e.SetAuthor(initSettings.MetaData.Author)
	e.SetLang("hu")

	cssPath, err := e.AddCSS("https://reader.dia.hu/online-reader/resources/epub-reader/extension/monocle.extension.css", "")

	if err != nil {
		fmt.Println("Error adding CSS:", err)
		return EpubResult{}, err
	}

	chunks := make(map[string]Chunk)
	g := new(errgroup.Group)

	for _, component := range initSettings.View.Components {
		g.Go(func() error {
			chunk, err := getChunk(component, &cookie)
			if err != nil {
				fmt.Println("Error getting chunk:", err)
				return err
			}
			chunks[component] = chunk
			return nil
		})

		// chunk, err := getChunk(component, &cookie)
		// if err != nil {
		// 	fmt.Println("Error getting chunk:", err)
		// 	return EpubResult{}, err
		// }
		// _, err = e.AddSection(chunk.Body, chunk.Title, "", cssPath)
		// if err != nil {
		// 	fmt.Println("Error adding section:", err)
		// 	return EpubResult{}, err
		// }
	}

	if err := g.Wait(); err != nil {
		fmt.Println("Error getting chunks:", err)
		return EpubResult{}, err
	}

	for _, content := range initSettings.View.Contents {
		parent, err := e.AddSection(chunks[content.Src].Body, content.Title, "", cssPath)
		if err != nil {
			fmt.Println("Error adding section:", err)
			return EpubResult{}, err
		}
		err = addChildrenRecursively(e, parent, content.Children, chunks, cssPath)
		if err != nil {
			fmt.Println("Error adding children:", err)
			return EpubResult{}, err
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

	filename := urlParts[len(urlParts)-1] + ".epub"

	return EpubResult{Epub: e, FileName: filename}, nil
}
