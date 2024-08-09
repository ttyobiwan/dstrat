package sqlite

import "strings"

func ParseGroupConcat[T any](value, rowSep, valueSep string, parser func(v []string) T) []T {
	rows := strings.Split(value, rowSep)
	dest := make([]T, len(rows))
	for idx, row := range rows {
		dest[idx] = parser(strings.Split(row, valueSep))
	}
	return dest
}
