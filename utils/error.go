package utils

import "fmt"

func CheckError(err error, message string) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error(), message)
		panic(err)
	}
}

func NoDeadError(err error, message string) bool{
	if err != nil {
		fmt.Println("Exists error ", err.Error(), message)
		return true
	}
	return false
}