package cacher

import (
	"encoding/json"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/irvankadhafi/talent-hub-service/utils"
	log "github.com/sirupsen/logrus"
)

type MultiResponse struct {
	IDs   []int64 `json:"ids"`
	Count int64   `json:"count"`
}

// NewMultiResponseFromByte converts interface to multi response model.
func NewMultiResponseFromByte(bt []byte) (mr *MultiResponse, err error) {
	if err := json.Unmarshal(bt, &mr); err != nil {
		log.WithField("bt", string(bt)).Error(err)
		return nil, err
	}

	return
}

func ToMultiResponse(ids []int64, count int64) *MultiResponse {
	return &MultiResponse{
		IDs:   ids,
		Count: count,
	}
}

func (mr *MultiResponse) ToByte() []byte {
	if mr == nil {
		return nil
	}

	return utils.ToByte(mr)
}

func get(client redigo.Conn, key string) (value any, err error) {
	defer utils.WrapCloser(client.Close)

	if err := client.Send("MULTI"); err != nil {
		return nil, err
	}

	if err := client.Send("EXISTS", key); err != nil {
		return nil, err
	}

	if err := client.Send("GET", key); err != nil {
		return nil, err
	}

	res, err := redigo.Values(client.Do("EXEC"))
	if err != nil {
		return nil, err
	}

	val, ok := res[0].(int64)
	if !ok || val <= 0 {
		return nil, ErrKeyNotExist
	}

	return res[1], nil
}

func getHashMember(client redigo.Conn, identifier, key string) (value any, err error) {
	defer func() {
		_ = client.Close()
	}()

	if err := client.Send("MULTI"); err != nil {
		return nil, err
	}

	if err := client.Send("HEXISTS", identifier, key); err != nil {
		return nil, err
	}

	if err := client.Send("HGET", identifier, key); err != nil {
		return nil, err
	}

	res, err := redigo.Values(client.Do("EXEC"))
	if err != nil {
		return nil, err
	}

	val, ok := res[0].(int64)
	if !ok || val <= 0 {
		return nil, ErrKeyNotExist
	}

	return res[1], nil
}
