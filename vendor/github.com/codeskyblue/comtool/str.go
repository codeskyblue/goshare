package comtool

import "strings"

func Template(s string, m map[string]string) string {
	for k, v := range m {
		s = strings.Replace(s, "{"+k+"}", v, -1)
	}
	return s
}
