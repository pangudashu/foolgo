package foolgo

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var rewrite_regexp []*Rewrite
var rewrite_static map[string]string

type Rewrite struct {
	pattern string
	match   string
	regex   *regexp.Regexp
}

func regRewrite(list map[string]string) {
	for p, m := range list {
		if strings.Index(p, "(") < 0 {
			if rewrite_static == nil {
				rewrite_static = make(map[string]string)
			}
			rewrite_static[p] = m
			continue
		}
		r := regexp.MustCompile(p)
		if r == nil {
			continue
		}
		reg := &Rewrite{
			pattern: p,
			match:   m,
			regex:   r,
		}
		rewrite_regexp = append(rewrite_regexp, reg)
	}
}

func matchRewrite(r *http.Request) {
	url := r.URL.Path
	var rewrite_url string = ""
	var ok bool

	if rewrite_url, ok = rewrite_static[url]; ok == true {
		goto RESET_URI
	}

	for _, rewrite := range rewrite_regexp {
		match := rewrite.regex.FindAllStringSubmatch(url, -1)
		if match == nil {
			continue
		}
		match_cnt := len(match[0])
		if match_cnt == 1 {
			return
		}

		rewrite_url = rewrite.match

		for n := 1; n < match_cnt; n++ {
			replace_val := "[" + strconv.Itoa(n) + "]"
			rewrite_url = strings.Replace(rewrite_url, replace_val, match[0][n], -1)
		}
		break
	}
	if rewrite_url == "" {
		return
	}

RESET_URI:

	rewrite_url = strings.Replace(rewrite_url, "[args]", r.URL.RawQuery, -1)
	uri_map := strings.SplitN(rewrite_url, "?", 2)

	if len(uri_map) == 2 {
		r.URL.Path = uri_map[0]
		r.URL.RawQuery = uri_map[1]
	} else {
		r.URL.Path = uri_map[0]
	}
}
