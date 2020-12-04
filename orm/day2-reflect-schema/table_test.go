package day2_reflect_schema

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `vgoorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	engine, _ := NewEngine("sqlite3", "vgo.db")

	s := engine.NewSession().Model(&User{})

	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}

}
