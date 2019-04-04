package naiveIPAM

import (
	"strings"
	"strconv"
	"bytes"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"log"
)

func StringIpToInt(ipstring string) uint32 {
	ipSegs := strings.Split(ipstring, ".")
	var ipInt uint32 = 0
	var pos uint32 = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | uint32(tempInt)
		pos -= 8
	}
	return ipInt
}

func IpIntToString(ipInt uint32) string {
	ipSegs := make([]string, 4)
	var len int = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < len; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[len-i-1] = strconv.Itoa(int(tempInt))
		ipInt = ipInt >> 8
	}
	for i := 0; i < len; i++ {
		buffer.WriteString(ipSegs[i])
		if i < len-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

func InitState() {
	b, _ := ioutil.ReadFile("yamls/RangeAllocRangeMap.yaml")
	yaml.Unmarshal(b, &RangeAllocRangeMap)
	b2, _ := ioutil.ReadFile("yamls/RangeTagMap.yaml")
	yaml.Unmarshal(b2, &RangeTagMap)
	b3, _ := ioutil.ReadFile("yamls/uuidRangeMap.yaml")
	yaml.Unmarshal(b3, &uuidRangeMap)

}

func SaveState() {
	res, _ := yaml.Marshal(RangeAllocRangeMap)
	log.Println(res)
	ioutil.WriteFile("yamls/RangeAllocRangeMap.yaml", []byte(res), 0644)
	res, _ = yaml.Marshal(RangeTagMap)
	log.Println(res)
	ioutil.WriteFile("yamls/RangeTagMap.yaml", []byte(res), 0644)
	res, _ = yaml.Marshal(uuidRangeMap)
	log.Println(res)
	ioutil.WriteFile("yamls/uuidRangeMap.yaml", []byte(res), 0644)
}
