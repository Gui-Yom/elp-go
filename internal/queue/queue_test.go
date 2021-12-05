package queue

import (
	"math/rand"
	"testing"
)

func benchmarkQueue(queue PriorityQueue, n int, b *testing.B) {
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
	benchmarkQueue(NewLinked(), 100, b)
}

func BenchmarkLinked500(b *testing.B) {
	benchmarkQueue(NewLinked(), 500, b)
}

func BenchmarkLinked1000(b *testing.B) {
	benchmarkQueue(NewLinked(), 1000, b)
}

func BenchmarkLinked10000(b *testing.B) {
	benchmarkQueue(NewLinked(), 10000, b)
}

func BenchmarkPairingQueue100(b *testing.B) {
	benchmarkQueue(NewPairing(), 100, b)
}

func BenchmarkPairingQueue500(b *testing.B) {
	benchmarkQueue(NewPairing(), 500, b)
}

func BenchmarkPairingQueue1000(b *testing.B) {
	benchmarkQueue(NewPairing(), 1000, b)
}

func BenchmarkPairingQueue10000(b *testing.B) {
	benchmarkQueue(NewPairing(), 10000, b)
}
