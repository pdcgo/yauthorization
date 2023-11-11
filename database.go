package yauthorization

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func RegisterEntity(tx *gorm.DB, entities ...Entity) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		for _, ent := range entities {
			err := tx.Save(&EntityInfo{
				Key: ent.GetEntityID(),
				Action: datatypes.NewJSONSlice[Action]([]Action{
					Create,
					Read,
					Update,
					Delete,
				}),
			}).Error

			if err != nil {
				return err
			}
		}

		return nil
	})
}
