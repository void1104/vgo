package session

import (
	"database/sql"
	"strings"
	"vgo/orm/day1-database-sql/log"
)

/**
db: 即使用sql.Open()方法连接数据库成功之后返回的指针。
sql,sqlVars：用来拼接SQL语句和SQL语句中占位符的对应值。用户调用Raw()方法即可改变这两个变量的值
*/
type Session struct {
	db      *sql.DB
	sql     strings.Builder
	sqlVars []interface{}
}

func New(db *sql.DB) *Session {
	return &Session{db: db}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

/**
封装Exec(),Query(),QueryRow()三个原生方法
封装有两个目的，一是统一打印日志（包括执行的SQL语句和错误日志）
二是执行完成后，清除(s *Session).sql和(s *Session).sqlVars两个变量。这样Session可以复用，
开启一次会话，可以执行多次SQL。
*/

// Exec raw sql with sqlVars
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRow gets a record from db
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows gets a list of records from db
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
