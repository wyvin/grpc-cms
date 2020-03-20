package util

import (
	"bytes"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

/*
用于判断请求是来源于Rpc客户端还是Restful Api的请求，
根据不同的请求注册不同的ServeHTTP服务；
r.ProtoMajor == 2也代表着请求必须基于HTTP/2
*/
func GrpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	//log.Println(otherHandler)
	//if otherHandler != nil {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		grpcServer.ServeHTTP(w, r)
	//	})
	//}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			log.Println(r.Method, r.Host, r.URL.String())
			if r.Method == "POST" {
				body, _ := ioutil.ReadAll(r.Body)
				log.Printf("%s", body)
				// 将字节切片内容写入相应报文
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
			otherHandler.ServeHTTP(w, r)
		}
	})
}
