package utils

import "regexp"

func ReplaceWholeWord(originalString string, oldWord string, newWord string) string {
	re := regexp.MustCompile(`\b` + oldWord + `\b`)
	return re.ReplaceAllString(originalString, newWord)
}
