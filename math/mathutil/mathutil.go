package mathutil

import (
	"fmt"
	"strconv"
)

type Comparable interface {
	Less(Comparable) bool
	Sub(Comparable) Comparable
}

type Int64 int64

func (v Int64) Less(other Comparable) bool {
	return v < other.(Int64)
}
func (v Int64) Sub(other Comparable) Comparable {
	return v - other.(Int64)
}

const defaultPercentBase = 10000

type Percent struct {
	numerator   int64
	denominator int64
}

func NewPercent(x, y int64) Percent {
	if y < 0 {
		x, y = -x, -y
	}
	return Percent{numerator: x, denominator: y}
}

func (p Percent) Normalize() int64 {
	return p.NormalizeWithBase(defaultPercentBase)
}

func (p Percent) NormalizeWithBase(base int64) int64 {
	if p.denominator == 0 {
		return 0
	}
	x1, y1 := p.numerator/p.denominator, p.numerator%p.denominator
	x2, y2 := base/p.denominator, base%p.denominator
	return x1*x2*p.denominator + (x1*y2 + x2*y1) + (y1*y2)/p.denominator
}

func (p Percent) Less(other Comparable) bool {
	p2 := other.(Percent)
	x1 := float64(0)
	if p.denominator != 0 {
		x1 = float64(p.numerator) / float64(p.denominator)
	}
	x2 := float64(0)
	if p2.denominator != 0 {
		x2 = float64(p2.numerator) / float64(p2.denominator)
	}
	return x1 < x2
}

func (p Percent) Sub(other Comparable) Comparable {
	p2 := other.(Percent)
	x1 := float64(0)
	if p.denominator != 0 {
		x1 = float64(p.numerator) / float64(p.denominator)
	}
	x2 := float64(0)
	if p2.denominator != 0 {
		x2 = float64(p2.numerator) / float64(p2.denominator)
	}
	x := x1 - x2
	return NewPercent(int64(x*defaultPercentBase), defaultPercentBase)
}

func (p Percent) String() string {
	value := p.NormalizeWithBase(10000)
	abs := value
	sign := ""
	if abs < 0 {
		abs = -value
		sign = "-"
	}
	return fmt.Sprintf("%s%d.%02d%%", sign, abs/100, abs%100)
}

func (p Percent) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(p.Normalize(), 10)), nil
}

func (p *Percent) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	p.numerator = i
	p.denominator = defaultPercentBase
	return nil
}

// C(m, n)
func Comb(m, n int) int64 {
	sum := int64(1)
	if m > n*2 {
		n = m - n
	}
	for i := 0; i < n; i++ {
		sum *= int64(m - i)
	}
	for i := 0; i < n; i++ {
		sum /= int64(i + 1)
	}
	return sum
}

// CombSet C(m, n)
func CombSet(m, n int) [][]int {
	type node struct {
		parent *node
		value  int
		depth  int
	}

	var (
		result [][]int
		nodes  []*node
	)
	root := new(node)
	root.value = -1
	nodes = append(nodes, root)
	for len(nodes) > 0 {
		x := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		if x.depth == n {
			values := make([]int, 0, n)
			curr := x
			for curr != nil && curr.value >= 0 {
				values = append(values, curr.value)
				curr = curr.parent
			}
			result = append(result, values)
			continue
		}

		for value := x.value + 1; value < m+x.depth+1-n; value++ {
			newNode := new(node)
			newNode.depth = x.depth + 1
			newNode.parent = x
			newNode.value = value
			nodes = append(nodes, newNode)
		}
	}

	return result
}

func MultiCombSet(nums []int, n int) [][]int {
	remainSum := 0
	for _, x := range nums {
		remainSum += x
	}
	return multiCombSet(nums, 0, n, remainSum)
}

func multiCombSet(nums []int, begin, n, remainSum int) [][]int {
	size := len(nums)
	if n == 0 {
		return [][]int{make([]int, size)}
	}
	if n > remainSum {
		return nil
	}
	if remainSum == n {
		values := make([]int, size)
		for i := begin; i < size; i++ {
			values[i] = nums[i]
		}
		return [][]int{values}
	}
	var result [][]int
	for x := 0; x <= nums[begin]; x++ {
		if n >= x && begin+1 < size {
			tmpResult := multiCombSet(nums, begin+1, n-x, remainSum-nums[begin])
			for _, tmp := range tmpResult {
				tmp[begin] = x
				result = append(result, tmp)
			}
		}
	}
	return result
}
