package model

import (
	"fmt"
	"github.com/google/uuid"
	"slices"
)

func NewId() uuid.UUID {
	return uuid.New()
}

func ParseId(value string) uuid.UUID {
	id, _ := uuid.Parse(value)
	return id
}

func Strings[T fmt.Stringer](values []T) []string {
	if values == nil {
		return nil
	}
	strings := make([]string, len(values))
	for i, value := range values {
		strings[i] = value.String()
	}
	return strings
}

func ParseIds(values []string) []uuid.UUID {
	if values == nil {
		return nil
	}
	uuids := make([]uuid.UUID, len(values))
	for i, value := range values {
		id, _ := uuid.Parse(value)
		uuids[i] = id
	}
	return uuids
}

func Union[T ~[]E, E comparable](items T, add ... E) (T, bool) {
	added := false
	for _, value := range add {
		if !slices.Contains(items, value) {
			items = append(items, value)
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