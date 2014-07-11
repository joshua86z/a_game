package token

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

type adapter interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Token struct {
	prefix   string
	adapter  adapter
	isUnique bool
}

func NewToken(a adapter) *Token {
	return &Token{"token_", a, true}
}

func (this *Token) NotUnique() {
	this.isUnique = false
}

func (this *Token) SetPrefix(prefix string) {
	this.prefix = prefix
}

// get uid from token
func (this *Token) GetUid(token string) (int64, error) {

	if str, err := this.adapter.Get(this.prefix + token); err != nil {
		return 0, err
	} else {
		return strconv.ParseInt(str, 10, 0)
	}
}

// create new token
func (this *Token) AddToken(uid int64) (string, error) {

	m := md5.New()
	m.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatInt(uid, 10)))
	token := hex.EncodeToString(m.Sum(nil))

	if this.isUnique {
		this.setUidToken(uid, token)
	}

	return token, this.adapter.Set(this.prefix+token, strconv.Itoa(int(uid)))
}

func (this *Token) setUidToken(uid int64, token string) error {

	key := "uid_token_" + strconv.Itoa(int(uid))

	if oldToken, err := this.adapter.Get(key); err == nil {
		this.adapter.Delete(this.prefix + oldToken)
	}

	return this.adapter.Set(key, token)
}
