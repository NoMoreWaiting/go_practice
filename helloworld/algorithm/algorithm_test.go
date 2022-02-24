package algorithm

import "testing"

func Test_SingleNonDuplicate(t *testing.T) {
	type args struct {
		nums []int
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{"case1", args{[]int{3, 3, 7, 7, 10, 11, 11, 12, 12, 13, 13}}, 10},
		{"case2", args{[]int{3, 7, 7, 10, 10, 11, 11, 12, 12, 13, 13}}, 3},
		{"case3", args{[]int{3, 3, 7, 7, 11, 11, 12, 12, 13, 13, 14}}, 14},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SingleNonDuplicate(tt.args.nums); got != tt.want {
				t.Errorf("singleNonDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}
