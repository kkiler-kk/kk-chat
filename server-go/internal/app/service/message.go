package service

var MessageService = &messageService{}

type messageService struct {
}

// 发送消息
//func (m *messageService) SendMessage(c *gin.Context) error {
//	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
//	if err != nil {
//		log.Error().Err(err)
//		return err
//	}
//	defer func(ws *websocket.Conn) {
//		err = ws.Close()
//		if err != nil {
//			log.Error().Err(err)
//		}
//	}(ws)
//	return m.MsgHandler(c, ws)
//}
//
//func (m *messageService) MsgHandler(c *gin.Context, ws *websocket.Conn) error {
//	redisUtil.Subscribe(c, consts.WebsocketPublish, func(msg *redis.Message, err error) {
//		tm := time.Now().Format(consts.TimeFormat)
//		m := fmt.Sprintf("[ws][%s]:[%s]", tm, msg)
//		err = ws.WriteMessage(1, []byte(m))
//		if err != nil {
//			log.Error().Err(err).Str("err", "发送消息失败").Msg("发送消息失败")
//			return
//		}
//	})
//	return nil
//}
