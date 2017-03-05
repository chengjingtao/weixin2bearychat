package action

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	"crypto/sha1"

	"weixin2bearychat/logger"

	"strings"

	"os"

	"io/ioutil"

	"text/template"

	"time"

	"github.com/codegangsta/cli"
)

var (
	host         = "NA"
	port         = 0
	tmplpath     = "NA"
	bearychatURL = ""
)

func Serve(cli *cli.Context) {
	serve(cli)
}

func serve(cli *cli.Context) {
	host = cli.String("host")
	port = cli.Int("port")
	tmplpath = cli.String("tmplpath")

	bearychatURL = cli.String("target")

	addRoute("/", "GET", validHandler)
	addRoute("/", "POST", msgHandler)

	initHandler()
	startServer(host, port)
}

func initHandler() {
	for key, _ := range routes {
		http.HandleFunc(key, h(key))
	}
}

func h(parten string) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		hs := routes[parten]
		for _, h := range hs {
			if h.method == req.Method {
				h.handler(resp, req)
			}
		}
	}
}

type router struct {
	parten  string
	method  string
	handler func(http.ResponseWriter, *http.Request)
}

var routes = map[string][]router{}

func addRoute(parten string, method string, handler func(http.ResponseWriter, *http.Request)) {
	routes[parten] = append(routes[parten], router{
		parten:  parten,
		method:  method,
		handler: handler,
	})
}

func startServer(host string, port int) {
	log := logger.New()
	log.Info("server is start...")
	err := http.ListenAndServe(fmt.Sprintf("%s:%v", host, port), nil)
	if err != nil {
		log.Error("error ", err.Error())
	}
}

func msgHandler(w http.ResponseWriter, req *http.Request) {
	log := logger.New()

	err := req.ParseForm()

	w.WriteHeader(200)

	buffer := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))
	_, err = buffer.ReadFrom(req.Body)
	if err != nil {
		log.Error("read from body error : ", err.Error())
		return
	}

	if buffer.Len() == 0 {
		log.Info("receive post body content is empty ")
		return
	}

	msg := buffer.String()

	instance, err := parseMsg(msg)
	if err != nil {
		log.Error("parse msg error: ", err.Error())
		return
	}

	//TODO control queue
	go func() {
		err := postToBearyChat(log, bearychatURL, instance)
		if err != nil {
			log.Error("post to bearychat error : ", err.Error())
		}
	}()
}

type weixinMsg struct {
	XMLNAME      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int      `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	PicUrl       string   `xml:"PicUrl"`
	MediaId      string   `xml:"MediaId"`
	MsgId        string   `xml:"MsgId"`
}

func parseMsg(msg string) (*weixinMsg, error) {
	if strings.TrimSpace(msg) == "" {
		return nil, errors.New("msg should not be empty")
	}

	instance := weixinMsg{}
	err := xml.Unmarshal([]byte(msg), &instance)
	if err != nil {
		return nil, errors.New("error in unmarshal msg " + msg + " " + err.Error())
	}

	return &instance, nil
}

var client = http.Client{
	Timeout: time.Second * 10,
}

func postToBearyChat(log *logger.Logger, url string, instance *weixinMsg) error {
	log.Info("posting < ", instance.MsgId, " > to bearychat ")
	buffer := bytes.NewBufferString("")
	encoder := json.NewEncoder(buffer)
	err := encoder.Encode(&instance)
	if err != nil {
		return errors.New("Encode instance error " + err.Error())
	}

	content, err := rebuildMsg(instance)
	if err != nil {
		return errors.New("rebuildMsg error " + err.Error())
	}
	begin := time.Now()
	_, err = client.Post(url, "application/json", bytes.NewBufferString(content))

	if err != nil {
		return errors.New("client post " + instance.MsgId + " error :" + err.Error())
	}

	log.Info("post < ", instance.MsgId, " >  √ ， using ", time.Now().Sub(begin).String())

	return nil
}

var unsupportMsgType = errors.New("msg type unsupport")

func rebuildMsg(msgInstance *weixinMsg) (string, error) {

	bts, err := ioutil.ReadFile(tmplpath + msgInstance.MsgType)
	if os.IsNotExist(err) {
		return "", errors.New(tmplpath + msgInstance.MsgType + " 模板文件不存在")
	}
	if err != nil {
		return "", errors.New("读取 " + tmplpath + msgInstance.MsgType + " 文件出现错误 " + err.Error())
	}
	content := string(bts)
	t, err := template.New(msgInstance.MsgType).Parse(content)
	if err != nil {
		return "", errors.New("get " + msgInstance.MsgType + " template instance error")
	}
	buffer := bytes.NewBufferString("")
	err = t.Execute(buffer, msgInstance)
	if err != nil {
		return "", errors.New("execute template " + msgInstance.MsgType + " error : " + err.Error())
	}

	return buffer.String(), nil
}

func _valid(signature, timestamp, nonce, token string) bool {

	str := nonce + timestamp + token
	sha := sha1.New()
	sha.Write([]byte(str))
	return fmt.Sprintf("%x", sha.Sum(nil)) == signature
}
func simpleGetQueryParamter(request *http.Request, key string) string {
	if len(request.Form[key]) <= 0 {
		return ""
	}
	return request.Form[key][0]
}

func validHandler(rw http.ResponseWriter, request *http.Request) {
	log := logger.New()
	err := request.ParseForm()
	if err != nil {
		fmt.Println("ParseForm error : " + err.Error())
		http.Error(rw, "parse form error ", http.StatusInternalServerError)
		return
	}
	var signature = simpleGetQueryParamter(request, "signature")
	var timestamp = simpleGetQueryParamter(request, "timestamp")
	var nonce = simpleGetQueryParamter(request, "nonce")
	var token = "abcchengjt"

	var ok = _valid(signature, timestamp, nonce, token)
	if !ok {
		log.Info("signature valid fauilure ! signature=", signature, " , timestamp="+timestamp, " , nonce="+nonce)
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}
	log.Info("signature valid success")

	if len(request.Form["echostr"]) <= 0 {
		log.Warn("echostr is empty!")
		http.Error(rw, "echostr is empty ", http.StatusBadRequest)
		return
	}

	echostr := request.Form["echostr"][0]

	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes.NewBufferString(echostr).Bytes())

	log.Debug("get echostr : "+echostr, ", valid response over")
}

var (
	weixinAPPID     = "wxb0a748e41e85cff6"
	weixinAppSecret = "51f6835173bd24f8914b09d63657b096"
)

// func getWeixinUsername(openID string) (string, error) {
// 	var url = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
// 	// client.Get()
// 	return "", nil
// }

// func getAccessToken(appSecret string) (string, error) {
// 	var url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", weixinAPPID, weixinAppSecret)
// 	resp, err := client.Get(url)
// 	if err != nil {
// 		return "", errors.New("client get access_token error " + err.Error())
// 	}
// 	buffer := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))
// 	_, err = buffer.ReadFrom(resp.Body)
// 	if err != nil {
// 		return "", errors.New("read from body error " + err.Error())
// 	}
// 	var tokenResp interface{}
// 	err = json.NewDecoder(resp.Body).Decode(tokenResp)
// 	if err != nil {
// 		return "", errors.New("decode error " + err.Error())
// 	}
// 	if ok, _ := tokenResp["errcode"]; ok {
// 		return "", errors.New("client get access token error " + tokenResp["errmsg"])
// 	}

// 	return tokenResp["access_token"], nil
// }
