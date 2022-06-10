package conceal

const minKeySize = 256

type Conceal struct {
	key []byte
}

func Init(key []byte) (*Conceal, error) {
	if len(key) < minKeySize {
		return nil, ErrKeySize
	}
	c := &Conceal{
		key: key,
	}
	return c, nil
}
