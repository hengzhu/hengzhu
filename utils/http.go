package utils

import (
	"strings"
	"net/url"
)

func QueryString2Map(que string) (set map[string]string) {
	set = map[string]string{}
	if !strings.Contains(que, "&") {
		return
	}
	for _, kv := range strings.Split(que, "&") {
		kAv := strings.Split(kv, "=")
		if len(kAv) == 2 {
			k, err := url.QueryUnescape(kAv[0])
			v, err2 := url.QueryUnescape(kAv[1])
			if err == nil && err2 == nil {
				set[k] = v
			}
		}
	}
	return
}
