package funcs

import (
	"fmt"
	"net/http"
	"strings"

	cmt "github.com/kacpekwasny/commontools"
)

var defaultLang = "en"
var allowedLangs = []string{"pl", "en"}

func GetLang(r *http.Request) string {
	ck, err := r.Cookie("lang")
	if err != nil {
		return defaultLang
	}
	lang := strings.Split(ck.String(), "=")[1]
	if _, is := cmt.InSlice(lang, allowedLangs); is {
		return lang
	}
	return defaultLang
}

func Cookie2Str(r *http.Request) string {
	ret := ""
	for _, c := range r.Cookies() {
		var val = c.Value
		if len(c.Value) > 10 && c.Name != "login" {
			val = c.Value[:10] + "..."
		}
		ret += fmt.Sprintf("%v=%v;  ", c.Name, val)
	}
	return ret
}
