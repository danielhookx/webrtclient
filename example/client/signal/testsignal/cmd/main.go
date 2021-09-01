package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/oofpgDLD/webrtclient/example/client/signal/testsignal"
	"io"
	"net"
	"net/http"
)

var upgrader = websocket.Upgrader{
	// allow cross-origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options
var sig *testsignal.Signal

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*") //Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,FZM-APP-ID
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}

var port = flag.String("port", ":19801", "http service port")

func main() {
	flag.Parse()
	sig = testsignal.NewSignal()
	root := gin.Default()
	root.Use(Cors())
	root.POST("/pub", Pub)
	root.GET("/ws", WsHandle)

	fmt.Println("Listen ", *port)
	l, err := net.Listen("tcp", *port)
	err = http.Serve(l, root)
	if err != nil {
		fmt.Println(err)
	}
}

func dispatch(conn *websocket.Conn, ch *testsignal.Channel) {
	var err error
	for {
		message := ch.Ready()
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("write:", err)
			break
		}
	}

	if err != nil && err != io.EOF {
		fmt.Println("dispatch ws error:" + err.Error())
	}
	conn.Close()
}

func WsHandle(cxt *gin.Context) {
	c, err := upgrader.Upgrade(cxt.Writer, cxt.Request, nil)
	if err != nil {
		fmt.Print("upgrade:", err)
		return
	}
	defer c.Close()
	ch := testsignal.NewChannel()
	go dispatch(c, ch)
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}
		fmt.Printf("recv: %s", message)
		switch mt {
		case websocket.BinaryMessage, websocket.TextMessage:
			//解析
			msg, err := testsignal.Parse(message)
			if err != nil {
				fmt.Println(err, "got:", string(message))
				continue
			}
			switch msg.Type {
			case 1:
				msg, err := testsignal.ParseRegisterMsg([]byte(msg.Data))
				if err != nil {
					fmt.Println(err)
					continue
				}
				//注册
				sig.Reg(msg.Name, ch)
				fmt.Println(msg.Name, " registered")
			}
		case websocket.PingMessage, websocket.PongMessage, websocket.CloseMessage:
		default:

		}
	}
}

func Pub(c *gin.Context) {
	var params = struct {
		Name string `json:"name"`
		SDP  string `json:"sdp"`
	}{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.String(http.StatusBadRequest, "", err.Error())
		return
	}
	sig.Push(params.Name, params.SDP)
	c.String(http.StatusOK, "", "ok")
}
