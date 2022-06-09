package pto

import (
	"context"
	"time"

	"tgo/pto/v1"
)

type TgoClient struct {
	c v1.TgoClient
}

func NewTgoClient(c v1.TgoClient) *TgoClient {
	return &TgoClient{c: c}
}

func (t TgoClient) Tgo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := t.c.Tg(ctx, &v1.Request{
		Uid:  "j",
		Name: "jx",
	})
	if err != nil {
		return err
	}

	return nil
}
