package main

import (
	"fmt"
	"html/template"
	"io"
	"regexp"
)

var tmpl *template.Template
var badChars, _ = regexp.Compile(`[',:]| \(FOIL\)`)
var spaceChars, _ = regexp.Compile(`[ -/]`)
var aChars, _ = regexp.Compile("[\u0061\u24D0\uFF41\u1E9A\u00E0\u00E1\u00E2\u1EA7\u1EA5\u1EAB\u1EA9\u00E3\u0101\u0103\u1EB1\u1EAF\u1EB5\u1EB3\u0227\u01E1\u00E4\u01DF\u1EA3\u00E5\u01FB\u01CE\u0201\u0203\u1EA1\u1EAD\u1EB7\u1E01\u0105\u2C65\u0250]")
var eChars, _ = regexp.Compile("[\u0065\u24D4\uFF45\u00E8\u00E9\u00EA\u1EC1\u1EBF\u1EC5\u1EC3\u1EBD\u0113\u1E15\u1E17\u0115\u0117\u00EB\u1EBB\u011B\u0205\u0207\u1EB9\u1EC7\u0229\u1E1D\u0119\u1E19\u1E1B\u0247\u025B\u01DD]")
var iChars, _ = regexp.Compile("[\u0069\u24D8\uFF49\u00EC\u00ED\u00EE\u0129\u012B\u012D\u00EF\u1E2F\u1EC9\u01D0\u0209\u020B\u1ECB\u012F\u1E2D\u0268\u0131]")
var oChars, _ = regexp.Compile("[\u006F\u24DE\uFF4F\u00F2\u00F3\u00F4\u1ED3\u1ED1\u1ED7\u1ED5\u00F5\u1E4D\u022D\u1E4F\u014D\u1E51\u1E53\u014F\u022F\u0231\u00F6\u022B\u1ECF\u0151\u01D2\u020D\u020F\u01A1\u1EDD\u1EDB\u1EE1\u1EDF\u1EE3\u1ECD\u1ED9\u01EB\u01ED\u00F8\u01FF\u0254\uA74B\uA74D\u0275]")
var uChars, _ = regexp.Compile("[\u0075\u24E4\uFF55\u00F9\u00FA\u00FB\u0169\u1E79\u016B\u1E7B\u016D\u00FC\u01DC\u01D8\u01D6\u01DA\u1EE7\u016F\u0171\u01D4\u0215\u0217\u01B0\u1EEB\u1EE9\u1EEF\u1EED\u1EF1\u1EE5\u1E73\u0173\u1E77\u1E75\u0289]")
var aeChars, _ = regexp.Compile("[\u00c6\u00e6]")

func init() {
	tmpl = template.New("template")
	funcs := map[string]interface{}{
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
	c = aChars.ReplaceAll(c, []byte("a"))
	c = eChars.ReplaceAll(c, []byte("e"))
	//c = iChars.ReplaceAll(c, []byte("i"))
	//c = oChars.ReplaceAll(c, []byte("o"))
	//c = uChars.ReplaceAll(c, []byte("u"))
	c = aeChars.ReplaceAll(c, []byte("ae"))

	if set == "GTC" && (card == "Mountain" || card == "Plains" || card == "Island" || card == "Forest" || card == "Swamp") {
		set = "RTR"
	}
	if (set == "BNG" || set == "JOU") && (card == "Mountain" || card == "Plains" || card == "Island" || card == "Forest" || card == "Swamp") {
		set = "THS"
	}

	url := fmt.Sprintf("http://www.wizards.com/global/images/magic/%s/%s.jpg", set, c)
	return url
}

func makePage(draft *Draft, wr io.Writer) {
	tmpl.Execute(wr, draft)
}
