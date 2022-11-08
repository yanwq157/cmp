package tools

import (
	"fmt"
	"strconv"
)

func ParseStringToInt64(value string) int64 {

	newValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		fmt.Println("string转换int64出错")
		return 0
	}
	return newValue
}
