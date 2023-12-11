package model

import (
	"fmt"
	"github.com/google/uuid"
	"slices"
)

func NewId() uuid.UUID {
	return uuid.New()
}

func Strings[T fmt.Stringer](values []T) []string {
	strings := make([]string, len(values))
	for i, value := range values {
		strings[i] = value.String()
	}
	return strings
}

func Union[T ~[]E, E comparable](items T, add ... E) (T, bool) {
	added := false
	for _, value := range items {
		if !slices.Contains(add, value) {
			add = append(add, value)
			added = true
		}
	}
	return add, added
}

func Difference[T ~[]E, E comparable](items T, remove ...E) (T, bool) {
	removed := false
	for _, item := range remove {
		if len(items) == 0 {
			return items, false
		}
		for ii, s := range items {
			if s == item {
				items[ii] = items[len(items)-1]
				items = items[:len(items)-1]
				removed = true
				break
			}
		}
	}
	return items, removed
}