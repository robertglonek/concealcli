package aerospike

import "concealcli/conceal"

type Aerospike struct {
	c           *conceal.Conceal
	m           map[string]string
	mapCallback func(s string, d string) error
}

func Init(key []byte, MapCallback func(s string, d string) error) (*Aerospike, error) {
	a := &Aerospike{
		m:           make(map[string]string),
		mapCallback: MapCallback,
	}
	var err error
	a.c, err = conceal.Init(key)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Aerospike) NodeID(s string) (string, error) {
	d, e := a.c.StringHex(s)
	if e == nil {
		if _, ok := a.m[s]; !ok {
			a.m[s] = d
			e = a.mapCallback(s, d)
		}
	}
	return d, e
}

func (a *Aerospike) IP(s string) (string, error) {
	d, e := a.c.IP(s)
	if e == nil {
		if _, ok := a.m[s]; !ok {
			a.m[s] = d
			e = a.mapCallback(s, d)
		}
	}
	return d, e
}

func (a *Aerospike) Port(s string) (string, error) {
	d, e := a.c.Port(s)
	if e == nil {
		if _, ok := a.m[s]; !ok {
			a.m[s] = d
			e = a.mapCallback(s, d)
		}
	}
	return d, e
}

func (a *Aerospike) Addr(s string) (string, error) {
	d, e := a.c.Addr(s)
	if e == nil {
		if _, ok := a.m[s]; !ok {
			a.m[s] = d
			e = a.mapCallback(s, d)
		}
	}
	return d, e
}

/* TODO add the following definitions here that correctly use StringRanges
Namespace name
Set/bin names
XDR DC names
cluster name
security - username
... and others ...
*/
