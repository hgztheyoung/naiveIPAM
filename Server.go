package naiveIPAM

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
)

func Serve() {
	InitState()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		<-c
		SaveState()
	}()

	r := gin.Default()
	r.POST("/alloc_range/", InitAllocRangeReq)
	r.GET("/alloc_range/:id", GetAllocRangeReq)
	r.GET("/uuid_range_map", GetUuidRangeMap)
	r.POST("/alloc_nips/:id", PostAllocNIPs)
	r.Run("0.0.0.0:443")
}

func InitAllocRangeReq(c *gin.Context) {
	var nr NumRange
	c.BindYAML(&nr)
	log.Println(nr)
	ret := InitAllocRange(nr)
	c.YAML(200, ret)
}

func GetAllocRangeReq(c *gin.Context) {
	uid := c.Params.ByName("id")
	uuidRangeMapMutex.RLock()
	defer uuidRangeMapMutex.RUnlock()
	c.YAML(200, uuidRangeMap[uid])
}

func GetUuidRangeMap(c *gin.Context) {
	c.YAML(200, uuidRangeMap)
}

func PostAllocNIPs(c *gin.Context) {
	args := &struct {
		IpCount uint32
	}{}
	c.BindYAML(args)
	uid := c.Params.ByName("id")
	uuidRangeMapMutex.RLock()
	defer uuidRangeMapMutex.RUnlock()
	ar := uuidRangeMap[uid]
	res := ar.AllocNIPs(args.IpCount)
	c.YAML(200, res)
}
