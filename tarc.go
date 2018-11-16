package tarc

import (
	"sync/atomic"
	"time"
)

const (
	nent = 16 // MUST be a power of 2
	mask = nent - 1
	cacheLinePad = 64
)

type Ring struct {
	x uint32
	c [nent]entry
	time.Duration
	_ [cacheLinePad]byte
}

func (c *Ring) Put(key, value string) {
	c.c[atomic.AddUint32(&c.x, 1)&mask] = entry{key, value, time.Duration(time.Now().UnixNano())}
}

func (c *Ring) Get(key string) (value string, ok bool) {
	i := atomic.LoadUint32(&c.x) & mask
	si := i
	ei := (si + 1) & mask
	for si != ei {
		v := c.c[si]
		if key == v.key {
			if time.Duration(time.Now().UnixNano())-v.Duration > c.Duration {
				return "", false
			}
			atomic.CompareAndSwapUint32(&c.x, i, si)
			return v.value, true
		}
		si = (si - 1) & mask
	}
	return "", false
}

type entry struct {
	key, value string
	time.Duration
}
