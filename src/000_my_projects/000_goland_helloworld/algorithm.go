package main

import "fmt"

func testSingleNonDuplicate() {
	nums := []int{3, 3, 7, 7, 10, 11, 11, 12, 12, 13, 13}
	fmt.Println(singleNonDuplicate(nums))

}

func singleNonDuplicate(nums []int) int {
	fmt.Println(nums)
	len := len(nums)
	if 1 == len {
		return nums[0]
	} else if len%2 == 0 {
		return -1 // 输入有误
	}

	if mid := len / 2; mid%2 == 0 {
		if nums[mid] == nums[mid+1] {
			nums = nums[mid+2:]
		} else if nums[mid] == nums[mid-1] {
			nums = nums[:mid-1]
		} else {
			return nums[mid]
		}
	} else {
		if nums[mid] == nums[mid+1] {
			nums = nums[:mid]
		} else if nums[mid] == nums[mid-1] {
			nums = nums[mid+1:]
		} else {
			return nums[mid]
		}
	}
	return singleNonDuplicate(nums)
}
