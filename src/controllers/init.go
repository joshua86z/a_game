package controllers

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"libs/log"
	"libs/token"
	_ "models"
//	"protodata"
	"runtime"
	"sync"
	"time"
)

// 运行变量
var (
	gameToken       *token.Token
	playLock        sync.Mutex
	playerMap       map[int64]*Connect
	request_log_map map[int32]string
)

func init() {
	gameToken = token.NewToken(&token.Adapter{})
	playerMap = make(map[int64]*Connect)
	request_log_map = make(map[int32]string)
	CountOnline()
	log.Info("Program Run !")
}

func Handler(ws *websocket.Conn) {
	Connect := &Connect{Conn: ws, Chan: make(chan []byte, 10)}
	Connect.pushToClient()
	Connect.PullFromClient()
	Connect.Close()
}

func SendMessage(uid int64, s []byte) error {
	playLock.Lock()
	if Connect, ok := playerMap[uid]; !ok {
		return fmt.Errorf("uid : %d not online", uid)
	} else {
		Connect.Chan <- s
	}
	playLock.Unlock()
	return nil
}

func CountOnline() {
	go func() {
		t := time.Tick(time.Second * 10)
		//t = time.Tick(time.Minute * 5)
		for {
			select {
			case <-t:
				fmt.Println("online num : ", len(playerMap))
				//models.DB().Exec("INSERT INTO `stat_online`(`online_num`,`online_time`) VALUES (? , NOW())", online)
			}
		}
	}()
}

func lineNum() int {
	_, _, line, ok := runtime.Caller(1)
	if ok {
		return line
	}
	return -1
}
