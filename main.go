package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	diaEpub "github.com/bartaldavid/pim-dia-dl/pkg/dia-epub"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	mux.HandleFunc("/epub", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		urlParts := strings.Split(url, "/")

		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, urlParts[len(urlParts)-1]+".epub"))

		err := diaEpub.UrlToEpub(url, w)

		if err != nil {
			http.Error(w, "Error creating EPUB", http.StatusInternalServerError)
		}

	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	http.ListenAndServe("0.0.0.0:"+port, mux)
}
