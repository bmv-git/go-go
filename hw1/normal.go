// normal определяет и выводит максимальное,минимальное и среднее арифметическое
// значения элементов (статического) массива.
package main

import "fmt"

func main() {
	a := [...]float64{1, 2, 3, 4, 5}
	N := len(a)
	aMax := a[0]
	aMin := a[0]
	aAvr := a[0]
	for i := 1; i < N; i++ {
		aAvr += a[i]
		if a[i] > aMax {
			aMax = a[i]
		}
		if a[i] < aMin {
			aMin = a[i]
		}
	}
	fmt.Printf("макс. = %v\n", aMax)
	fmt.Printf("мин. = %v\n", aMin)
	fmt.Printf("среднее арифметическое = %v\n", aAvr/float64(N))
}
