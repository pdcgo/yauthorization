package yauthorization

import (
	"errors"

	"gorm.io/gorm"
)

var ErrPermission = errors.New("permission error")

type PermissionError struct {
	err              error
	NeedPermissions  []*EntityPermission `json:"need_permission"`
	ActualPermission []*EntityPermission `json:"actual_permission"`
}

// Error implements error.
func (permerr *PermissionError) Error() string {
	if len(permerr.ActualPermission) == 0 {
		return "permission kosong"
	}

	return "permission restricted"
}

func (err *PermissionError) Unwrap() error {
	return err.err
}

type Entity interface {
	Permission(identity Identity, action Action) *EntityPermission
	GetDomainID() uint
	GetEntityID() string
}

type EntityUpdateBy interface {
	SetUpdateByID(idnya uint)
}

type Identity interface {
	IsSuperUser() bool
	IdentityID() uint
	GetUserID() uint
	SetRole(tx *gorm.DB, role *RoleIdentity) error
	DeleteRole(tx *gorm.DB, roleID uint) error
	GetRole(tx *gorm.DB, domainID uint) (*RoleIdentity, error)
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
		return &PermissionError{
			NeedPermissions:  q.Permission,
			ActualPermission: permission,
			err:              ErrPermission,
		}
	}

	for _, policy := range permission {
		if policy.Policy == Deny {
			return &PermissionError{
				NeedPermissions:  q.Permission,
				ActualPermission: permission,
				err:              ErrPermission,
			}
		}
	}

	return nil

}

type RawQuery interface {
	Entity
	Raw() string
}

func (q *SecQuery) Preload(query RawQuery, args ...interface{}) *SecQuery {
	perm := query.Permission(q.identity, Read)
	entityID := perm.EntityID
	perm = q.PermHandler(perm)
	perm.EntityID = entityID

	newq := q.copy()

	newq.Permission = append(q.Permission,
		perm,
	)

	newq.Tx = newq.Tx.Preload(query.Raw(), args...)
	return newq
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

func (q *SecQuery) Delete(value Entity) *gorm.DB {
	q.Permission = append(q.Permission,
		q.PermHandler(value.Permission(q.identity, Delete)),
	)

	err := q.CheckPermission()
	if err != nil {
		return &gorm.DB{
			Error: err,
		}
	}
	return q.Tx.Delete(value, value)
}

func (q *SecQuery) Find(value Entity) *gorm.DB {
	q.Permission = append(q.Permission,
		q.PermHandler(value.Permission(q.identity, Read)),
	)

	err := q.CheckPermission()
	if err != nil {
		return &gorm.DB{
			Error: err,
		}
	}
	return q.Tx.Find(value)
}

func (q *SecQuery) Save(value Entity) *gorm.DB {
	q.Permission = append(q.Permission,
		q.PermHandler(value.Permission(q.identity, Create)),
		q.PermHandler(value.Permission(q.identity, Update)),
	)

	v, upby := value.(EntityUpdateBy)
	if upby {
		v.SetUpdateByID(q.identity.GetUserID())
	}

	err := q.CheckPermission()
	if err != nil {
		return &gorm.DB{
			Error: err,
		}
	}

	return q.Tx.Save(value)
}
