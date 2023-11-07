package main

import (
	"errors"

	"gorm.io/gorm"
)

var ErrPermission = errors.New("permission error")

type Policy int

const (
	Allow Policy = 1
	Deny  Policy = 0
)

type Action string

const (
	Create Action = "create"
	Update Action = "update"
	Read   Action = "read"
	Delete Action = "delete"
)

type Entity interface {
	Permission(identity Identity, action Action) *EntityPermission
	GetDomainID() string
}

type Identity interface {
	IsSuperUser() bool
	IdentityID() string
}

type EntityPermission struct {
	IdentityID string `gorm:"primaryKey"`
	DomainID   string `gorm:"primaryKey"`
	EntityID   string `gorm:"primaryKey"`
	Action     Action `gorm:"primaryKey"`
	Policy     Policy
}

type SecQuery struct {
	identity    Identity
	SecTx       *gorm.DB
	Tx          *gorm.DB
	Permission  []*EntityPermission
	PermHandler func(perm *EntityPermission) *EntityPermission
}

func NewSecQuery(
	identity Identity,
	tx *gorm.DB,

) *SecQuery {
	return &SecQuery{
		identity:    identity,
		SecTx:       tx,
		Tx:          tx,
		Permission:  []*EntityPermission{},
		PermHandler: func(perm *EntityPermission) *EntityPermission { return perm },
	}
}

func (q *SecQuery) copy() *SecQuery {
	return &SecQuery{
		identity:    q.identity,
		SecTx:       q.SecTx,
		Tx:          q.Tx,
		Permission:  q.Permission,
		PermHandler: q.PermHandler,
	}
}

func (q *SecQuery) CheckPermission() error {
	if q.identity.IsSuperUser() {
		return nil
	}

	permission := []*EntityPermission{}

	query := q.SecTx.Model(&EntityPermission{})
	for _, perm := range q.Permission {
		query = query.Or(perm)
	}

	err := query.Find(&permission).Error
	if err != nil {
		return err
	}

	if len(q.Permission) != len(permission) {
		return ErrPermission
	}

	for _, policy := range permission {
		if policy.Policy == Deny {
			return ErrPermission
		}
	}

	return err

}

func (q *SecQuery) Model(value Entity) *SecQuery {
	newq := q.copy()
	newq.PermHandler = func(perm *EntityPermission) *EntityPermission {
		permb := value.Permission(q.identity, "")
		perm.DomainID = value.GetDomainID()
		perm.EntityID = permb.EntityID

		return perm
	}

	newq.Tx = newq.Tx.Model(value).Where(value)

	return newq
}

func (q *SecQuery) Find(value Entity) *gorm.DB {
	q.Permission = append(q.Permission,
		q.PermHandler(value.Permission(q.identity, Read)),
	)

	err := q.CheckPermission()
	if err != nil {
		return &gorm.DB{
			Error: ErrPermission,
		}
	}
	return q.Tx.Find(value)
}

func (q *SecQuery) Save(value Entity) *gorm.DB {
	q.Permission = append(q.Permission,
		q.PermHandler(value.Permission(q.identity, Create)),
		q.PermHandler(value.Permission(q.identity, Update)),
	)

	err := q.CheckPermission()
	if err != nil {
		return &gorm.DB{
			Error: ErrPermission,
		}
	}

	return q.Tx.Save(value)
}
