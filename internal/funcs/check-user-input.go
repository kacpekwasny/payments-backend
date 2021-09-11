package funcs

import cmt "github.com/kacpekwasny/commontools"

var allowedChars = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "-", "_"}

func NoForbidenChars(str string) bool {
	for _, char := range str {
		_, found := cmt.InSlice(string(char), allowedChars)
		if !found {
			return false
		}
	}
	return true
}
