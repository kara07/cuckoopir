// package cuckoopir

// import (
// 	"math/rand"
// 	"testing"
// 	"time"
// 	"fmt"
// )

// func randomUint32() uint32 {
// 	return rand.Uint32()
// }

// func BenchmarkMultiplyTwoRandomUint32(b *testing.B) {
// 	rand.Seed(time.Now().UnixNano())

// 	b.ResetTimer()
// 	start := time.Now()
// 	for i := 0; i < b.N; i++ {
// 		_ = randomUint32() * randomUint32()
// 	}
// 	b.StopTimer()
//     elapsed := time.Since(start)
//     b.ReportMetric(float64(b.N), "iterations")
//     fmt.Printf("BenchmarkExample1 took %v\n", elapsed)
// }

// func BenchmarkMultiplyRandomUint32AndSmallNumber(b *testing.B) {
// 	rand.Seed(time.Now().UnixNano())

// 	b.ResetTimer()
// 	start := time.Now()
// 	for i := 0; i < b.N; i++ {
// 		_ = randomUint32() * uint32(rand.Intn(3)-1)
// 	}
// 	b.StopTimer()
//     elapsed := time.Since(start)
//     b.ReportMetric(float64(b.N), "iterations")
//     fmt.Printf("BenchmarkExample2 took %v\n", elapsed)
// }

package cuckoopir

import (
	"math/rand"
	"testing"
	"time"
)

func randomUint32() uint32 {
	return rand.Uint32()
}

func BenchmarkMultiplyTwoRandomUint32(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < b.N; i++ {
		_ = randomUint32() * randomUint32()
	}
}

func BenchmarkMultiplyRandomUint32AndSmallNumber(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < b.N; i++ {
		_ = randomUint32() * uint32(0)
	}
}
