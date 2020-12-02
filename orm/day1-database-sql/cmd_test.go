package day1_database_sql

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestEngine(t *testing.T) {
	engine, _ := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/orm?charset=utf8")
	defer engine.Close()
	s := engine.NewSession()

	result, _ := s.Raw("INSERT INTO user(name) values (?),(?)", "dnw", "Alex").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
