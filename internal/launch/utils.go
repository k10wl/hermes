package launch

import (
	"strings"
)

func countTruthyValues(vals ...any) int {
	truthyValues := 0
	for _, val := range vals {
		switch val.(type) {
		case string:
			if strings.Trim(val.(string), " ") != "" {
				truthyValues++
			}
		case bool:
			if val.(bool) {
				truthyValues++
			}
		default:
			panic("unhandled truthy data type")
		}
	}
	return truthyValues
}
