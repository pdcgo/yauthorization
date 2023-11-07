package yauthorization

import "gorm.io/gorm"

type Action string

const (
	Create Action = "create"
	Update Action = "update"
	Read   Action = "read"
)

type User struct {
	Role       []*Role
	Permission []*Permission
}

type Entity struct {
	EntityName string
}

type Division struct {
	DivisionName string
}

type Role struct {
	Division    *Division
	RoleName    string
	Permissions []*Permission
}

type Permission struct {
	Entity *Entity
	Action Action
}

type Query struct {
}

func (query *Query) CheckAuthorization() *gorm.DB {
	return nil
}
