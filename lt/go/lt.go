package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	messageList         map[string][]string //频道加消息 每次刷新时调用
	messageListPassword map[string]string	//频道加密码
	socketMap           map[string][]socket //频道加websoket连接
)

type socket struct {
	socket  *websocket.Conn
	message chan string
	channel    string
}

func init() {
	//message = make(chan string, 1000)
	messageListPassword = make(map[string]string, 0)
	socketMap = make(map[string][]socket, 0)
	messageList = make(map[string][]string, 0)
}

func (this *socket) read() {
	defer this.socket.Close()
	//接收
	for {
		_, _, err := this.socket.ReadMessage()
		if err != nil {
			fmt.Println(err)
			delete(socketMap, this.channel)
			return
		}
		//fmt.Println(string(messageByte))
		//message <- string(messageByte)
	}
}

func (this *socket) write() {
	defer this.socket.Close()
	this.message <- "进来了"
	//发送
	for {
		if err := this.socket.WriteMessage(websocket.TextMessage, []byte(<-this.message)); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func upper(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}

	req.ParseForm()
	channel := req.Form.Get("channel")

	socketMod := socket{socket: conn, channel: channel, message: make(chan string, 1000)}

	if len(socketMap[channel])==0{
		socketMap[channel] = make([]socket,0)
	}
	socketMap[channel] = append(socketMap[channel], socketMod)

	go socketMod.read()

	go socketMod.write()

	//监控ws是否健康
	/*for {
		time.Sleep(3*time.Second)
		if reply == ""{
			fmt.Println("关闭ws")
			ws.Close()
			return
		}
		reply = ""
	}*/

}

//channel频道
//消息
func InsertMessage(w http.ResponseWriter, r *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	if err := decoder.Decode(&params); err != nil {
		fmt.Println(err)
	}
	messageList[params["channel"]] = append(messageList[params["channel"]], params["message"]+"   时间:"+time.Now().Format("2006-01-02 15:04:05"))
	for _, v := range socketMap[params["channel"]] {
		v.insert(params["message"] + "   时间:" + time.Now().Format("2006-01-02 15:04:05"))
	}
	w.Write([]byte("成功"))
}

func (this *socket) insert(message string) {
	this.message <- message
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	if err := decoder.Decode(&params); err != nil {
		fmt.Println(err)
	}
	messageList[params["channel"]] = make([]string, 0)
	messageList[params["channel"]] = append(messageList[params["channel"]], "开始啦----")
	w.Write([]byte("成功"))
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	if err := decoder.Decode(&params); err != nil {
		fmt.Println(err)
	}
	message, err := json.Marshal(messageList[params["channel"]])
	if err != nil {
		panic(err)
	}
	w.Write(message)
}

func CreateChannel(w http.ResponseWriter, r *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	if err := decoder.Decode(&params); err != nil {
		fmt.Println(err)
	}
	//如果没有改频道则创建该频道 否则进入该频道
	if len(messageList[params["channel"]]) == 0 {
		messageList[params["channel"]] = append(messageList[params["channel"]], "开始啦----")
		messageListPassword[params["channel"]] = params["password"]
	} else {
		if messageListPassword[params["channel"]] != params["password"] {
			w.Write([]byte("加入频道密码输入错误"))
		}
	}
	w.Write([]byte("成功"))
}

func DeleteChannel(w http.ResponseWriter, r *http.Request) {

	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	if err := decoder.Decode(&params); err != nil {
		fmt.Println(err)
	}
	//密码错误，则删除平带失败
	if messageListPassword[params["channel"]] != params["password"] {
		w.Write([]byte("删除频道失败"))
	}
	delete(messageList,params["channel"])
	delete(messageListPassword,params["channel"])
	delete(socketMap,params["channel"])
	w.Write([]byte("成功"))

}

func main() {

	r := chi.NewRouter()
	r.Use(CORS)
	r.HandleFunc("/ws", upper)
	r.Group(func(r chi.Router) {
		r.Use(checkPassword)
		r.Post("/insert/message", InsertMessage) //新增消息
		r.Post("/delete/message", DeleteMessage) //删除消息
		r.Post("/get/message", GetMessage)       //查询消息
	})
	r.Post("/create/channel", CreateChannel) //创建频道
	r.Post("/delete/channel", DeleteChannel) //删除频道
	if err := http.ListenAndServe(":4000", r); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}



func checkPassword(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		buf, err :=ioutil.ReadAll(request.Body)
		if err != nil{
			writer.Write([]byte(err.Error()))
			return
		}
		var params map[string]string
		err = json.Unmarshal(buf,&params)
		if err != nil{
			writer.Write([]byte(err.Error()))
			return
		}

		//如果没有改频道则创建该频道 否则进入该频道
		if messageListPassword[params["channel"]] == params["password"] && params["password"]!=""{
			request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			next.ServeHTTP(writer, request)
		}else{
			writer.Write([]byte("请重新选择频道"))
			return
		}
	})
}


func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		h := writer.Header()
		if request.Header.Get(HeaderOrigin) == "" {
			next.ServeHTTP(writer, request)
			return
		}
		h.Set(HeaderAccessControlAllowOrigin, request.Header.Get(HeaderOrigin))
		h.Set(HeaderAccessControlAllowCredentials, "true")
		switch request.Method {
		case "OPTIONS":
			s := request.Header.Get(HeaderAccessControlRequestHeaders)
			if s != "" {
				h.Set(HeaderAccessControlAllowHeaders, s)
			}
			s = request.Header.Get(HeaderAccessControlRequestMethod)
			if s != "" {
				h.Set(HeaderAccessControlAllowMethods, s)
			}
			writer.WriteHeader(http.StatusOK)
		default:
			h.Set(HeaderVary, "Origin")
			next.ServeHTTP(writer, request)
		}
	})
}

const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderXRequestedWith      = "X-Requested-With"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)
