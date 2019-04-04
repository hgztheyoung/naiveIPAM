package naiveIPAM

import (
	"fmt"
)

// a range [Low,High)
type NumRange struct {
	Low  uint32
	High uint32
}

func (numRange *NumRange) Len() uint32 {
	return numRange.High - numRange.Low
}

func (numRange *NumRange) ContainsNumber(n uint32) bool {
	return n >= numRange.Low &&
		n <= numRange.High
}

func (numRange *NumRange) GetHighLow() uint64 {
	return uint64(numRange.High)<<32 | uint64(numRange.Low)
}

func EqualNumRange(n1, n2 NumRange) bool {
	return n1.GetHighLow() == n2.GetHighLow()
}

func (numRange *NumRange) Cmp(n2 *NumRange) int {
	hl1 := numRange.GetHighLow()
	hl2 := n2.GetHighLow()
	if hl1 > hl2 {
		return 1
	}
	if hl1 < hl2 {
		return -1
	}
	return 0
}

func NumRangeIntersect(n1, n2 NumRange) *NumRange {
	ret := new(NumRange)
	if n1.Low >= n2.High || n1.High <= n2.Low {
		return ret
	}
	ret.High = n1.High
	if ret.High > n2.High {
		ret.High = n2.High
	}
	ret.Low = n1.Low
	if ret.Low < n2.Low {
		ret.Low = n2.Low
	}
	return ret
}

type NumRangeSlice []NumRange

func (p NumRangeSlice) Len() int { return len(p) }

func (p NumRangeSlice) Less(i, j int) bool { return p[i].GetHighLow() < p[j].GetHighLow() }

func (p NumRangeSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (numRange NumRange) String() string {
	return fmt.Sprintf("[%v,%v)", numRange.Low, numRange.High)
}
