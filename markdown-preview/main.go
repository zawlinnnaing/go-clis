package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html><html><head><meta http-equiv="content-type" content="text/html; charset=utf-8"> <title>{{ .Title }}</title> </head> <body> {{ .Body }} </body> </html>`
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	templateFile := flag.String("t", "", "Alternative HTML template file")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout, *skipPreview, *templateFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(filename string, writer io.Writer, skipPreview bool, templateFile string) error {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData, err := parseContent(input, templateFile)
	if err != nil {
		return err
	}
	temp, err := ioutil.TempFile("", "mdp-*.html")
	if err != nil {
		return err
	}
	if err = temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()
	fmt.Fprintln(writer, outName)

	err = saveHTML(outName, htmlData)
	if err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)

	return preview(outName)
}

func parseContent(input []byte, templateFileName string) ([]byte, error) {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)
	htmlTemplate, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}
	if templateFileName != "" {
		htmlTemplate, err = template.ParseFiles(templateFileName)
		if err != nil {
			return nil, err
		}
	}

	templateContent := content{
		Title: "Markdown Preview Tool from Template",
		Body:  template.HTML(body),
	}

	var buffer bytes.Buffer
	if err = htmlTemplate.Execute(&buffer, templateContent); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func saveHTML(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

func preview(filename string) error {
	cName := ""
	cParams := []string{}

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	cParams = append(cParams, filename)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	err = exec.Command(cPath, cParams...).Run()

	time.Sleep(2 * time.Second)

	return err
}
