package utils

import "fmt"

func IsNil(i interface{}) {
	if i != nil {
		err := fmt.Sprintf("Expected nil, got: %v", i);
		panic(err);
	}
}
