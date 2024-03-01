package main

import "fmt"

// Easy - функция принимает число-возраст,
// определяет совершеннолетний ли человек
// и возвращает результат
func Easy(a int) string {
	if a >= 18 {
		return "совершеннолетний"
	}
	return "несовершеннолетний"
}

func main() {
	var age int
	fmt.Print("Введите возраст (полных лет):")
	if _, err := fmt.Scan(&age); err != nil {
		fmt.Println(err)
	}
	status := Easy(age)
	fmt.Println("человек", status)
}
