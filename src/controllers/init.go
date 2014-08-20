package controllers

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"libs/log"
	"libs/token"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// 运行变量
var (
	gameToken       *token.Token
	playerMap       *PlayerMap
	request_log_map map[int32]string
)

type PlayerMap struct {
	Lock *sync.RWMutex
	Map  map[int64]*Connect
}

func (p *PlayerMap) Get(uid int64) *Connect {
	p.Lock.RLock()
	var result *Connect
	if val, ok := playerMap.Map[uid]; ok {
		result = val
	}
	p.Lock.RUnlock()
	return result
}

func (p *PlayerMap) Set(uid int64, connect *Connect) {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	if val, ok := p.Map[uid]; ok {
		if connect.Conn != val.Conn {
			val.Conn.Close()
		}
	}
	p.Map[uid] = connect
}

func (p *PlayerMap) Delete(uid int64, connect *Connect) {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	delete(playerMap.Map, uid)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	gameToken = token.NewToken(&token.Adapter{})
	playerMap = new(PlayerMap)
	playerMap.Lock = new(sync.RWMutex)
	playerMap.Map = make(map[int64]*Connect)
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
	playerMap.Lock.RLock()
	defer playerMap.Lock.RUnlock()
	if Connect, ok := playerMap.Map[uid]; !ok {
		return fmt.Errorf("uid : %d not online", uid)
	} else {
		Connect.Chan <- s
	}
	return nil
}

func CountOnline() {
	go func() {
		t := time.Tick(time.Second * 10)
		//t = time.Tick(time.Minute * 5)
		for {
			select {
			case <-t:
				playerMap.Lock.RLock()
				fmt.Println("online num : ", len(playerMap.Map))
				playerMap.Lock.RUnlock()
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
