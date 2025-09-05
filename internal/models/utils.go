package models

import (
	"errors"
	"fmt"
)

func GetFieldFromAny(name string, src any) (any, error) {
	m, ok := src.(map[string]any)
	if !ok {
		return nil, errors.New("src is not a map[string]any")
	}
	if m[name] == nil {
		return nil, errors.New(fmt.Sprintf("no field %s found", name))
	}
	return m[name], nil
}
