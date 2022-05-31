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

//RemoveSpaces from an input string
func RemoveSpaces(str string) string {
	return strings.Replace(str, " ", "", -1)
}
