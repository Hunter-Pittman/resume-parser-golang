// Name: Hunter Pittman
// Date: 12/5/2020

package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

type Resume struct {
	filename string
	count    []int
}

func main() {
	pdf.DebugOn = true
	var pdfPath string
	var keywordList string
	var outputCSV string

	flag.StringVar(&pdfPath, "p", "pdf", "Specify directory to pdf files")
	flag.StringVar(&keywordList, "k", "keywords", "Specify file with keywords")
	flag.StringVar(&outputCSV, "o", "output", "Specify file path to the csv output")
	flag.Parse()

	//Test value only
	//pdfPath = "C:\\Users\\hunte\\Documents\\pdf_test_dir2\\"
	//keywords := wordlistSeperate("C:\\Users\\hunte\\Documents\\list.txt")

	_, err := os.Stat(pdfPath)
	if err == nil {

	}
	if os.IsNotExist(err) {
		panic("Path to PDF's does not exist!")
	}

	_, err = os.Stat(keywordList)
	if err == nil {

	}
	if os.IsNotExist(err) {
		panic("Path to keyword list does not exist!")
	}

	keywords := wordlistSeperate(keywordList)

	analyzedResumes := searchPdf(pdfPath, keywords)

	generateCSV(analyzedResumes, keywords, outputCSV)

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
			pdfName := pdfPath + f.Name()
			content, err := readPdf(pdfName) // Read local pdf file
			fullPdfPath := pdfPath + f.Name()
			if err != nil {
				fmt.Printf("Error: %v, File Path: %v\n", err, fullPdfPath)
			} else {
				finalCounts := make([]int, 0)
				for _, keyword := range keywords {
					count := strings.Count(content, keyword)

					finalCounts = append(finalCounts, count)

				}
				var sum int = 0
				for _, num := range finalCounts {
					sum += num
				}

				finalCounts = append(finalCounts, sum)

				completedResume := Resume{pdfName, finalCounts}
				output = append(output, completedResume)
			}
		}
	} else {
		content, err := readPdf(pdfPath) // Read local pdf file
		if err != nil {
			permissionError := "malformed PDF: reading at offset 0: stream not present"
			if err.Error() == permissionError {
				fmt.Printf("Try chaging the page extraction permission on the PDF, File Path: %v\n", err, pdfPath)
			} else {
				fmt.Printf("Error: %v, File Path: %v\n", err, pdfPath)
			}

		}

		finalCounts := make([]int, 0)
		for _, keyword := range keywords {
			count := strings.Count(content, keyword)

			finalCounts = append(finalCounts, count)

		}

		var sum int = 0
		for _, num := range finalCounts {
			sum += num
		}

		finalCounts = append(finalCounts, sum)

		completedResume := Resume{pdfPath, finalCounts}
		output = append(output, completedResume)
	}

	return output
}

func generateCSV(resumeData []Resume, keywords []string, outputCSV string) {
	resumeCSVStructure := [][]string{
		{"Filename"},
	}

	// Add each keyword to the initial row

	resumeCSVStructure[0] = append(resumeCSVStructure[0], keywords...)
	resumeCSVStructure[0] = append(resumeCSVStructure[0], "Total")

	for _, resume := range resumeData {
		resumeRow := []string{resume.filename}
		for _, count := range resume.count {
			resumeRow = append(resumeRow, fmt.Sprintf("%d", count))
		}
		resumeCSVStructure = append(resumeCSVStructure, resumeRow)
	}

	csvFile, err := os.Create(outputCSV + "parsedresumes_" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10) + ".csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for _, row := range resumeCSVStructure {
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
	csvFile.Close()
}

func wordlistSeperate(wordlistPath string) []string {
	// Read the contents of the file into a byte slice.
	b, err := ioutil.ReadFile(wordlistPath)
	if err != nil {
		panic(err)
	}

	// Convert the byte slice to a string.
	s := string(b)

	// Split the string into a slice of lines.
	lines := strings.Split(s, "\n")

	return lines

}
