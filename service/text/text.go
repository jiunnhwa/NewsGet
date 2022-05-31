package text

import (
	"strconv"
	"strings"
)

//Enquote an input string
func Enquote(str string) string {
	if strings.HasPrefix(str, "\"") {
		return str
	}
	return strconv.Quote(str)
	//return fmt.Sprintf("%q", str)
}

//RemoveSpaces of an input string
func RemoveSpaces(str string) string {
	return strings.Replace(str, " ", "", -1)
}

//ToLower an input string
func ToLower(str string) string {
	return strings.ToLower(str)
}

//ToUpper an input string
func ToUpper(str string) string {
	return strings.ToUpper(str)
}
