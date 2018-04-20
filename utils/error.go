package utils

import "fmt"

func CheckError(err error, message string) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		panic(err)
	}
	fmt.Printf("%s\n", message)
}
