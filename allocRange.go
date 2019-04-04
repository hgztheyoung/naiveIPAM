package naiveIPAM

import (
	"sort"
	"fmt"
	"google/uuid"
	"sync"
)

var RangeAllocRangeMap = make(map[uint64]*AllocRange)
var RangeTagMap = make(map[uint64]interface{})
var uuidRangeMap = make(map[string]*AllocRange)
var uuidRangeMapMutex sync.RWMutex

var allocMutex sync.Mutex

type AllocRange struct {
	InitialAllocRange NumRange
	AllocRanges       []NumRange
	UsedRanges        []NumRange
}

func InitAllocRange(nr NumRange) string {
	ret := new(AllocRange)
	ret.InitialAllocRange = nr
	ret.AllocRanges = cutNumRangeIntoCidrPieces(nr, 0)
	uid := uuid.New().String()
	uuidRangeMapMutex.RLock()
	uuidRangeMap[uid] = ret
	uuidRangeMapMutex.RUnlock()
	RangeAllocRangeMap[ret.InitialAllocRange.GetHighLow()] = ret
	return uid
}

func (allocRange *AllocRange) AddToUsedRange(nr NumRange) {
	allocRange.UsedRanges = append(allocRange.UsedRanges, nr)
	sort.Sort(NumRangeSlice(allocRange.UsedRanges))
}

func (allocRange *AllocRange) AllocNIPs(n uint32) NumRange {
	allocMutex.Lock()
	defer allocMutex.Unlock()
	n = ceiling(n)
	for i, r := range allocRange.AllocRanges {
		if r.Len() >= n {
			newalloc := make([]NumRange, 0)
			newalloc = append(newalloc, allocRange.AllocRanges[:i]...)
			newalloc = append(newalloc, cutNumRangeIntoCidrPieces(NumRange{Low: r.Low + n, High: r.High}, 0)...)
			newalloc = append(newalloc, allocRange.AllocRanges[i+1:]...)
			allocRange.AllocRanges = newalloc
			//allocRange.AllocRanges = append(allocRange.AllocRanges, allocRange.AllocRanges[i+1:]...)

			ret := NumRange{Low: r.Low, High: r.Low + n}
			//allocRange.UsedRanges = append(allocRange.UsedRanges, ret)
			allocRange.AddToUsedRange(ret)
			InitAllocRange(ret)
			return ret
		}
	}
	return NumRange{}
}

func (allocRange *AllocRange) AllocSpecifiedCidrNumRange(nr NumRange) NumRange {
	allocMutex.Lock()
	defer allocMutex.Unlock()
	for i, a := range allocRange.AllocRanges {
		if a.Low <= nr.Low && a.High >= nr.High {
			//allocRange.UsedRanges = append(allocRange.UsedRanges, nr)
			allocRange.AddToUsedRange(nr)
			temp := a
			allocRange.AllocRanges = append(allocRange.AllocRanges[:i], allocRange.AllocRanges[i+1:]...)
			if temp.Low < nr.Low {
				allocRange.AllocRanges = append(allocRange.AllocRanges, NumRange{temp.Low, nr.Low})
			}
			if temp.High > nr.High {
				allocRange.AllocRanges = append(allocRange.AllocRanges, NumRange{nr.High, temp.High})
			}
			sort.Sort(NumRangeSlice(allocRange.AllocRanges))
			return nr
		}
	}
	return NumRange{}
}

func (allocRange *AllocRange) DeAllocNumRange(nr NumRange) {
	for i, u := range allocRange.UsedRanges {
		if EqualNumRange(u, nr) {
			allocRange.UsedRanges = append(allocRange.UsedRanges[:i], allocRange.UsedRanges[i+1:]...)
			allocRange.AllocRanges = append(allocRange.AllocRanges, u)
			sort.Sort(NumRangeSlice(allocRange.AllocRanges))
		}
		temp := make([]NumRange, 0)
		if len(allocRange.AllocRanges) >= 1 {
			f := allocRange.AllocRanges[0].Low
			l := allocRange.AllocRanges[0].High
			for _, a := range allocRange.AllocRanges[1:] {
				if a.Low == l {
					l = a.High
				} else {
					temp = append(temp, cutNumRangeIntoCidrPieces(NumRange{f, l}, 0)...)
					f = a.Low
					l = a.High
				}
			}
			temp = append(temp, cutNumRangeIntoCidrPieces(NumRange{f, l}, 0)...)
			allocRange.AllocRanges = temp
		}

	}
}

func PrintGlobal() {
	//var RangeAllocRangeMap = make(map[uint64]*AllocRange)
	//var RangeTagMap = make(map[uint64]interface{})
	fmt.Println(RangeAllocRangeMap)
	fmt.Println(RangeTagMap)
	fmt.Println(uuidRangeMap)
}
