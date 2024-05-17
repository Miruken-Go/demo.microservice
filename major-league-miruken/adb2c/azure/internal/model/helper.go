package model

import (
	"slices"

	"github.com/google/uuid"
)

func NewId() string {
	return uuid.New().String()
}

func Union[T ~[]E, E comparable](items T, add ...E) (T, bool) {
	added := false
	for _, item := range add {
		if !slices.Contains(items, item) {
			items = append(items, item)
			added = true
		}
	}
	return items, added
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

