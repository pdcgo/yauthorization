package yauthorization_test

import (
	"errors"
	"log"
	"testing"

	"github.com/pdcgo/yauthorization"
	"github.com/pdcgo/yauthorization/mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(t *testing.T) {
	mock.DbScenario(t, func(tx *gorm.DB) {

		t.Run("test delete from database", func(t *testing.T) {
			domainID := 2
			order := mock.MockOrder{
				Name:     "order 2",
				DomainID: uint(domainID),
			}
			err := tx.Save(&order).Error
			assert.Nil(t, err)

			identity := mock.MockIdentity{}
			err = yauthorization.NewSecQuery(&identity, tx).Delete(&order).Error
			assert.ErrorIs(t, errors.Unwrap(err), yauthorization.ErrPermission)

			t.Run("test dengan permission", func(t *testing.T) {
				orderperm := order.Permission(&identity, yauthorization.Delete)
				orderperm.Policy = yauthorization.Allow

				identity.WithPermission(t, []*yauthorization.EntityPermission{
					orderperm,
				}, tx, func(perms []*yauthorization.EntityPermission) {

					secquery := yauthorization.NewSecQuery(&identity, tx)

					err := secquery.Delete(&order).Error
					assert.Nil(t, err)
				})
			})

		})

		t.Run("testing create database", func(t *testing.T) {
			order := mock.MockOrder{
				Name:     "order 1",
				DomainID: 1,
			}

			identity := mock.MockIdentity{}

			secquery := yauthorization.NewSecQuery(&identity, tx)

			err := secquery.Save(&order).Error
			assert.ErrorIs(t, errors.Unwrap(err), yauthorization.ErrPermission)

			t.Run("test dengan superuser", func(t *testing.T) {
				identity := mock.MockIdentity{
					SuperUser: true,
				}

				secquery := yauthorization.NewSecQuery(&identity, tx)

				err := secquery.Save(&order).Error
				assert.Nil(t, err)
			})

			t.Run("testing create database dengan ada permission", func(t *testing.T) {
				orderperm := order.Permission(&identity, yauthorization.Create)
				orderperm.Policy = yauthorization.Allow

				orderperup := order.Permission(&identity, yauthorization.Update)
				orderperup.Policy = yauthorization.Allow

				identity.WithPermission(t, []*yauthorization.EntityPermission{
					orderperm,
					orderperup,
				}, tx, func(perms []*yauthorization.EntityPermission) {

					secquery := yauthorization.NewSecQuery(&identity, tx)

					err := secquery.Save(&order).Error
					assert.Nil(t, err)
				})

				t.Run("test dengan domain scope salah", func(t *testing.T) {
					secquery := yauthorization.NewSecQuery(&identity, tx)

					order := mock.MockOrder{
						Name:     "order 1",
						DomainID: 2,
					}

					err := secquery.Save(&order).Error
					assert.ErrorIs(t, errors.Unwrap(err), yauthorization.ErrPermission)
				})

			})

			t.Run("test dengan identity deny", func(t *testing.T) {

				orderperm := order.Permission(&identity, yauthorization.Create)
				orderperm.Policy = yauthorization.Deny

				orderperup := order.Permission(&identity, yauthorization.Update)
				orderperup.Policy = yauthorization.Deny

				identity.WithPermission(t, []*yauthorization.EntityPermission{
					orderperm,
					orderperup,
				}, tx, func(perms []*yauthorization.EntityPermission) {

					secquery := yauthorization.NewSecQuery(&identity, tx)

					err := secquery.Save(&order).Error
					assert.ErrorIs(t, errors.Unwrap(err), yauthorization.ErrPermission)
				})

			})

		})

		t.Run("test support by id", func(t *testing.T) {
			identity := mock.MockIdentity{}

			dataup := mock.MockUpBy{}
			secquery := yauthorization.NewSecQuery(&identity, tx)

			secquery.Save(&dataup)

			assert.NotEmpty(t, dataup.UpdatedByID)
			assert.NotEqual(t, 0, dataup.UpdatedByID)
			log.Println("dataasdasd", dataup.UpdatedByID)

		})

	})

}
