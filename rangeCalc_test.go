package naiveIPAM

import (
	"testing"
	"fmt"
)

func TestCalc(t *testing.T) {
	//fmt.Println(getEndingZerocount(5424))
	//fmt.Println(getEndingZerocount(0))
	//fmt.Println(1 << uint(getLogBase(33)))
	//fmt.Println(cutRangeIntoCidrPieces(16, 512, 1))
	//fmt.Println(cutRangeIntoCidrPieces(16, 64, 8))
	//fmt.Println(cutRangeIntoCidrPieces(110, 517, 4))

	//fmt.Println(allocNIPs(16, 512, 128, 16))
	//fmt.Println(IpIntToString(math.MaxUint32))
	//fmt.Println(StringIpToInt("255.255.255.255"))
	fmt.Println(allocCIDRs("0.0.0.0", "0.0.1.0", 100, 0))
	fmt.Println(rangeToCidr(5435734, 5435789))
}
