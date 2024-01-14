package validation

import (
	"net/mail"
	"strings"
	"unicode"
)

func IsValidEmail(s string) bool {
	if _, err := mail.ParseAddress(s); err != nil {
		return false
	}
	return true
}

func IsValidPassword(s string) bool {
	letters := 0
	var sevenOrMore, number, upper, special bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		}
	}
	sevenOrMore = letters >= 8 && letters <= 64
	if sevenOrMore && number && upper && special {
		return true
	}
	return false
}

func IsValidUsername(s string) bool {
	if len(s) > 0 && len(s) <= 20 && !strings.Contains(s, ":") {
		return true
	}
	return false
}
