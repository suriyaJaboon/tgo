package rd

import (
	"context"
	"errors"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRD(t *testing.T) {
	const tk = "tk"
	const token = "token"

	rbd, mock := redismock.NewClientMock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	nrb := NewRB(rbd)

	t.Run("success", func(t *testing.T) {
		mock.ExpectGet(tk).SetVal(token)

		tkStr := nrb.Read(ctx, tk)
		assert.Equal(t, token, tkStr)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("not found key", func(t *testing.T) {
		mock.ExpectGet("x").SetVal("token")

		tkStr := nrb.Read(ctx, tk)
		assert.Equal(t, "", tkStr)
		assert.NotNil(t, mock.ExpectationsWereMet())
	})
}

func TestRDUX(t *testing.T) {
	const k = "ux"

	rbd, mock := redismock.NewClientMock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	nrb := NewRB(rbd)

	t.Run("success", func(t *testing.T) {
		mock.ExpectGet(k).SetVal(`{"u": "ux", "x": 1}`)

		u, err := nrb.UX(ctx, k)
		assert.Nil(t, err)

		assert.NotNil(t, u)
		assert.Equal(t, "ux", u.U)

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectGet(k).SetErr(errors.New("error"))

		_, err := nrb.UX(ctx, k)
		assert.NotNil(t, err)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("error convert", func(t *testing.T) {
		mock.ExpectGet(k).SetVal(``)

		_, err := nrb.UX(ctx, k)
		assert.NotNil(t, err)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func BenchmarkRb_UX(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	const k = "ux"

	rbd, mock := redismock.NewClientMock()
	nrb := NewRB(rbd)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mock.ExpectGet(k).SetVal(`{"u": "ux", "x": 1}`)

	for i := 0; i < b.N; i++ {
		_, _ = nrb.UX(ctx, k)
	}
}
