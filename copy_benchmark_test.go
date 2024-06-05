package quickhttp

import (
	"bytes"
	"testing"
)

func Benchmark32(b *testing.B) {
	buf := make([]byte, 32)
	src := bytes.Repeat([]byte("1"), 32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(buf, src)
	}
}

func Benchmark64(b *testing.B) {
	buf := make([]byte, 64)
	src := bytes.Repeat([]byte("1"), 64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(buf, src)
	}
}

func Benchmark128(b *testing.B) {
	buf := make([]byte, 128)
	src := bytes.Repeat([]byte("1"), 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(buf, src)
	}
}

func Benchmark256(b *testing.B) {
	buf := make([]byte, 256)
	src := bytes.Repeat([]byte("1"), 256)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(buf, src)
	}
}

func Benchmark512(b *testing.B) {
	buf := make([]byte, 512)
	src := bytes.Repeat([]byte("1"), 512)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(buf, src)
	}
}

func Benchmark1K(b *testing.B) {
	buf := make([]byte, 1024)
	src := bytes.Repeat([]byte("1"), 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(buf, src)
	}
}

func Benchmark2K(b *testing.B) {
	buf := make([]byte, 1024*2)
	src := bytes.Repeat([]byte("1"), 1024*2)
	for i := 0; i < b.N; i++ {
		copy(buf, src)
	}
}
