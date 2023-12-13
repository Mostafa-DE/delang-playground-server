package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Mostafa-DE/delang/evaluator"
	"github.com/Mostafa-DE/delang/lexer"
	"github.com/Mostafa-DE/delang/object"
	"github.com/Mostafa-DE/delang/parser"
)

func initAppRoutes() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/api/exec", codeExecHandler)
	http.HandleFunc("/api/examples/", examplesHandler)
	http.HandleFunc("/api/de/ubuntu/de_install.sh", DEDownloadUbuntuHandler)
}

func mainHandler(resW http.ResponseWriter, req *http.Request) {
	http.Redirect(resW, req, "https://delang.mostafade.com", http.StatusMovedPermanently)
}

func codeExecHandler(resW http.ResponseWriter, req *http.Request) {
	var res map[string]string

	if req.Method != "POST" {
		methodNotAllowedHandler(resW, req)
		return
	}

	returnObj := createFileToExecFromReqBody(req)
	fileName := returnObj["fileName"]
	errorMessage := returnObj["error"]

	if errorMessage != "" {
		res = map[string]string{
			"error": errorMessage,
		}

		jsonHandler(resW, req, res, fileName)

		return
	}

	fileContents, err := ioutil.ReadFile(fileName)

	if err != nil {
		res = map[string]string{
			"error": errorMessage,
		}

		jsonHandler(resW, req, res, fileName)

		return
	}

	fileContentString := string(fileContents)

	l := lexer.New(fileContentString)
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		res = map[string]string{
			"error": fmt.Sprintf("Parser errors: %s", p.Errors()[0]),
		}

		jsonHandler(resW, req, res, fileName)

		return
	}

	env := object.NewEnvironment()
	env.Set("timeoutLoop", &object.Boolean{Value: true}, false)

	eval := evaluator.Eval(program, env)

	if eval == nil {
		eval = &object.Null{}
	}

	if eval.Type() == object.ERROR_OBJ {
		res = map[string]string{
			"error": fmt.Sprintf("Evaluation error: %s", eval.Inspect()),
		}

		jsonHandler(resW, req, res, fileName)

		return
	}

	logs, logsOk := env.Get("bufferLogs")
	timeOutExceeded, timeoutOk := env.Get("timeoutExceeded")

	if !logsOk {
		logs = &object.Buffer{}
	}

	if !timeoutOk {
		timeOutExceeded = &object.Boolean{Value: false}
	}

	res = map[string]string{
		"logs":    logs.Inspect(),
		"data":    eval.Inspect(),
		"timeout": fmt.Sprintf("%t", timeOutExceeded.(*object.Boolean).Value),
	}

	jsonHandler(resW, req, res, fileName)
}

func examplesHandler(resW http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		methodNotAllowedHandler(resW, req)
		return
	}

	pathname := req.URL.Path
	exampleNumber := pathname[len("/api/examples/"):]

	absPath, err := filepath.Abs(fmt.Sprintf("examples/%s.md", exampleNumber))
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		res := map[string]string{
			"error": "Something went wrong while getting the file, please try again later",
		}

		jsonHandler(resW, req, res, "")

		return
	}

	fileContents, err := ioutil.ReadFile(absPath)

	if err != nil {
		fmt.Println("Error reading file:", err)
		res := map[string]string{
			"error": "Something went wrong while reading the file, please try again later",
		}

		jsonHandler(resW, req, res, "")

		return
	}

	mds := string(fileContents)
	html := mdToHTML([]byte(mds))

	respose := map[string]string{
		"html": string(html),
	}

	jsonHandler(resW, req, respose, "")
}

func DEDownloadUbuntuHandler(resW http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		methodNotAllowedHandler(resW, req)
		return
	}

	filePath := "./de.sh"

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(resW, err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()

	resW.Header().Set("Content-Disposition", "attachment; filename=de_install.sh")

	// Copy the file content to the response writer
	_, err = io.Copy(resW, file)
	if err != nil {
		http.Error(resW, err.Error(), http.StatusInternalServerError)
	}
}
