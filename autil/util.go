package autil

import (
	"reflect"

	"github.com/google/uuid"
)

func GenerateUuid() string {
	return uuid.New().String()
}

func GenerateName() string {
	return GenerateUuid()
}

func valueOf(i interface{}) reflect.Value {
	return reflect.ValueOf(i)
}

func In(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

func GetPointer(v interface{}) interface{} {
	vv := valueOf(v)
	if vv.Kind() == reflect.Ptr {
		return v
	}
	return reflect.New(vv.Type()).Interface()
}

func MustBePointer(v interface{}) {
	vv := valueOf(v)
	if vv.Kind() != reflect.Ptr {
		panic("must be a pointer")
	}
}

func Strlen(str string) int {
	return len([]rune(str))
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)

	if rl == 0 {
		return ""
	}

	if start < 0 {
		start = rl + start
	}
	if start < 0 {
		start = 0
	}
	if start > rl-1 {
		return ""
	}

	end := rl

	if length < 0 {
		end = rl + length
	} else if length > 0 {
		end = start + length
	}

	if end < 0 || start >= end {
		return ""
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
