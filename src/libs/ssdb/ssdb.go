package ssdb

import (
	"bytes"
	"fmt"
	"libs/lua"
	"net"
	"strconv"
)

const maxClient int = 20

var (
	ssdb [maxClient]*Client
	i    int
)

func SSDB() *Client {

	i++
	if ssdb[maxClient-1] == nil {
		L, err := lua.NewLua("conf/app.lua")
		if err != nil {
			panic(err)
		}
		host := L.GetString("ssdbhost")
		port := L.GetInt("ssdbport")

		L.Close()

		for i := 0; i < maxClient; i++ {
			if ssdb[i], err = Connect(host, port); err != nil {
				panic(err)
			}
		}

	}

	return ssdb[i%maxClient]
}

type Client struct {
	sock     *net.TCPConn
	recv_buf bytes.Buffer
}

func Connect(ip string, port int) (*Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}
	sock, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	var c Client
	c.sock = sock
	return &c, nil
}

func (c *Client) Do(args ...interface{}) ([]string, error) {
	err := c.send(args)
	if err != nil {
		return nil, err
	}
	resp, err := c.recv()
	return resp, err
}

func (c *Client) Set(key string, val string) error {
	resp, err := c.Do("set", key, val)
	if err != nil {
		return err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return fmt.Errorf("bad response")
}

//
func (c *Client) Get(key string) (string, error) {
	resp, err := c.Do("get", key)
	if err != nil {
		return "", err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1], nil
	}
	if resp[0] == "not_found" {
		return "", nil
	}
	return "", fmt.Errorf("bad response")
}

func (c *Client) Del(key string) error {
	resp, err := c.Do("del", key)
	if err != nil {
		return err
	}
	if len(resp) == 1 && resp[0] == "ok" {
		return nil
	}
	return fmt.Errorf("bad response")
}

func (c *Client) Qslice(key string, start int, end int) ([]string, error) {
	resp, err := c.Do("qslice", key, start, end)
	if err != nil || len(resp) <= 1 {
		return make([]string, 0), err
	}
	if resp[0] == "ok" {
		return resp[1:], nil
	}
	return make([]string, 0), fmt.Errorf("bad response")
}

func (c *Client) QpushFront(key string, v ...interface{}) (int, error) {
	args := []interface{}{"qpush_front", key}
	args = append(args, v...)
	resp, err := c.Do(args...)
	if err != nil {
		return 0, err
	}
	if resp[0] == "ok" {
		return strconv.Atoi(resp[1])
	}
	return 0, fmt.Errorf("bad response")
}

func (c *Client) Qpush(key string, v ...interface{}) (int, error) {
	args := []interface{}{"qpush", key}
	args = append(args, v...)
	resp, err := c.Do(args...)
	if err != nil {
		return 0, err
	}
	if resp[0] == "ok" {
		return strconv.Atoi(resp[1])
	}
	return 0, fmt.Errorf("bad response")
}

func (c *Client) QpopFront(key string) (string, error) {
	resp, err := c.Do("qpop_front", key)
	if err != nil {
		return "", err
	}
	if resp[0] == "ok" {
		return resp[1], nil
	}
	return "", fmt.Errorf("bad response")
}

func (c *Client) QpopBack(key string) (string, error) {
	resp, err := c.Do("qpop_back", key)
	if err != nil {
		return "", err
	}
	if resp[0] == "ok" {
		return resp[1], nil
	}
	return "", fmt.Errorf("bad response")
}

func (c *Client) Qsize(key string) (int, error) {
	resp, err := c.Do("qsize", key)
	if err != nil {
		return 0, err
	}
	if resp[0] == "ok" {
		return strconv.Atoi(resp[1])
	}
	return 0, fmt.Errorf("bad response")
}

func (c *Client) send(args []interface{}) error {
	var buf bytes.Buffer
	for _, arg := range args {
		var s string
		switch arg := arg.(type) {
		case string:
			s = arg
		case []byte:
			s = string(arg)
		case []string:
			for _, s := range arg {
				buf.WriteString(fmt.Sprintf("%d", len(s)))
				buf.WriteByte('\n')
				buf.WriteString(s)
				buf.WriteByte('\n')
			}
			continue
		case int:
			s = fmt.Sprintf("%d", arg)
		case int64:
			s = fmt.Sprintf("%d", arg)
		case float64:
			s = fmt.Sprintf("%f", arg)
		case bool:
			if arg {
				s = "1"
			} else {
				s = "0"
			}
		case nil:
			s = ""
		default:
			return fmt.Errorf("bad arguments")
		}
		buf.WriteString(fmt.Sprintf("%d", len(s)))
		buf.WriteByte('\n')
		buf.WriteString(s)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	_, err := c.sock.Write(buf.Bytes())
	return err
}

func (c *Client) recv() ([]string, error) {
	var tmp [8192]byte
	for {
		n, err := c.sock.Read(tmp[0:])
		if err != nil {
			return nil, err
		}
		c.recv_buf.Write(tmp[0:n])
		resp := c.parse()
		if resp == nil || len(resp) > 0 {
			return resp, nil
		}
	}
}

func (c *Client) parse() []string {
	resp := []string{}
	buf := c.recv_buf.Bytes()
	var idx, offset int
	idx = 0
	offset = 0

	for {
		idx = bytes.IndexByte(buf[offset:], '\n')
		if idx == -1 {
			break
		}
		p := buf[offset : offset+idx]
		offset += idx + 1
		//fmt.Printf("> [%s]\n", p);
		if len(p) == 0 || (len(p) == 1 && p[0] == '\r') {
			if len(resp) == 0 {
				continue
			} else {
				c.recv_buf.Next(offset)
				return resp
			}
		}

		size, err := strconv.Atoi(string(p))
		if err != nil || size < 0 {
			return nil
		}
		if offset+size >= c.recv_buf.Len() {
			break
		}

		v := buf[offset : offset+size]
		resp = append(resp, string(v))
		offset += size + 1
	}

	return []string{}
}

// Close The Client Connection
func (c *Client) Close() error {
	return c.sock.Close()
}
