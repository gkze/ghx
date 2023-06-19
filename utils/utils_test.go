package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type sortee struct {
	ID     uint64
	Name   string
	Values []string
}

func TestSortSliceBy(t *testing.T) {
	for _, tc := range []struct {
		input    []sortee
		expected []sortee
		sortKey  string
	}{
		{
			[]sortee{{2, "", []string{}}, {1, "", []string{}}},
			[]sortee{{1, "", []string{}}, {2, "", []string{}}},
			"ID",
		},
		{
			[]sortee{{0, "name", []string{}}, {0, "anothername", []string{}}},
			[]sortee{{0, "anothername", []string{}}, {0, "name", []string{}}},
			"name",
		},
		{
			[]sortee{{0, "", []string{"hi", "there"}}, {0, "", []string{"check", "ing"}}},
			[]sortee{{0, "", []string{"check", "ing"}}, {0, "", []string{"hi", "there"}}},
			"values",
		},
	} {
		SortSlice[sortee](&SortSliceOpts[sortee]{tc.input, tc.sortKey, nil})
		assert.Equal(t, tc.expected, tc.input)
	}
}
