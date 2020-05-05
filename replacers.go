package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/html"
)

func processImport(basePath, text string) []byte {
	l := log.WithField("base_path", basePath).WithField("import_string", text)
	l.Debug("Processing imports")

	// Make a Regex to say we only want letters and numbers
	procString := strings.TrimLeft(text, " @import ")
	procString = strings.TrimLeft(procString, " @include ")

	// We strip off all quotes
	reg, err := regexp.Compile(`["']+`)
	if err != nil {
		log.Fatal(err)
	}
	files := reg.ReplaceAllString(procString, "")
	allFiles := strings.Split(files, " ")

	buf := bytes.NewBuffer([]byte{})

	for _, file := range allFiles {
		if file == "" {
			continue
		}
		fileNameWithPath := filepath.Join(basePath, file)
		l2 := l.WithField("import_file", fileNameWithPath)

		l2.Debug("Importing file")

		impFile, err := os.Open(fileNameWithPath)
		if err != nil {
			nf := filepath.Join(filepath.Dir(fileNameWithPath), "_"+filepath.Base(fileNameWithPath))
			impFile, err = os.Open(nf)
			if err != nil {
				l2.WithError(err).Error("Error opening import file")
				return []byte{}
			}
		}

		data, err := ioutil.ReadAll(impFile)
		if err != nil {
			l2.WithError(err).Error("Reading import file")
			return []byte{}
		}
		buf.Write(processFileImports(bytes.NewBuffer(data), fileNameWithPath, filepath.Dir(fileNameWithPath)).Bytes())
		l2.Trace("Importing file processed")
	}

	return buf.Bytes()
}

func replaceVariables(data []byte) []byte {
	replaced := bytes.NewBuffer([]byte{})
	buf := bytes.NewBuffer(data)
	z := html.NewTokenizer(buf)
	varDef, _ := regexp.Compile(`(\$[a-zA-Z0-9-]+)\:(.*)`)
	// varRepl, _ := regexp.Compile(`<!--[ \n]*(\$[a-zA-Z0-9-]+)[ ]*-->`)

	vars := map[string][]byte{}

forLoop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err().Error() != io.EOF.Error() {
				log.WithField("input", string(data)).WithError(z.Err()).Error("Error parsing buffer")
			}
			break forLoop
		case html.CommentToken:
			comment := z.Text()
			comment = []byte(strings.ReplaceAll(string(comment), "\n", ""))
			if varDef.Match(comment) {
				varKey := varDef.ReplaceAll(comment, []byte("$1"))
				varValue := varDef.ReplaceAll(comment, []byte("$2"))
				trimmedKey := strings.Trim(string(varKey), " ")
				log.WithField("name", trimmedKey).Trace("Defining variable")
				vars[trimmedKey] = []byte(strings.Trim(string(varValue), " "))
			} else {
				log.WithField("data", string(z.Raw())).Trace("No replacement found")
				replaced.Write(z.Raw())
			}
		default:
			replaced.Write(z.Raw())
		}
	}

	rplst := replaced.String()
	for k, v := range vars {
		tag := fmt.Sprintf("<!-- %s -->", k)
		rplst = strings.ReplaceAll(rplst, tag, string(v))
		log.WithField("variable", tag).Trace("Replaced variable")
	}

	return []byte(rplst)
}
