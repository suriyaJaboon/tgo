package service

const host = "localhost:5000"

func NewGRPC(addr string) (bool, error) {
	return addr == host, nil
}
