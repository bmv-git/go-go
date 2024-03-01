package sandbox

import (
	"fmt"
	"strconv"
)

// ItoA - это функция, проверяющая варианты преобразования
// целого числа в строку...
func ItoA(a int) {
	fmt.Printf("strconv.Itoa(a): %T, %v\n", strconv.Itoa(a), strconv.Itoa(a))
	fmt.Printf("fmt.Sprint(a): %T, %v\n", fmt.Sprint(a), fmt.Sprint(a))
}
