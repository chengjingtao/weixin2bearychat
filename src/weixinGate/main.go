package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"encoding/xml"

	"time"

	"encoding/json"

	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "weixin gate"
	app.Author = "chengjt"
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "Server",
			ShortName: "Server",
			Usage:     "start weixin gate server",
			Action:    server,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host,H",
					Usage: "server bind host to",
					Value: "0.0.0.0",
				},
				cli.IntFlag{
					Name:  "port,p",
					Usage: "bind port",
					Value: 8080,
				},
				cli.StringFlag{
					Name:  "target beary char url,t",
					Usage: "which url that post msg to",
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("app run eror ", err.Error())
	}
}

func server(cli *cli.Context) {
	host := cli.String("host")
	port := cli.Int("port")

	initHandler()
	startServer(host, port)

}
func initHandler() {
	http.HandleFunc("/", handlePost)
}
func startServer(host string, port int) {
	fmt.Println("server is start...")
	err := http.ListenAndServe(fmt.Sprintf("%s:%v", host, port), nil)
	if err != nil {
		fmt.Println("error ", err.Error())
	}
}

func handlePost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	valid(w, req)

	buffer := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))
	_, err := buffer.ReadFrom(req.Body)
	if err != nil {
		fmt.Println("read from body error ", err.Error())
		return
	}

	if buffer.Len() == 0 {
		fmt.Println("receive post body content is empty ")
		return
	}

	msg := buffer.String()

	instance, err := parseMsg(msg)
	if err != nil {
		fmt.Println("parse msg error ", err.Error())
		return
	}

	go func() {
		fmt.Println("post ", instance.MsgId)
		err := postToBearyChat("", instance)
		if err != nil {
			fmt.Println("postToBearyChat error ", err.Error())
			return
		}
		fmt.Println(instance.MsgId + " √")
	}()
}

type weixinMsg struct {
	XMLNAME      xml.Name  `xml:"xml"`
	ToUserName   string    `xml:"ToUserName"`
	FromUserName string    `xml:"FromUserName"`
	CreateTime   time.Time `xml:"CreateTime"`
	MsgType      string    `xml:"MsgType"`
	Content      string    `xml:"Content"`
	PicUrl       string    `xml:"PicUrl"`
	MediaId      string    `xml:"MediaId"`
	MsgId        string    `xml:"MsgId"`
}

func parseMsg(msg string) (*weixinMsg, error) {
	fmt.Sprintln("get msg " + msg)
	instance := weixinMsg{}
	err := xml.Unmarshal([]byte(msg), &instance)
	if err != nil {
		fmt.Println("error in unmarshal msg " + msg + " " + err.Error())
		return nil, err
	}

	return &instance, nil
}

var client = http.Client{}

func postToBearyChat(url string, instance *weixinMsg) error {
	buffer := bytes.NewBufferString("")
	encoder := json.NewEncoder(buffer)
	err := encoder.Encode(&instance)
	if err != nil {
		return errors.New("Encode instance error " + err.Error())
	}

	resp, err := client.Post(url, "application/json", buffer)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("post to beary chat error " + err.Error())
	}
	fmt.Println(resp.StatusCode)
	return nil
}

func valid(rw http.ResponseWriter, request *http.Request) {
	fmt.Println("get request " + request.URL.RawQuery)
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
