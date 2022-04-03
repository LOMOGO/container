package utils

import (
	"fmt"
	"testing"
)

func TestCreateUnionFS(t *testing.T) {
	imagePath := "/home/lomogo/todo/personal/golang/busybox.tar"
	mergePath, err := CreateUnionFS(imagePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mergePath)
}
