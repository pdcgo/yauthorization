package mock

import (
	"testing"

	"github.com/pdcgo/yauthorization"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type MockIdentity struct {
	ID        uint
	SuperUser bool
	Role      *yauthorization.RoleIdentity
}

// DeleteRole implements yauthorization.Identity.
func (*MockIdentity) DeleteRole(tx *gorm.DB, roleID uint) error {
	panic("unimplemented")
}

// GetRole implements yauthorization.Identity.
func (*MockIdentity) GetRole(tx *gorm.DB, domainID uint) (*yauthorization.RoleIdentity, error) {
	panic("unimplemented")
}

// SetRole implements yauthorization.Identity.
func (*MockIdentity) SetRole(tx *gorm.DB, role *yauthorization.RoleIdentity) error {
	panic("unimplemented")
}

func (ident *MockIdentity) WithPermission(
	t *testing.T,
	perms []*yauthorization.EntityPermission,
	db *gorm.DB,
	handler func(perms []*yauthorization.EntityPermission),
) {
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, perm := range perms {
			err := tx.Save(perm).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	assert.Nil(t, err)

	defer func() {
		db.Transaction(func(tx *gorm.DB) error {
			for _, perm := range perms {
				err := tx.Delete(perm).Error
				if err != nil {
					return err
				}
			}
			return nil
		})
	}()

	handler(perms)

}

func (ident *MockIdentity) IdentityID() uint {
	if ident.Role != nil {
		return ident.Role.IdentityID()
	}

	return 1
}

func (ident *MockIdentity) IsSuperUser() bool {
	return ident.SuperUser
}

func (ident *MockIdentity) GetUserID() uint {
	if ident.ID == 0 {
		return 3
	}
	return ident.ID
}

type MockOrder struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	DomainID uint
}

// GetEntityID implements yauthorization.Entity.
func (*MockOrder) GetEntityID() string {
	return "MockOrder"
}

func (mo *MockOrder) GetDomainID() uint {
	return mo.DomainID
}

func (mo *MockOrder) Permission(identity yauthorization.Identity, action yauthorization.Action) *yauthorization.EntityPermission {
	return &yauthorization.EntityPermission{
		IdentityID: identity.IdentityID(),
		DomainID:   mo.GetDomainID(),
		EntityID:   mo.GetEntityID(),
		Policy:     yauthorization.Deny,
		Action:     action,
	}
}

// func (mo *MockOrder) EntityID() string {
// 	return "MockOrder"
// }

// func (mo *MockOrder) DomainID() string {
// 	return "Team1"
// }

type MockUpBy struct {
	ID          uint `gorm:"primarykey"`
	UpdatedByID uint
}

// GetEntityID implements yauthorization.Entity.
func (*MockUpBy) GetEntityID() string {
	return "MockUpBy"
}

func (m *MockUpBy) Permission(identity yauthorization.Identity, action yauthorization.Action) *yauthorization.EntityPermission {
	return &yauthorization.EntityPermission{
		IdentityID: identity.IdentityID(),
		DomainID:   1,
		EntityID:   m.GetEntityID(),
		Policy:     yauthorization.Deny,
		Action:     action,
	}
}

func (mo *MockUpBy) GetDomainID() uint {
	return 1
}

func (mo *MockUpBy) SetUpdateByID(idnya uint) {
	mo.UpdatedByID = idnya
}

type MockDomain struct {
	ID uint
}
