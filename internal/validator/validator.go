package validator

import (
	"context"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator interface {
	Valid(context.Context) Evaluator
}

var EmailRx = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Evaluator map[string]string

func (e *Evaluator) AddFieldError(key, message string) {
	if *e == nil {
		*e = make(map[string]string)
	}

	// acessar maps em ponteiros
	if _, exists := (*e)[key]; !exists {
		(*e)[key] = message
	}
}

func (e *Evaluator) CheckField(ok bool, key, message string) {
	if !ok {
		e.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
