package processor

import (
	"reflect"
	"testing"
)

func TestProcessNumbers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		nums       []int
		processors []func([]int) []int
		expected   []int
	}{
		{
			name:       "No processors",
			nums:       []int{5, 3, 9, 1, 3, 7, 5},
			processors: nil,
			expected:   []int{5, 3, 9, 1, 7},
		},
		{
			name: "Remove duplicates only",
			nums: []int{5, 3, 9, 1, 3, 7, 5},
			processors: []func([]int) []int{
				removeDuplicates,
			},
			expected: []int{1, 3, 5, 7, 9},
		},
		{
			name: "Sort only",
			nums: []int{5, 3, 9, 1, 7},
			processors: []func([]int) []int{
				sort,
			},
			expected: []int{1, 3, 5, 7, 9},
		},
		{
			name: "Floor only",
			nums: []int{5, 3, 9, 1, -2, 7},
			processors: []func([]int) []int{
				floor(3),
			},
			expected: []int{3, 5, 7, 9},
		},
		{
			name: "Floor and remove duplicates",
			nums: []int{5, 3, 9, 3, -2, 7, 9},
			processors: []func([]int) []int{
				floor(4),
				removeDuplicates,
			},
			expected: []int{5, 7, 9},
		},
		{
			name: "All processors",
			nums: []int{5, 3, 9, 1, -2, 7, 3, 9, 0},
			processors: []func([]int) []int{
				floor(3),
				removeDuplicates,
				sort,
			},
			expected: []int{3, 5, 7, 9},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()
				result := processNumbers(test.nums, test.processors...)
				if !reflect.DeepEqual(result, test.expected) {
					t.Errorf("expected %v, got %v", test.expected, result)
				}
			},
		)
	}
}
