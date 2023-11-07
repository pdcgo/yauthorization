package main

import (
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(t *testing.T) {
	dbScenario(t, func(tx *gorm.DB) {

		t.Run("testing create database", func(t *testing.T) {
			order := MockOrder{
				Name:     "order 1",
				DomainID: "team1",
			}

			identity := MockIdentity{}

			secquery := NewSecQuery(&identity, tx)

			err := secquery.Save(&order).Error
			assert.ErrorIs(t, err, ErrPermission)

			t.Run("test dengan superuser", func(t *testing.T) {
				identity := MockIdentity{
					SuperUser: true,
				}

				secquery := NewSecQuery(&identity, tx)

				err := secquery.Save(&order).Error
				assert.Nil(t, err)
			})

			t.Run("testing create database dengan ada permission", func(t *testing.T) {
				orderperm := order.Permission(&identity, Create)
				orderperm.Policy = Allow

				orderperup := order.Permission(&identity, Update)
				orderperup.Policy = Allow

				identity.WithPermission(t, []*EntityPermission{
					orderperm,
					orderperup,
				}, tx, func(perms []*EntityPermission) {

					secquery := NewSecQuery(&identity, tx)

					err := secquery.Save(&order).Error
					assert.Nil(t, err)
				})

				t.Run("test dengan domain scope salah", func(t *testing.T) {
					secquery := NewSecQuery(&identity, tx)

					order := MockOrder{
						Name:     "order 1",
						DomainID: "team2",
					}

					err := secquery.Save(&order).Error
					assert.ErrorIs(t, err, ErrPermission)
				})

			})

		})

	})

}

func dbScenario(t *testing.T, handle func(tx *gorm.DB)) {
	id := uuid.New()
	os.MkdirAll("db_test", os.ModeDir)
	fname := "db_test/" + id.String() + ".db"

	db, err := gorm.Open(sqlite.Open(fname), &gorm.Config{})
	assert.Nil(t, err)

	defer os.Remove(fname)
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	db.AutoMigrate(
		&MockOrder{},
		&EntityPermission{},
	)

	handle(db)
}

type MockIdentity struct {
	SuperUser bool
}

func (ident *MockIdentity) WithPermission(
	t *testing.T,
	perms []*EntityPermission,
	db *gorm.DB,
	handler func(perms []*EntityPermission),
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

func (ident *MockIdentity) IdentityID() string {
	return "user1"
}

func (ident *MockIdentity) IsSuperUser() bool {
	return ident.SuperUser
}

type MockOrder struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	DomainID string
}

func (mo *MockOrder) GetDomainID() string {
	return mo.DomainID
}

func (mo *MockOrder) Permission(identity Identity, action Action) *EntityPermission {
	return &EntityPermission{
		IdentityID: identity.IdentityID(),
		DomainID:   mo.GetDomainID(),
		EntityID:   "MockOrder",
		Policy:     Deny,
		Action:     action,
	}
}

// func (mo *MockOrder) EntityID() string {
// 	return "MockOrder"
// }

// func (mo *MockOrder) DomainID() string {
// 	return "Team1"
// }
