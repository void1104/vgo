package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

/**
sql.Open()连接数据库，第一个参数是驱动名称，import语句_"github.com/mattn/go-sqlite3"包导入时会注册sqlite3的驱动
第二个参数是数据库的名称，对于SQLite来说，也就是文件名，不存在会新建。返回一个sql.DB实例的指针

1. Exec()用于执行SQL语句，如果是查询语言，不会返回相关的记录。所以查询语句通常使用Query()和QueryRow()，前者可以返回多条记录，后者只返回一条记录
2. Exec(),Query(),QueryRow()接受1个或多个入参，第一个入参是SQL语句，后面的入参是SQL语句中的占位符？对应的值，占位符一般用来防SQL注入
3. QueryRow()的返回值类型是*sql.Row(),row.Scan()接收1个或多个指针作为参数，可以获取对应列的值，
   在这个示例中，只有Name一列，因此传入字符串指针&name即可获得查询的结果
*/
func main() {
	db, _ := sql.Open("sqlite3", "vgo.db")
	defer func() {
		_ = db.Close()
	}()

	_, _ = db.Exec("DROP TABLE IF EXISTS User;")
	_, _ = db.Exec("CREATE TABLE User(Name text);")
	result, err := db.Exec("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam")

	if err != nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}
	row := db.QueryRow("SELECT Name FROM User LIMIT 1")
	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	}
}
