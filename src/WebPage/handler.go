package WebPage

import (
	_ "embed"
	"strconv"
	"strings"
)

//go:embed Home/Home.html
var homeHtml string

//go:embed Home/Home.js
var homeJS string

//go:embed Home/Home.css
var homeCSS string

func GetHomeHtml(port int32) string {
	home_app := string(homeHtml) + "\n" + "<script>" + string(homeJS) + "</script>" + "\n" + "<style>" + string(homeCSS) + "</style>"

	return strings.Replace(home_app, "IP_ADDRESS", "localhost:"+strconv.Itoa(int(port)), -1)
}
