package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/dslipak/pdf"
)

func main() {
	pwd := os.Args[1]
	if pwd == "" {
		fmt.Println("Please provide a password")
		os.Exit(1)
	}

	path := os.Args[2]
	if path == "" {
		fmt.Println("Please provide a path to a PDF file")
		os.Exit(1)
	}

	content, err := readPdf(path, pwd) // Read local pdf file
	if err != nil {
		panic(err)
	}

	r, err := regexp.Compile(`(?P<NAME>.{0,30}?)(?:PARC\.(?P<INST>\d+?\/\d+?))*?(?P<DATE>\d{2}\/\d{2})(?P<VALUE>\d+?,\d{2})`)
	if err != nil {
		panic(err)
	}

	matches := r.FindAllStringSubmatch(content, -1)

	f, err := os.Create("expenses.csv")

	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("NAME,INST,DATE,VALUE\n")
	for _, match := range matches {
		f.WriteString(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\"\n", match[1], match[2], match[3], match[4]))
	}
}

func readPdf(path string, pwd string) (string, error) {
	f, err := os.Open(path)
	// defer f.Close()
	if err != nil {
		return "", err
	}
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	pwdf := func() string { return pwd }
	r, err := pdf.NewReaderEncrypted(f, fi.Size(), pwdf)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
