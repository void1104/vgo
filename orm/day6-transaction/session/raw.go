package session

import (
	"database/sql"
	"strings"
	"vgo/orm/day4-chain-operation/clause"
	"vgo/orm/day4-chain-operation/dialect"
	"vgo/orm/day4-chain-operation/schema"
)

/**
GoLang中的database/sql中支持原生事务:
	1) db.Begin() 获得*sql.Tx对象
	2) tx.Exec()
	3) tx.Commit()
	4) tx.Rollback()
*/

type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	tx       *sql.Tx
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder
	sqlVars  []interface{}
}

// CommonDB is a minimal function set of db
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// DB returns tx if a tx begins. otherwise return *sql.DB
func (s *Session) DB() CommonDB{
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

