package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	initHandler()
	startServer()

}
func initHandler() {
	http.HandleFunc("/", valid)
}
func startServer() {
	fmt.Println("server is start...")
	http.ListenAndServe("0.0.0.0:8080", nil)

}

func valid(rw http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Println("ParseForm error : " + err.Error())
		http.Error(rw, "parse form error ", http.StatusInternalServerError)
		return
	}
	if len(request.Form["echostr"]) <= 0 {
		fmt.Println("请求中 echostr 为空")
		http.Error(rw, "echostr is empty ", http.StatusBadRequest)
		return
	}

	echostr := request.Form["echostr"][0]

	rw.WriteHeader(http.StatusOK)

	fmt.Println("echostr is " + echostr)

	rw.Write(bytes.NewBufferString(echostr).Bytes())
	fmt.Println("valid over")
}
