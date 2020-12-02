package day1_database_sql

import (
	"database/sql"
	"vgo/orm/day1-database-sql/log"
	"vgo/orm/day1-database-sql/session"
)

/**
ORM框架相对于对象和数据库中间的一个桥梁，借助ORM可以避免写繁琐的SQL语言，
仅仅通过操作具体的对象，就能够完成对关系型数据库的操作。

ORM框架是通用的，也就是说可以将任意合法的对象转换成数据库中的表和记录：
如果根据任意类型的指针，得到其对应的结构体的信息。这涉及到了Go语言的反射机制(reflect)，
通过反射，可以获取到对象的结构体名称，成员变量，方法等信息。

设计ORM框架应该关注的问题：
（1）多种数据库的适配问题
（2）如果对象的字段发生改变，数据库表结构能够自动更新，即是否支持数据库自动迁移？
（3）数据库支持的功能很多，例如事务（transaction）,ORM框架能实现哪些？
*/

type Engine struct {
	db *sql.DB
}

/**
Engine的逻辑非常简单，最重要的方法是NewEngine,NewEngine主要做了两件事
1. 连接数据库
2. 调用db.Ping(), 检查数据库是否能够正常连接
*/
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	e = &Engine{db: db}
	log.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// NewSession() 可以通过Engine实例创建会话，进而与数据库进行交互。
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db)
}
