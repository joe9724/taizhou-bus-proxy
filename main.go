package main

import (
	"flag"
	"log"
	"github.com/valyala/fasthttp"
	"taizhou-bus/models"
    "github.com/json-iterator/go"
	"fmt"
)
var json = jsoniter.ConfigCompatibleWithStandardLibrary
var (
	proxyAddr   string
	proxyClient = &fasthttp.HostClient{
		IsTLS: false,
		Addr:  "api.yourdomain.com",

		// set other options here if required - most notably timeouts.
		// ReadTimeout: 60, // 如果在生产环境启用会出现多次请求现象
	}
)

func ReverseProxyHandler(ctx *fasthttp.RequestCtx) {

	//log.Println(ctx, "Hello, world! Requested path is %q", string(ctx.Path()))
	req := &ctx.Request
	resp := &ctx.Response

	prepareRequest(req)

	defer resp.SetConnectionClose()
	if err := proxyClient.Do(req, resp); err != nil {
		ctx.Logger().Printf("error when proxying the request: %s", err)
	}

	postprocessResponse(resp)
}

func prepareRequest(req *fasthttp.Request) {
	// do not proxy "Connection" header.
	//req.Header.Del("Connection")
	// strip other unneeded headers.

	// alter other request params before sending them to upstream host
	req.Header.SetHost(proxyAddr)

	req.SetRequestURI("http://61.132.47.90:8998/BusService/Require_RouteStatData/?RouteID=11")
}

func postprocessResponse(resp *fasthttp.Response) {
	// do not proxy "Connection" header
	//resp.Header.Del("Connection")
	//resp.SkipBody = false
	//resp.AppendBody([]byte("abc"))
	//fmt.Println(string(resp.Body()))
	//先把数据解析出来
	var data []models.LineStationModel
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	json.Unmarshal(resp.Body(), &data)
	if len(data)>0{
	   temp := data[0]
       var ls models.BtkLineStation
       ls.ID = temp.RouteID
       ls.Name = temp.RouteName
       ls.Dir = temp.RouteType
       if len(temp.SegmentList)==0{
		   resp.SetBody([]byte("emptyDir"))
	   }else if len(temp.SegmentList) == 1{
           ls.SEG[0].Info =temp.SegmentList[0].FirstLastShiftInfo
           ls.SEG[0].Price = temp.SegmentList[0].RoutePrice
		   ls.SEG[0].SegID = temp.SegmentList[0].SegmentID
		   ls.SEG[0].SegName = temp.SegmentList[0].SegmentName
		   //重写stations
		   var tempstation models.BtkStation
		   for i:=0; i<len(temp.SegmentList[0].StationList);i++  {
			   tempstation.Lat = temp.SegmentList[0].StationList[i].StationPosition.Latitude
			   tempstation.Lon = temp.SegmentList[0].StationList[i].StationPosition.Longitude
			   tempstation.StID = temp.SegmentList[0].StationList[i].StationID
			   tempstation.StName = temp.SegmentList[0].StationList[i].StationName
			   tempstation.StNo = temp.SegmentList[0].StationList[i].StationNo
			   ls.SEG[0].Stations = append(ls.SEG[0].Stations,tempstation)
		   }
	   	}else if len(temp.SegmentList) == 2{
	   		var seg0 models.Seg
		   seg0.Info =temp.SegmentList[0].FirstLastShiftInfo
		   seg0.Price = temp.SegmentList[0].RoutePrice
		   seg0.SegID = temp.SegmentList[0].SegmentID
		   seg0.SegName = temp.SegmentList[0].SegmentName
		   //重写stations
		   var tempstation models.BtkStation
		   for i:=0; i<len(temp.SegmentList[0].StationList);i++  {
		   	   fmt.Println(temp.SegmentList[0].StationList[i].StationPosition.Longitude)
			   tempstation.Lat = temp.SegmentList[0].StationList[i].StationPosition.Latitude
			   tempstation.Lon = temp.SegmentList[0].StationList[i].StationPosition.Longitude
			   tempstation.StID = temp.SegmentList[0].StationList[i].StationID
			   tempstation.StName = temp.SegmentList[0].StationList[i].StationName
			   tempstation.StNo = temp.SegmentList[0].StationList[i].StationNo
			   seg0.Stations = append(seg0.Stations,tempstation)
		   }
		   ls.SEG = append(ls.SEG,seg0)
		   //
		   var seg1 models.Seg
		   seg1.Info =temp.SegmentList[1].FirstLastShiftInfo
		   seg1.Price = temp.SegmentList[1].RoutePrice
		   seg1.SegID = temp.SegmentList[1].SegmentID
		   seg1.SegName = temp.SegmentList[1].SegmentName
		   //重写stations
		   var tempstation1 models.BtkStation
		   for i:=1; i<len(temp.SegmentList[1].StationList);i++  {
			   tempstation1.Lat = temp.SegmentList[1].StationList[i].StationPosition.Latitude
			   tempstation1.Lon = temp.SegmentList[1].StationList[i].StationPosition.Longitude
			   tempstation1.StID = temp.SegmentList[1].StationList[i].StationID
			   tempstation1.StName = temp.SegmentList[1].StationList[i].StationName
			   tempstation1.StNo = temp.SegmentList[1].StationList[i].StationNo
			   seg1.Stations = append(seg1.Stations,tempstation1)
		   }
		   ls.SEG = append(ls.SEG,seg1)
	   }
		result,err := json.Marshal(ls)
		if err != nil {
			resp.SetBody([]byte("emptyres"))
		}
		resp.SetBody([]byte(result))

	}else{
		resp.SetBody([]byte("empty"))
	}
	//fmt.Println(data[0].RouteID)
	//resp.SetBody([]byte("abc"))

	// strip other unneeded headers

	// alter other response data if needed
	// resp.Header.Set("Access-Control-Allow-Origin", "*")
	// resp.Header.Set("Access-Control-Request-Method", "OPTIONS,HEAD,POST")
	// resp.Header.Set("Content-Type", "application/json; charset=utf-8")
}

func main() {
	port := flag.String("port", "523", "listen port")
	targetAddr := flag.String("target", "61.132.47.90:8998", "your server domain")
	flag.Parse()

	proxyClient.Addr = *targetAddr

	log.Println("port:", *port)
	log.Println("target:", *targetAddr)

	//setupMiddlewares(ReverseProxyHandler)
	/*requestHandler := func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
	}*/

	// 创建自定义服务器。
	s := &fasthttp.Server{
		Handler: ReverseProxyHandler,
		//MaxConnsPerIP:1,
		//MaxRequestsPerConn:1,
		// Every response will contain 'Server: My super server' header.
		Name: "My super server",
		//DisableKeepalive:true,


		// Other Server settings may be set here.
	}


	if err := s.ListenAndServe(":"+*port); err != nil {
		log.Fatalf("error in fasthttp server: %s", err)
	}
	log.Println("start server...")
}
