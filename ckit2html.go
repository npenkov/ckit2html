package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/html"
)

var (
	inFolder, outFolder string
	verboseOutput       bool
	showVersion         bool
)

var (
	version string
)

func main() {
	flag.StringVar(&inFolder, "in", ".", "Input folder")
	flag.StringVar(&outFolder, "out", ".", "Output folder")
	flag.BoolVar(&verboseOutput, "v", false, "Set to verbose output")
	flag.BoolVar(&showVersion, "version", false, "Show the version")
	flag.Parse()

	log.SetLevel(log.InfoLevel)
	if verboseOutput {
		log.SetLevel(log.TraceLevel)
	}
	if showVersion {
		fmt.Fprintf(os.Stderr, "Running ckit2html version %s\n", version)
		os.Exit(0)
	}

	log.Debug("Transforming .kit files from folder " + inFolder + " to " + outFolder)

	files := []string{}
	err := filepath.Walk(inFolder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".kit") {
				files = append(files, path)

			}
			return nil
		})
	if err != nil {
		log.Error(err)
	}
	for _, p := range files {
		parseFile(p)
	}
}

func parseFile(path string) {
	// If filename starts with `_` then we don't write it
	if strings.HasPrefix(filepath.Base(path), "_") {
		return
	}
	l := log.WithField("file", path)
	l.Debug("Parsing")
	r, err := os.Open(path)
	if err != nil {
		l.WithError(err).Error("Error loading file")
		return
	}
	defer r.Close()

	replacedContents := processFileImports(r, path, inFolder)

	outputFileName := strings.ReplaceAll(strings.ReplaceAll(path, ".kit", ".html"), inFolder, outFolder)
	l2 := l.WithField("out_file", outputFileName)

	outputFile, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		l2.WithError(err).Error("Cannot create output file")
		return
	}
	defer outputFile.Close()

	outputFile.Write(replaceVariables(replacedContents.Bytes()))
	l2.Trace("Output file written")
}

func processFileImports(r io.Reader, path, basePath string) *bytes.Buffer {
	l := log.WithField("base_path", basePath).WithField("file", path)

	fileContent := bytes.NewBuffer([]byte{})
	z := html.NewTokenizer(r)

forLoop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err().Error() != io.EOF.Error() {
				l.WithError(z.Err()).Error("Error parsing file")
			}
			break forLoop
		case html.CommentToken:
			tn := string(z.Text())
			if strings.HasPrefix(tn, " @import ") || strings.HasPrefix(tn, "@include ") {
				res := processImport(basePath, tn)
				fileContent.Write(res)
			} else {
				fileContent.Write(z.Raw())
			}
		default:
			fileContent.Write(z.Raw())
		}
	}
	return fileContent
}
