package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	diaEpub "github.com/bartaldavid/pim-dia-dl/pkg/dia-epub"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	mux.HandleFunc("/epub", func(w http.ResponseWriter, r *http.Request) {
		urlQuery := r.URL.Query().Get("url")

		if urlQuery == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		url, err := url.Parse(urlQuery)

		if err != nil {
			http.Error(w, fmt.Sprintf("Malformed URL: %s", urlQuery), http.StatusBadRequest)
			return
		}

		url.Scheme = "https"

		if url.Hostname() != "reader.dia.hu" {
			http.Error(w, fmt.Sprintf("Invalid URL: %s", urlQuery), http.StatusBadRequest)
			return
		}

		epub, err := diaEpub.UrlToEpub(urlQuery)

		if err != nil {
			http.Error(w, "Error creating EPUB - "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/epub+zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, epub.FileName))

		_, err = epub.Epub.WriteTo(w)

		if err != nil {
			http.Error(w, "Error writing EPUB", http.StatusInternalServerError)
			return
		}

		fmt.Print(epub.FileName)
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, mux)
}
