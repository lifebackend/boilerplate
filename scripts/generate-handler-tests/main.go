package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	serviceName, exists := os.LookupEnv("SERVICE_NAME")
	if !exists {
		log.Fatal("expected ENV variable SERVICE_NAME")
	}

	srcFileWithHandlers := fmt.Sprintf("internal/app/%s/server/restapi/configure_%s.go", serviceName, serviceName)
	dstFilePath := fmt.Sprintf("internal/app/%s/handlers/", serviceName)

	data, err := ioutil.ReadFile(srcFileWithHandlers)
	if err != nil {
		log.Fatal(err)
	}

	dataStr := string(data)

	// locate imports
	var importsBody string
	{
		// locate handlers
		beginPattern := "/* default handlers import"
		endPattern := "default handlers import */"

		startPosition := strings.Index(dataStr, beginPattern)
		if startPosition == 0 {
			log.Fatalf("expected pattern '%s' not found", beginPattern)
		}
		startPosition += len(beginPattern) + 1

		endPosition := strings.Index(dataStr, endPattern)
		if endPosition == 0 {
			log.Fatalf("expected pattern '%s' not found", endPattern)
		}
		if endPosition < startPosition {
			log.Fatalf("expected pattern '%s' not found", endPattern)
		}

		importsBody = dataStr[startPosition:endPosition]
	}

	var handlersBody []string
	for {
		// locate handlers
		beginPattern := "/* default handler for"
		endPattern := "default handler */"

		startPosition := strings.Index(dataStr, beginPattern)
		if startPosition == -1 {
			break
		}
		startPosition += len(beginPattern) + 1

		endPosition := strings.Index(dataStr, endPattern)
		if endPosition == 0 {
			log.Fatalf("expected pattern '%s' not found", endPattern)
		}
		if endPosition < startPosition {
			log.Fatalf("expected pattern '%s' not found", endPattern)
		}

		handlersBody = append(handlersBody, dataStr[startPosition:endPosition])
		dataStr = dataStr[endPosition+len(endPattern):]
	}

	for _, handlerBody := range handlersBody {
		fileName, rawHandlerBody, err := getFilenameFromHandlerBody(handlerBody)
		if err != nil {
			log.Fatalf("unable to get filename from the handler '%s' not found, Err: %s", handlerBody, err)
		}

		if fileName == "" {
			continue
		}

		testName, err := getTestNameFromHandlerBody(rawHandlerBody)
		if err != nil {
			log.Fatalf("unable to get test name from the handler '%s', Err: %s", handlerBody, err)
		}

		dstFile := filepath.Join(dstFilePath, fileName)

		if _, err = os.Stat(dstFile); !os.IsNotExist(err) {
			// don't overwrite existing files
			continue
		}

		// create target file
		f, err := os.Create(dstFile)
		if err != nil {
			log.Fatalf("unable to create a file '%s'. Err: %s", dstFile, err)
		}

		operationStructureName := fmt.Sprintf(
			"%s%s",
			strings.ToUpper(serviceName[:1]),
			strings.ToLower(serviceName[1:]),
		)

		// render template and write into the file
		err = packageTemplate.Execute(
			f,
			struct {
				RawImportsBody         template.HTML
				RawHandlerBody         template.HTML
				RawTestName            template.HTML
				OperationStructureName string
			}{
				RawImportsBody:         template.HTML(importsBody),    //nolint:gosec
				RawHandlerBody:         template.HTML(rawHandlerBody), //nolint:gosec
				RawTestName:            template.HTML(testName),       //nolint:gosec
				OperationStructureName: operationStructureName,
			},
		)
		if err != nil {
			log.Println(err)
		}

		err = f.Close()
		if err != nil {
			log.Fatalf("unable to close the file '%s'. Err: %s", dstFile, err)
		}

		err = exec.Command("goimports", "-w", dstFile).Start()
		if err != nil {
			log.Printf("unable format file '%s'. Err: %s", dstFile, err)
		}
	}
}

func getFilenameFromHandlerBody(data string) (string, string, error) {
	lines := strings.Split(data, "\n")
	if len(lines) == 0 {
		return "", "", fmt.Errorf("unable to split data on lines")
	}

	line := lines[0]
	line = strings.ReplaceAll(line, " ", "")
	line = strings.ReplaceAll(line, "/", "-")
	line = strings.ReplaceAll(line, "}", "")
	line = strings.ReplaceAll(line, "{", "")
	line = strings.ReplaceAll(line, "_", "")
	line = strings.ReplaceAll(line, "--", "-")
	line = strings.ToLower(line)
	line = line[1:]
	line = fmt.Sprintf("%s_test.go", line)

	fileBody := strings.Join(lines[1:], "\n")

	return line, fileBody, nil
}

func getTestNameFromHandlerBody(data string) (string, error) {
	lines := strings.Split(data, "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("unable to split data on lines")
	}

	line := lines[0]
	line = strings.ReplaceAll(line, "//", "")
	line = strings.TrimSpace(line)
	line = strings.Split(line, " ")[0]

	return line, nil
}

var (
	packageTemplate = template.Must(template.New("").Parse(tpl))
	tpl             = `package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	{{ .RawImportsBody }}
)

	func (s *{{ .OperationStructureName }}Suite) Test{{ .RawTestName }}() {
		s.T().SkipNow()
	
		t := s.T()
		t.Parallel()
}
`
)
