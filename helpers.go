package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/google/uuid"
)

func createFileToExecFromReqBody(req *http.Request) map[string]string {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)

	if err != nil {
		panic(err)
	}

	fileName := fmt.Sprintf("index_%s.de", uuid.New().String())
	ioutil.WriteFile(fileName, buf.Bytes(), 0644)

	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to read the file: %s", err)
		fmt.Println(errMessage)
		os.Remove(fileName)

		return map[string]string{
			"error":    errMessage,
			"fileName": "",
		}
	}

	fileContentString := string(fileContents)

	ioutil.WriteFile(fileName, []byte(fileContentString), 0644)

	if isFileEmpty(fileName) {
		errMessage := "Code is required"
		fmt.Println(errMessage)
		os.Remove(fileName)

		return map[string]string{
			"error":    errMessage,
			"fileName": "",
		}
	}

	return map[string]string{
		"error":    "",
		"fileName": fileName,
	}
}

func isFileEmpty(fileName string) bool {
	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Failed to read the file:", err)
		os.Remove(fileName)
		return false
	}

	fileContentString := string(fileContents)

	fileContentString = strings.ReplaceAll(fileContentString, " ", "")
	fileContentString = strings.ReplaceAll(fileContentString, "\n", "")

	return len(fileContentString) == 0
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	return port
}

func enableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", " https://delang.mostafade.com/")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)

			return
		}

		// Continue with the actual handler
		handler.ServeHTTP(w, r)
	})
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	pars := parser.NewWithExtensions(extensions)
	doc := pars.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	options := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(options)

	return markdown.Render(doc, renderer)
}
