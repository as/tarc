package tarc

import (
	"fmt"
	"testing"
	"time"
)

func TestTARC(t *testing.T) {
	c := TARC{Duration: time.Second / 2}
	lo := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	hi := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
	i := 0
	for i = range lo {
		c.Put(fmt.Sprint(lo[i]), hi[i])
		ch, ok := c.Get(lo[i])
		if !ok || ch != hi[i] {
			t.Fatalf("have %q, want %q", ch, hi[i])
		}
	}
	if _, ok := c.Get("A"); ok {
		t.Fatal("space: stale entry in cache")
	}
	c.Put("?", "!")
	if ch, ok := c.Get("?"); !ok || ch != "!" {
		t.Fatalf("have %q, want %q", ch, "!")
	}
	time.Sleep(time.Second)
	if _, ok := c.Get("?"); ok {
		t.Fatal("time: stale entry in cache")
	}
}
func BenchmarkTARC(b *testing.B) {
	c := TARC{Duration: time.Second / 2}
	lo := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	hi := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
	i := 0

	b.Run("Put", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if i > len(lo) {
				i = 0
			}
			c.Put(lo[i], hi[i])
		}
	})
	i = 0

	b.Run("Get", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if i > len(lo) {
				i = 0
			}
			c.Get(lo[i])
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for _, cpu := range []int{0, 1, 2, 4} {
			b.Run(fmt.Sprintf("%dWriters", cpu), func(b *testing.B) {
				done := make(chan bool)
				defer close(done)
				for x := 0; x < cpu; x++ {
					go func() {
						for {
							select {
							case <-done:
								return
							default:
								c.Put("x", "y")
								c.Put("y", "x")
							}
						}
					}()
				}
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						c.Get("x")
					}
				})
			})
		}
	})

}
