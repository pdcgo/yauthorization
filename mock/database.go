package mock

import (
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/pdcgo/yauthorization"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func DbScenario(t *testing.T, handle func(tx *gorm.DB)) {
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
		&yauthorization.EntityPermission{},
		&MockUpBy{},
		&yauthorization.RoleIdentity{},
	)

	handle(db)
}
