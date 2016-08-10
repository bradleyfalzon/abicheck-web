package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/bradleyfalzon/abicheck"
)

var tmpls *template.Template // tmpls contains all the html templates

func main() {
	log.Println("Starting...")

	listen := flag.String("listen", ":80", "Listen address")
	flag.Parse()

	// initialise html templates
	var err error
	log.Println("Parsing templates...")
	if tmpls, err = template.ParseGlob("*.tmpl"); err != nil {
		log.Fatalf("could not parse html templates: %s", err)
	}

	http.HandleFunc("/", homeHandler)

	log.Println("Listening on:", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var (
		prev string = r.Form.Get("prev")
		next string = r.Form.Get("next")
		buf  bytes.Buffer
	)

	if prev != "" {
		// check it
		vcs := abicheck.StrVCS{}
		vcs.SetFile("rev1", "example.go", []byte(prev))
		vcs.SetFile("rev2", "example.go", []byte(next))

		abi := abicheck.New(abicheck.SetVLog(&buf), abicheck.SetVCS(vcs))
		changes, err := abi.Check("", "rev1", "rev2")

		if err != nil {
			fmt.Fprint(&buf, err)
			fmt.Println(err)
		}

		for _, change := range changes {
			fmt.Fprint(&buf, change)
		}
	}

	page := struct {
		Prev    string
		Next    string
		Results string
	}{prev, next, buf.String()}

	if err := tmpls.ExecuteTemplate(w, "home.tmpl", page); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing home template: %s", err)
	}
}
