package glogger

import "testing"

func benchmarkTask(b *testing.B, fn func(int)) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(i)
	}
}

func BenchmarkFormat(b *testing.B) {
	formatter := NewDefaultFormatter()
	rec := &Record{}

	benchmarkTask(b, func(i int) {
		formatter.Format(rec)
	})
}
