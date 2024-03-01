package main

import (
	"errors"
	"fmt"
)

// Hard - функция принимает карту, где имя-ключ,
// возраст-значение. Функция возвращает имена несовершеннолетних
// и ошибку с произвольным текстом, если такие имеются (errors.New())
func Hard(a map[string]int) ([]string, error) {
	b := make([]string, 0, len(a))
	for name, age := range a {
		if age < 18 {
			b = append(b, name)
		}
	}

	var err error // err := nil
	if len(b) != 0 {
		err = errors.New("ошибка в Hard(): в списке 18+ есть несовершеннолетние")
	}
	return b, err
}

func main() {
	aMap := map[string]int{
		"Антон": 3,
		"Иван":  5,
		"Олег":  18,
	}
	fmt.Println("\t список:", aMap)

	if names, err := Hard(aMap); err != nil {
		fmt.Println("\t", err, names)
	} else {
		fmt.Println("\t в списке нет несовершеннолетних ")
	}
}
