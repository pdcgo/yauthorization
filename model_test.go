package yauthorization_test

import (
	"testing"

	"github.com/pdcgo/yauthorization"
	"github.com/pdcgo/yauthorization/mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func TestModel(t *testing.T) {
	mock.DbScenario(t, func(tx *gorm.DB) {
		ent := yauthorization.EntityInfo{
			Key: "OrderEntity",
			Action: datatypes.NewJSONSlice[yauthorization.Action]([]yauthorization.Action{
				yauthorization.Create,
				yauthorization.Delete,
			}),
		}

		err := tx.Debug().Save(&ent).Error
		assert.Nil(t, err)
	})
}
