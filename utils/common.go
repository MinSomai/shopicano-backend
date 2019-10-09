package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/teris-io/shortid"
	"time"
)

var sid *shortid.Shortid

func init() {
	v, err := shortid.New(10, shortid.DefaultABC, 2342)
	if err != nil {
		panic(err)
	}
	sid = v
}

func NewUUID() string {
	v := uuid.NewV4()
	return v.String()
}

func NewShortUUID() string {
	return sid.MustGenerate()
}

func NewToken() string {
	token := fmt.Sprintf("%s_%d_%s", NewUUID(), time.Now().Unix(), time.Now().UTC())
	return base64.StdEncoding.EncodeToString([]byte(token))
}
