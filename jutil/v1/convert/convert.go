package convert

import "strconv"

//ToInt converts string to integer
func ToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return i
}
