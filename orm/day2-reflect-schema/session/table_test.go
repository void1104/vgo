package session

import (
	"testing"
)

type User struct {
	Name string `vgoorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	//s := NewSession().Model(&User{})
}
