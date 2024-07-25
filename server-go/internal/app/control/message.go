package control

import (
	"github.com/gorilla/websocket"
	"net/http"
	_ "server-go/internal/app/service"
)

var MessageControl *messageControl

type messageControl struct {
}

var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
