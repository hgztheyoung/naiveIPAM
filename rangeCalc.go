package naiveIPAM

import (
	"math"
	"fmt"
)

// 32-bit word input to count zero bits on right
func getEndingZerocount(v uint32) uint {
	Mod37BitPosition := // map a bit value mod 37 to its position
		[]uint{
			32, 0, 1, 26, 2, 23, 27, 0, 3, 16, 24, 30, 28, 11, 0, 13, 4,
			7, 17, 0, 25, 22, 31, 15, 29, 10, 12, 6, 0, 21, 14, 9, 5,
			20, 8, 19, 18,
		}
	return Mod37BitPosition[(-v&v)%37]
}

func getLogBase(v uint32) (r uint) {
	MultiplyDeBruijnBitPosition := []uint{
		0, 9, 1, 10, 13, 21, 2, 29, 11, 14, 16, 18, 22, 25, 3, 30,
		8, 12, 20, 28, 15, 17, 24, 7, 19, 27, 23, 6, 26, 5, 4, 31,
	}
	v |= v >> 1 // first round down to one less than a power of 2 
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	r = MultiplyDeBruijnBitPosition[(uint32)(v*0x07C4ACDD)>>27]
	return
}

//cut [f,l) range into cidr ranges
//say [0,256) will return [0,256)
//[4,256) will return [[4 8) [8 16) [16 32) [32 64) [64 128) [128 256)]

func cutNumRangeIntoCidrPieces(numRange NumRange, maxSubnetLen uint32) []NumRange {
	return cutRangeIntoCidrPieces(numRange.Low, numRange.High, maxSubnetLen)
}

func GetFirstLargeEnoughNumRange(f, l uint32, wangtedSubNetLen uint32) NumRange {
	p := f
	for p < l {
		step := uint32(1 << getEndingZerocount(p))
		if step == 0 {
			step = 1 << 31
		}
		for p+step > l {
			step = step >> 1
		}
		if step >= wangtedSubNetLen {
			return NumRange{p, p + step}
		}
		p += step
	}
	return NumRange{}
}


func cutRangeIntoCidrPieces(f, l uint32, maxSubnetLen uint32) []NumRange {
	res := make([]NumRange, 0)
	if maxSubnetLen == 0 { // avoid infinite loop
		maxSubnetLen = math.MaxUint32
	}

	p := f
	for p < l {
		step := uint32(1 << getEndingZerocount(p))
		if step == 0 {
			step = 1 << 31
		}
		for p+step > l || step > maxSubnetLen {
			step = step >> 1
		}
		//fmt.Println(p, p+step-1)
		res = append(res, NumRange{p, p + step})
		p += step
	}
	return res
}

func ceiling(n uint32) uint32 {
	if 1<<uint32(getLogBase(n)) == n {
		return n
	}
	return 2 << uint32(getLogBase(n))
}

func allocNIPs(f, l uint32, n int, maxSubnetLen uint32) []NumRange {
	ranges := cutRangeIntoCidrPieces(f, l, maxSubnetLen)
	nrest := uint32(n)
	ret := make([]NumRange, 0)

	for _, r := range ranges {
		rcount := r.High - r.Low + 1
		want := ceiling(nrest)
		if rcount > want {
			ret = append(ret, NumRange{r.Low, r.Low + want})
			return ret
		}
		ret = append(ret, r)
		nrest -= rcount
	}
	return ret
}

func allocCIDRs(f, l string, n int, maxSubnetLen uint32) []string {
	fi := StringIpToInt(f)
	li := StringIpToInt(l)
	res := allocNIPs(fi, li, n, maxSubnetLen)
	ret := make([]string, 0)
	for _, r := range res {
		ret = append(ret, rangeToCidr(r.Low, r.High))
	}
	return ret
}
func rangeToCidr(u uint32, u2 uint32) string {
	size := u2 - u
	logsize := getLogBase(size)
	u = (math.MaxUint32 << logsize) & u
	return fmt.Sprintf("%v/%v", IpIntToString(u), 32-logsize)
}
