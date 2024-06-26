package utils

import (
	cryptoRand "crypto/rand"
	"fmt"
	mathRand "math/rand"
	"regexp"
	"strings"
	"time"
)

func StrRemoveSpace(str string) string {
	return regexp.MustCompile(`\s`).ReplaceAllString(str, "")
}

func StrAddBegin(str, subStr string) string {
	if str == "" {
		return str
	}

	if str[:len(subStr)] != subStr {
		return subStr + str
	}

	return str
}

func StrRemoveBegin(str, subStr string) string {
	if str == "" {
		return str
	}

	subStrLength := len(subStr)
	if str[:subStrLength] == subStr {
		return str[subStrLength:]
	}

	return str
}

func StrAddEnd(str, subStr string) string {
	if str == "" {
		return str
	}

	lastMatchedIndex := strings.LastIndex(str, subStr)
	isAtEnd := lastMatchedIndex+len(subStr) == len(str)
	if !isAtEnd {
		return str + subStr
	}

	return str
}

func StrRemoveEnd(str, subStr string) string {
	if str == "" {
		return str
	}

	lastMatchedIndex := strings.LastIndex(str, subStr)
	isAtEnd := lastMatchedIndex+len(subStr) == len(str)
	if isAtEnd {
		return str[:lastMatchedIndex]
	}

	return str
}

func StrWithCharset(length int, charset string) string {
	b := make([]byte, length)
	seededRand := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func StrRandom(length int) string {
	charset := "`~1!2@3#4$5%6^7&8*9(0)-_=+qQwWeErRtTyYuUiIoOpP[{]}\\|aAsSdDfFgGhHjJkKlL;:'zZxXcCvVbBnNmM,<.>/?"
	return StrWithCharset(length, charset)
}

func StrSegment(str string, sep byte, start int) (string, int) {
	if len(str) == 0 || start < 0 || start > len(str)-1 {
		return "", -1
	}

	i := strings.IndexByte(str[start+1:], sep)
	if i < 0 {
		return str[start:], i
	}

	next := i + start + 1
	return str[start+1 : next], next
}

func StrRemoveDup(str, pattern string) string {
	resp := ""
	for i, r := range str {
		if i != 0 {
			if string(r) == pattern && str[i-1:i] == pattern {
				continue
			}
		}
		resp += string(r)
	}

	return resp
}

func StrIsLower(str string) []bool {
	res := []bool{}
	for _, r := range str {
		s := string(r)
		res = append(res, s == strings.ToLower(s))
	}

	return res
}

func StrIsUpper(str string) []bool {
	res := []bool{}
	for _, r := range str {
		s := string(r)
		res = append(res, s == strings.ToUpper(s))
	}

	return res
}

func StrUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := cryptoRand.Read(uuid)
	if err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func ErrorMessage(actual, expected any, desc string) string {
	if desc != "" {
		desc = FmtGreen("\n%v", desc)
	}

	return desc +
		FmtDim("\nExpected: %v", FmtRed("%v", expected)) +
		FmtDim("\nActual: %v", FmtRed("%v", actual))
}

func StrCapitalizeFirstLetter(str string) string {
	if len(str) == 0 {
		return str
	}
	firstLetter := string(str[0] - 32)

	return firstLetter + str[1:]
}
