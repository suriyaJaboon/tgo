package pto

import (
	"context"
	"tgo/pto/v1"
)

type TgoServer struct {
	v1.UnimplementedTgoServer
}

func NewTgo() *TgoServer {
	return &TgoServer{}
}

func (t TgoServer) Tg(ctx context.Context, req *v1.Request) (*v1.Response, error) {
	return &v1.Response{
		Status:  true,
		Message: "success",
	}, nil
}
