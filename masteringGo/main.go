package main

import (
	"fmt"
)

func main() {
	var array = [...]int{45, 12, 85, 32, 89, 39, 69, 44, 42, 1, 6, 8}
	var tmp int
	for i := 1; i < len(array); i++ {
		for j := i; j > 0; j-- {
			if array[j] < array[j-1] {
				tmp = array[j]
				array[j] = array[j-1]
				array[j-1] = tmp
			}
		}
	}

	for k := 0; k < len(array); k++ {
		fmt.Print(array[k], " ")
	}
}
