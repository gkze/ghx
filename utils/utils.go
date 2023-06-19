package utils

import (
	"reflect"
	"sort"
	"strings"
)

var DefaultStructFieldNameMatcher = func(fieldName string) func(s string) bool {
	return func(s string) bool { return strings.EqualFold(fieldName, s) }
}

type SortSliceOpts[S any] struct {
	Slice                  []S
	SortField              string
	StructFieldNameMatcher func(string) func(string) bool
}

func SortSlice[S any](opts *SortSliceOpts[S]) {
	if opts.StructFieldNameMatcher == nil {
		opts.StructFieldNameMatcher = DefaultStructFieldNameMatcher
	}

	sort.Slice(opts.Slice, func(i, j int) bool {
		fieldI := reflect.ValueOf(opts.Slice[i]).FieldByNameFunc(
			opts.StructFieldNameMatcher(opts.SortField),
		)
		fieldJ := reflect.ValueOf(opts.Slice[j]).FieldByNameFunc(
			opts.StructFieldNameMatcher(opts.SortField),
		)

		switch fieldI.Kind() {
		case reflect.Uint64:
			return fieldI.Uint() < fieldJ.Uint()
		case reflect.String:
			return fieldI.String() < fieldJ.String()
		case reflect.SliceOf(reflect.TypeOf(reflect.String)).Kind():
			return strings.Compare(
				sort.StringSlice(fieldI.Interface().([]string))[0],
				sort.StringSlice(fieldJ.Interface().([]string))[0],
			) == -1
		}

		return false
	})
}
