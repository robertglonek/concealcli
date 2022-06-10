package conceal

import (
	"errors"
	"fmt"
)

var ErrKeySize = fmt.Errorf("minimum key size is %d", minKeySize)
var ErrInvalidAddress = errors.New("invalid address")
var ErrInvalidPort = errors.New("invalid port")
var ErrInvalidIP = errors.New("invalid IP")
