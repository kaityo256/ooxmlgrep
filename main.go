package main

import (
	"archive/zip"
	"encoding/xml"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const version = "ooxmlgrep 1.0.0"
const red = "\033[31m"
const reset = "\033[0m"

func extractText(xmlData string) string {
	var result strings.Builder
	decoder := xml.NewDecoder(strings.NewReader(xmlData))

	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "t" {
				var content string
				decoder.DecodeElement(&content, &se)
				result.WriteString(content + " ")
			}
		}
	}
	return result.String()
}

func highlight(text, keyword string, ignoreCase bool) string {
	if ignoreCase {
		// 正規表現で大文字小文字無視のマッチを強調
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(keyword))
		return re.ReplaceAllStringFunc(text, func(m string) string {
			return red + m + reset
		})
	} else {
		return strings.ReplaceAll(text, keyword, red+keyword+reset)
	}
}

func searchInFile(keyword, filename string, ignoreCase, onlyNumber bool) {
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return
	}
	defer reader.Close()

	re := regexp.MustCompile(`ppt/slides/slide([0-9]+)\.xml$`)
	for _, file := range reader.File {
		m := re.FindStringSubmatch(file.Name)
		if m != nil {
			slideNum := m[1]
			rc, err := file.Open()
			if err != nil {
				continue
			}
			data, err := ioutil.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}
			text := extractText(string(data))
			matchText := text
			matchKeyword := keyword

			if ignoreCase {
				matchText = strings.ToLower(text)
				matchKeyword = strings.ToLower(keyword)
			}

			if strings.Contains(matchText, matchKeyword) {
				if onlyNumber {
					fmt.Printf("%s:slide %s\n", filename, slideNum)
				} else {
					highlighted := highlight(text, keyword, ignoreCase)
					fmt.Printf("%s:slide %s:%s\n", filename, slideNum, strings.TrimSpace(highlighted))
				}
			}
		}
	}
}

func printUsage() {
	fmt.Println("Usage: ooxmlgrep [options] <keyword>")
	fmt.Println("\nOptions:")
	fmt.Println("  -n, --number        Show slide number only (like grep -n)")
	fmt.Println("  -i, --ignore-case   Ignore case distinctions")
	fmt.Println("  --version           Show version and exit")
}

func main() {
	onlyNumber := flag.Bool("n", false, "Show slide number only")
	ignoreCase := flag.Bool("i", false, "Ignore case distinctions")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Usage = printUsage
	flag.Parse()

	for _, arg := range os.Args[1:] {
		if arg == "--number" {
			*onlyNumber = true
		}
		if arg == "--ignore-case" {
			*ignoreCase = true
		}
	}

	if *showVersion {
		fmt.Println(version)
		return
	}

	args := flag.Args()
	if len(args) != 1 {
		printUsage()
		return
	}
	keyword := args[0]

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".pptx") {
			searchInFile(keyword, path, *ignoreCase, *onlyNumber)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
		os.Exit(1)
	}
}
