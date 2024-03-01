package main

import "fmt"

// Norm - функция принимает набор чисел-возрастов,
// и определяет есть ли хоть один совершеннолетний в наборе, возвращает «Да»
// или «Нет»
func Norm(a []int) string {
	for _, age := range a {
		if age >= 18 {
			return "Да"
		}
	}
	return "Нет"
}

func main() {
	aSlice := []int{1, 9, 18, 7}
	fmt.Println("\tсписок возрастов:", aSlice)
	switch Norm(aSlice) {
	case "Да":
		fmt.Println("\tв списке есть совершеннолетние")
	case "Нет":
		fmt.Println("\tв списке нет совершеннолетних")
	default:
	}
}
