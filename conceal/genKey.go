package conceal

import (
	"math/rand"
	"time"
)

func GenKey() ([]byte, error) {
	ret := make([]byte, 256)
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	_, err := s.Read(ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
