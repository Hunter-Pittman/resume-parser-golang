// Name: Hunter Pittman
// Date: 12/5/2020

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

type KeywordInfo struct {
	keywords, count interface{}
}

type Resume struct {
	filename    string
	keywordInfo KeywordInfo
}

func main() {
	pdf.DebugOn = true
	var pdfPath string
	var keywordList string

	flag.StringVar(&pdfPath, "p", "pdf", "Specify directory to pdf files")
	flag.StringVar(&keywordList, "k", "keywords", "Specify file with keywords")
	flag.Parse()

	//Test value only
	pdfPath = "C:\\Users\\hunte\\Documents\\pdf_test_dir\\"

	keywords := []string{"sponsor", "cyber", "security"}

	analyzedResumes := searchPdf(pdfPath, keywords)

	fmt.Printf("%v", analyzedResumes)

}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
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

func searchPdf(pdfPath string, keywords []string) []Resume {
	objectInfo, err := os.Stat(pdfPath)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	output := make([]Resume, 0)
	if objectInfo.IsDir() {
		files, err := ioutil.ReadDir(pdfPath)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		for _, f := range files {
			content, err := readPdf(pdfPath + f.Name()) // Read local pdf file
			fullPdfPath := pdfPath + f.Name()
			if err != nil {
				fmt.Printf("Error: %v, File Path: %v\n", err, fullPdfPath)
			}

			for _, keyword := range keywords {
				count := strings.Count(content, keyword)
				//fmt.Printf("Found %s %d times in %s\n", keyword, count, pdfPath)

				completedKeywordInfo := KeywordInfo{keyword, count}
				completedResume := Resume{fullPdfPath, completedKeywordInfo}
				output = append(output, completedResume)
			}
		}
	} else {
		content, err := readPdf(pdfPath) // Read local pdf file
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		for _, keyword := range keywords {
			count := strings.Count(content, keyword)
			//fmt.Printf("Found %s %d times in %s\n", keyword, count, pdfPath)

			completedKeywordInfo := KeywordInfo{keyword, count}
			completedResume := Resume{pdfPath, completedKeywordInfo}
			output = append(output, completedResume)
		}
	}

	return output
}

func generateCSV(resumeData []Resume) {

}
