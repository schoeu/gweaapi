package utils

import (
	"fmt"
)

func ErrHandle(err error){
	if err != nil {
		fmt.Println(err)
	}
}