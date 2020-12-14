package subscribe

import (
	"fmt"
	"strconv"
	"testing"
)

func TestPutSubNode(t *testing.T) {
	sr := GetSubscribeRoot()
	sr.PutSubscribeNode("a", "c1", "", 1, false)
	sr.PutSubscribeNode("a", "c2", "", 0, false)
	sr.PutSubscribeNode("a", "c1", "", 0, false)
	sr.PutSubscribeNode("a", "c2", "", 1, false)
	fmt.Println(sr.ShowSubscribeRoot())
}

func TestDelSubNode(t *testing.T) {
	sr := GetSubscribeRoot()
	sr.PutSubscribeNode("a", "c1", "", 1, false)
	sr.PutSubscribeNode("a/b/c", "c1", "", 1, false)
	sr.PutSubscribeNode("a", "c2", "", 1, false)
	sr.PutSubscribeNode("a/b", "c2", "", 1, false)
	fmt.Println(sr.ShowSubscribeRoot())

	sr.DelSubscribeNode("a", "c2")
	sr.DelSubscribeNode("a", "c3")
	sr.DelSubscribeNode("a/b", "c2")
	// sr.DelSubscribeNode("a/b/c", "c1")
	fmt.Println(sr.ShowSubscribeRoot())
}

func TestMatchSubNode(t *testing.T) {
	sr := GetSubscribeRoot()
	sr.PutSubscribeNode("#", "c0", "", 1, false)
	sr.PutSubscribeNode("a", "c1", "", 1, false)
	sr.PutSubscribeNode("a/b", "c2", "", 1, false)
	sr.PutSubscribeNode("a/+", "c3", "", 1, false)
	sr.PutSubscribeNode("a/#", "c4", "", 1, false)
	sr.PutSubscribeNode("a/b/c", "c5", "", 1, false)
	fmt.Println(sr.ShowSubscribeRoot())

	fmt.Println(sr.MatchSubscribeNode("a/c"))
}

func TestCleanSubscribeNode(t *testing.T) {
	sr := GetSubscribeRoot()
	sr.PutSubscribeNode("#", "c0", "", 1, false)
	sr.PutSubscribeNode("a", "c1", "", 1, false)
	sr.PutSubscribeNode("a/b", "c2", "", 1, false)
	sr.PutSubscribeNode("a/+", "c3", "", 1, false)
	sr.PutSubscribeNode("a/#", "c4", "", 1, false)
	sr.PutSubscribeNode("a/b/c", "c5", "", 1, false)
	fmt.Println(sr.ShowSubscribeRoot())

	sr.CleanSubscribeNode("c3")
	fmt.Println(sr.ShowSubscribeRoot())
}

func BenchmarkMatchSubNode(b *testing.B) {
	sr := GetSubscribeRoot()
	for i := 0; i < 10000; i++ {
		is := strconv.FormatInt(int64(i), 10)
		sr.PutSubscribeNode("#", "c"+is, "", 1, false)
		sr.PutSubscribeNode("a", "c"+is, "", 1, false)
		sr.PutSubscribeNode("a/b", "c"+is, "", 1, false)
		sr.PutSubscribeNode("a/+", "c"+is, "", 1, false)
		sr.PutSubscribeNode("a/#", "c"+is, "", 1, false)
		sr.PutSubscribeNode("a/b/c", "c"+is, "", 1, false)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sr.MatchSubscribeNode("a/b/c/d")
	}
}
