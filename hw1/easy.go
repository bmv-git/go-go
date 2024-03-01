// easy проверяет целое число с консоли на четность и выводит результат
package main

import (
	"fmt"
)

func main() {
	var a int64
	fmt.Print("введите целое число:")
	if _, err := fmt.Scan(&a); err != nil {
		fmt.Printf("ошибка в easy1: %v\n", err)
	} else {
		ref := "\tчисло четное"
		if a%2 != 0 {
			ref = "\tчисло нечетное"
		}
		fmt.Println(ref)
	}

}
