package sandbox

import (
	"fmt"
	"strconv"
	"testing"
)

// проверка времени выполнения варианта 1
func BenchmarkItoA1(b *testing.B) {
	var aInt int = 8080
	for i := 0; i < b.N; i++ {
		strconv.Itoa(aInt) // кстати, не возвращает ошибку, в отличие от Atoi()
	}
}

// проверка времени выполнения варианта 2
func BenchmarkItoA2(b *testing.B) {
	var aInt int = 8080
	for i := 0; i < b.N; i++ {
		fmt.Sprint(aInt)
	}
}
