package keyboard

import (
	"fmt"
	"strings"
)

type Layout struct {
	Size string
	Rows [][]string
	Keys map[string]bool
}

func Lookup(size string) (Layout, error) {
	size = strings.TrimSuffix(size, "%")
	rows, ok := presets[size]
	if !ok {
		return Layout{}, fmt.Errorf("unsupported keyboard size %q; choose %s", size, strings.Join(Sizes, ", "))
	}
	layout := Layout{Size: size, Rows: rows, Keys: make(map[string]bool)}
	for _, row := range rows {
		for _, key := range row {
			layout.Keys[key] = true
		}
	}
	return layout, nil
}
