// Path: mime.go
// A simple static HTTP server that supports custom mime-types.

package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func addMimeTypes() error {
	log.Println("Initializing mime-types ...")

	mimeTypes, err := loadMimeTypes(*mimeConfig)
	if err != nil {
		return err
	}

	for ext, typ := range mimeTypes {
		log.Println("Adding mime-type", typ, "for extension", "."+ext, "...")
		if err := mime.AddExtensionType("."+ext, typ); err != nil {
			return err
		}
	}

	return nil
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving file", *dir+r.URL.Path, "with HTTP method", r.Method, "...")
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.URL.Path)))
	http.ServeFile(w, r, *dir+r.URL.Path)
}

func loadMimeTypes(path string) (map[string]string, error) {
	log.Println("Loading mime-types from", path, "...")

	mimeTypesJSONFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer mimeTypesJSONFile.Close()
	log.Println("mimeTypesJSONFile:", mimeTypesJSONFile)

	var mimeTypes map[string][]string
	if err := json.NewDecoder(mimeTypesJSONFile).Decode(&mimeTypes); err != nil {
		return nil, err
	}
	log.Println("mimeTypes:", mimeTypes)

	// Convert the mime-types to a map.
	mimeTypesMap := make(map[string]string)
	for typ, exts := range mimeTypes {
		for _, ext := range exts {
			mimeTypesMap[ext] = typ
		}
	}
	log.Println("mimeTypesMap:", mimeTypesMap)

	return mimeTypesMap, nil
}
