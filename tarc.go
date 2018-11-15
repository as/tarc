package tarc

import (
	"sync/atomic"
	"time"
)

type cacheLinePad = [64]byte

type Entry struct {
	key, value string
	time.Duration
}

type TARC struct {
	x uint32
	c [8]Entry
	time.Duration
	_ cacheLinePad
}

func (c *TARC) Put(key, value string) {
	c.c[atomic.AddUint32(&c.x, 1)&7] = Entry{key, value, time.Duration(time.Now().UnixNano())}
}

func (c *TARC) Get(key string) (value string, ok bool) {
	mask := uint32(len(c.c) - 1)
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
