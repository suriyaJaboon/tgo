package rd

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

type rb struct {
	*redis.Client
}

type Rb interface {
	Read(ctx context.Context, key string) string
	UX(ctx context.Context, key string) (*ux, error)
}

func NewRB(c *redis.Client) Rb {
	return &rb{c}
}

func (r rb) Read(ctx context.Context, key string) string {
	value, err := r.Get(ctx, key).Result()
	if err != nil {
		return ""
	}

	return value
}

type ux struct {
	U string `json:"u"`
	X int    `json:"x"`
}

func (r rb) UX(ctx context.Context, key string) (*ux, error) {
	bt, err := r.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var u *ux
	if err = json.Unmarshal(bt, &u); err != nil {
		return nil, err
	}

	return u, nil
}
