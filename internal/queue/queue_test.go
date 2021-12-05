package queue

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func testQueue(t *testing.T, queue PriorityQueue) {
	queue.Push("second", 2)
	queue.Push("first", 0)
	queue.Push("third", 5)
	assert.Equal(t, "first", queue.Pop())
	assert.Equal(t, "second", queue.Pop())
	assert.Equal(t, "third", queue.Pop())
	// No more items
	assert.Equal(t, nil, queue.Pop())
}

func TestLinked(t *testing.T) {
	testQueue(t, NewLinked())
}

func TestPairing(t *testing.T) {
	testQueue(t, NewPairing())
}

func benchmarkQueue(b *testing.B, queue PriorityQueue, n int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			queue.Push(i, rand.Float32())
		}
		for j := 0; j <= n; j++ {
			queue.Pop()
		}
	}
}

func BenchmarkLinked100(b *testing.B) {
	benchmarkQueue(b, NewLinked(), 100)
}

func BenchmarkLinked500(b *testing.B) {
	benchmarkQueue(b, NewLinked(), 500)
}

func BenchmarkLinked1000(b *testing.B) {
	benchmarkQueue(b, NewLinked(), 1000)
}

func BenchmarkLinked10000(b *testing.B) {
	benchmarkQueue(b, NewLinked(), 10000)
}

func BenchmarkPairingQueue100(b *testing.B) {
	benchmarkQueue(b, NewPairing(), 100)
}

func BenchmarkPairingQueue500(b *testing.B) {
	benchmarkQueue(b, NewPairing(), 500)
}

func BenchmarkPairingQueue1000(b *testing.B) {
	benchmarkQueue(b, NewPairing(), 1000)
}

func BenchmarkPairingQueue10000(b *testing.B) {
	benchmarkQueue(b, NewPairing(), 10000)
}
