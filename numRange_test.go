package naiveIPAM

import (
	"testing"
	"fmt"
	"time"
)

func TestModel(t *testing.T) {
	//fmt.Println(NumRangeIntersect(NumRange{Low: 4, High: 10}, NumRange{Low: 4, High: 10}).ToString())
	//fmt.Println(NumRangeIntersect(NumRange{Low: 4, High: 10}, NumRange{Low: 5, High: 20}).ToString())
	nr := NumRange{Low: 1, High: 511}
	fmt.Println(cutNumRangeIntoCidrPieces(nr, 0))
	nr = NumRange{Low: 1, High: 512}
	fmt.Println(cutNumRangeIntoCidrPieces(nr, 0))
	nr = NumRange{Low: 0, High: 511}
	fmt.Println(cutNumRangeIntoCidrPieces(nr, 0))
}

func TestSort(t *testing.T) {
	//nr := []NumRange{
	//	{Low: 1, High: 511},
	//	{Low: 2, High: 255},
	//	NumRange{Low: 1, High: 255},
	//}
	fmt.Println(InitAllocRange(NumRange{Low: 1, High: 255}))
}

func TestAllocRange(t *testing.T) {
	uid := InitAllocRange(NumRange{Low: 1, High: 128})
	ar := uuidRangeMap[uid]
	nr := ar.AllocNIPs(4)
	fmt.Println(ar)
	nr = ar.AllocNIPs(16)
	fmt.Println(ar)
	nr2 := ar.AllocNIPs(8)
	fmt.Println(ar)
	nr = ar.AllocNIPs(8)
	fmt.Println(ar)
	nr = ar.AllocNIPs(4)
	fmt.Println(ar)
	//fmt.Println(RangeAllocRangeMap[nr.GetHighLow()])
	fmt.Println(nr)
	ar.DeAllocNumRange(nr2)
	fmt.Println(ar)
	ar.DeAllocNumRange(NumRange{Low: 16, High: 32})
	fmt.Println(ar)
	ar.AllocSpecifiedCidrNumRange(NumRange{20, 24})
	fmt.Println(ar)
	//PrintGlobal()
	PrintGlobal()
	//fmt.Println(RangeAllocRangeMap[nr.GetHighLow()])
	//fmt.Println(ceiling(4))
}

func TestMultiAllocRange(t *testing.T) {
	uid := InitAllocRange(NumRange{Low: 0, High: 40960})
	ar := uuidRangeMap[uid]
	for i := 0; i < 50; i++ {
		i := i
		go func() {

			fmt.Println(i, ar.AllocNIPs(uint32(i+1)))
			<-time.After(5 * time.Second)
			fmt.Println(uuidRangeMap)

		}()
	}
	<-time.After(500 * time.Second)
}

func TestServer(t *testing.T) {
	Serve()
}
