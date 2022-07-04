package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"tgo/pto"
	"tgo/pto/mocks"
	"tgo/pto/v1"
)

func TestGRPC_Client(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mtc := mocks.NewMockTgoClient(ctrl)

	t.Run("Tgo success", func(t *testing.T) {
		mtc.EXPECT().Tg(
			gomock.Any(),
			gomock.Any(),
		).Return(&v1.Response{
			Status:  true,
			Message: "success",
		}, nil)

		err := pto.NewTgoClient(mtc).Tgo()
		assert.Nil(t, err)
	})

	t.Run("Tgo error", func(t *testing.T) {
		mtc.EXPECT().Tg(
			gomock.Any(),
			gomock.Any(),
		).Return(&v1.Response{
			Status:  false,
			Message: "error",
		}, status.Error(codes.Internal, "error"))

		err := pto.NewTgoClient(mtc).Tgo()
		assert.NotNil(t, err)
	})
}

func TestTestGRPC_Server(t *testing.T) {
	ts := pto.NewTgo()
	assert.NotNil(t, ts)

	res, err := ts.Tg(context.Background(), &v1.Request{
		Uid:  "j",
		Name: "jx",
	})
	assert.Nil(t, err)
	assert.True(t, res.Status)
}

func FuzzNewGRPC(f *testing.F) {
	f.Add("localhost:5000")
	f.Fuzz(func(t *testing.T, host string) {
		rpc, err := NewGRPC(host)
		if !rpc && err != nil {
			t.Errorf("New gRPC type bool: %v", rpc)
		}
	})
}
