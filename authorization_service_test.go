package yauthorization_test

import (
	"encoding/json"
	"testing"

	"github.com/pdcgo/yauthorization"
	"github.com/pdcgo/yauthorization/mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRoleService(t *testing.T) {
	mock.DbScenario(t, func(tx *gorm.DB) {
		authsrv := yauthorization.NewAuthorizeService(tx)

		userdentity := mock.MockIdentity{}
		domain := mock.MockDomain{ID: 1}

		t.Run("test create role tanpa admin", func(t *testing.T) {
			role := yauthorization.RoleIdentity{
				Key:      "cs",
				DomainID: domain.ID,
			}
			err := authsrv.RoleCreate(&userdentity, &role)
			assert.ErrorIs(t, err, yauthorization.ErrPermission)

			t.Run("test create with admin", func(t *testing.T) {
				adminIdentityRole(t, authsrv, domain.ID, func(arole *yauthorization.RoleIdentity) {
					userdentity.Role = arole

					tx = tx.Debug()
					authsrv := yauthorization.NewAuthorizeService(tx)

					err := authsrv.RoleCreate(&userdentity, &role)
					data, _ := json.MarshalIndent(err, "", "\t")
					assert.Nil(t, err, string(data))

				})
			})

		})

		t.Run("test role list", func(t *testing.T) {
			t.Error("not implemented")
		})

		t.Run("test role delete", func(t *testing.T) {
			t.Error("not implemented")
		})

		t.Run("test update permission", func(t *testing.T) {
			t.Error("not implemented")
		})
	})
}

func adminIdentityRole(t *testing.T, authsrv *yauthorization.AuthorizeService, domainID uint, handler func(role *yauthorization.RoleIdentity)) {

	role := yauthorization.RoleIdentity{
		Key:      "admin",
		DomainID: domainID,
	}

	root := mock.MockIdentity{
		SuperUser: true,
		Role:      &role,
	}
	err := authsrv.RoleCreate(&root, &role)
	assert.Nil(t, err)

	ent := &yauthorization.RoleIdentity{
		DomainID: domainID,
	}

	perms := []*yauthorization.EntityPermission{
		ent.Permission(&root, yauthorization.Create),
		ent.Permission(&root, yauthorization.Read),
		ent.Permission(&root, yauthorization.Delete),
		ent.Permission(&root, yauthorization.Update),
	}

	for _, perm := range perms {
		perm.Policy = yauthorization.Allow
	}

	err = authsrv.RoleUpdatePermission(&root, &role, nil, perms)
	assert.Nil(t, err)

	defer func() {
		t.Run("test delete role", func(t *testing.T) {
			err := authsrv.RoleDelete(&root, &role)
			assert.Nil(t, err)
		})
	}()

	handler(&role)
}