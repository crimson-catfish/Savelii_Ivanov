package processer

import (
	"slices"
)

func sort(nums []int) []int {
	n := make([]int, len(nums))
	copy(n, nums)
	slices.Sort(n)

	return n
}

func removeDuplicates(nums []int) []int {
	n := make([]int, 0)
	appear := make(map[int]bool)
	for _, num := range nums {
		if appear[num] {
			continue
		}
		appear[num] = true
		n = append(n, num)
	}
	return n
}

func floor(min int) func(nums []int) []int {
	return func(nums []int) []int {
		n := make([]int, 0)
		for i := 0; i < len(nums); i++ {
			if nums[i] >= min {
				n = append(n, nums[i])
			}
		}
		return n
	}
}

func processNumbers(nums []int, processors ...func([]int) []int) []int {
	n := make([]int, len(nums))
	copy(n, nums)
	for _, process := range processors {
		n = process(n)
	}

	n = floor(0)(n)
	n = sort(n)
	n = removeDuplicates(n)

	return n
}
