package session

import (
	"reflect"
	"testing"
	"vgo/orm/day1-database-sql/log"
)

/**
目的：通过反射（reflect）获取结构体绑定的钩子（hooks），并调用
支持CRUD前后调用钩子

比如，设计一个Account类，Account包含有一个隐私字段Password,那么每次查询后都需要
做脱敏处理，才能继续使用。如果提供了AfterQuery的钩子，查询后，自动地将Password字段的值脱敏
就能省去很多冗余代码
*/

// Hooks constants

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

/**
CallMethod calls the registered hooks
1.钩子机制同样是通过反射来实现的，s.RefTable().Model或value即当前会话
正在操作的对象，使用MethodByNames方法反射得到该对象的方法
2.将s *Session作为入参调用。每一个钩子的入参类型均是*Session
 */

func (s *Session) CallMethod(method string, value interface{}) {
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}

type Account struct {
	ID int `vgoorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *Session) error {
	log.Info("before insert", account)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *Session) error {
	log.Info("after query", account)
	account.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	//s := NewSession().
}