package data

import (
	"encoding/json"
	"ercd-test/internal/conf"
	"ercd-test/internal/logger"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type TronData struct {
	mu      sync.RWMutex
	DatChan chan *TronTicker
}

func NewTronData() *TronData {
	return &TronData{
		DatChan: make(chan *TronTicker),
	}
}

func (s *TronData) GetData() {
	conn, _, err := websocket.DefaultDialer.Dial(conf.Wss().Addr, nil)
	if err != nil {
		panic(err)
		return
	}
	defer conn.Close()

	logger.Logrus.Infof("Connected webscoket addr: %v", conf.Wss().Addr)

	msg := "heartbeat"
	data, err := json.Marshal(msg)
	if err != nil {
		logger.Logrus.Error("err:", err)
		return
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
	conn.SetWriteDeadline(time.Now().Add(1 * time.Minute))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
		if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
			logger.Logrus.Error(err)
			return err
		}
		return nil
	})
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		logger.Logrus.Error(err)
		return
	}

	go func() {
		t := time.NewTicker(time.Second * 30)
		for range t.C {
			// log.Println("ping message::")
			conn.SetWriteDeadline(time.Now().Add(time.Minute))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Logrus.Error(err)
				return
			}
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Logrus.Error(err.Error())
			return
		}

		//fmt.Println("REC ==>", string(msg))
		go s.handleMsg(msg)

		//data := &TronTicker{}
		//if err := json.Unmarshal(msg, &data); err != nil {
		//	logger.Logrus.Error(err)
		//	return
		//}
		////fmt.Println("REC ==>", data)
		//
		//s.DatChan <- data

		//if msgType == websocket.TextMessage {
		//}
	}

}

func (s *TronData) handleMsg(msg []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("REC ==>", string(msg))
	data := &TronTicker{}
	if err := json.Unmarshal(msg, &data); err != nil {
		logger.Logrus.Error(err)
		return
	}

	s.DatChan <- data

}

type TronTicker struct {
	Contract    string  `json:"contract"`
	EventName   string  `json:"event_name"`
	Transaction string  `json:"transaction"`
	Block       int64   `json:"block"`
	FromAddress string  `json:"from_address"`
	ToAddress   string  `json:"to_address"`
	Amount      float64 `json:"amount"`
}
