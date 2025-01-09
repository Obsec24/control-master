package traffic

import (
	"bytes"
	. "github.com/Obsec24/control-master/common"
	"strings"
)

func ApkName() (sol string) {
	var out bytes.Buffer
	Command(&out, "aapt", "dump", "badging", "base.apk")
	sol = SplitSearch(out.String(), "\n", "package")
	sol = SplitSearch(sol, " ", "name")
	sol = strings.Trim(strings.Split(sol, "=")[1], "'")
	return
}

func SplitSearch(s, sep, query string) string {
	arr := strings.Split(s, sep)
	for i := 0; i < len(arr); i++ {
		if strings.Contains(arr[i], query) {
			return arr[i]
		}
	}
	return ""
}
