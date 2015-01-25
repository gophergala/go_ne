package shared

import "fmt"

func ExtractOptions(p interface{}) []string {
	var result []string
	raw := p.([]interface{})

	for _, v := range raw {
		result = append(result, fmt.Sprintf("%v", v))
	}

	return result
}

func ExtractBool(p interface{}) bool {
	return p.(bool)
}

func ExtractString(p interface{}) string {
	return p.(string)
}
