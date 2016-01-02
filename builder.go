package main

import (
	"io"
	"html/template"
	"fmt"
	"regexp"
)

var tmpl *template.Template
var badChars, _ = regexp.Compile(`[']| \(FOIL\)`)
var spaceChars, _ = regexp.Compile(`[ -]`)

func init() {
	tmpl = template.New("template")
	funcs := map[string] interface{} {
		"wizUrl": wizardsUrl,
	}
	tmpl = tmpl.Funcs(funcs)
	var err error
	tmpl, err = tmpl.ParseFiles("template")
	if err != nil {
		fmt.Println("template parsing error", err)
	}
}

func wizardsUrl(set, card string) string {
	c := []byte(card)
	c = badChars.ReplaceAll(c, []byte{})	
	c = spaceChars.ReplaceAll(c, []byte("_"))
	url := fmt.Sprintf("http://www.wizards.com/global/images/magic/%s/%s.jpg", set, c)
	return url
}

func makePage(draft *Draft, wr io.Writer) {
	tmpl.Execute(wr, draft)
}
