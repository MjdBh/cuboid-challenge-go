package migrations

import (
	"fmt"

	"cuboid-challenge/app/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "20210601021531",
		Migrate: func(tx *gorm.DB) error {
			fmt.Println("Running migration create_cuboid")
			type Cuboid struct {
				models.Model
				Width  uint
				Height uint
				Depth  uint
				BagID  uint
			}

			return tx.AutoMigrate(&Cuboid{})
		},
		Rollback: func(tx *gorm.DB) error {
			fmt.Println("Rollback migration create_cuboid")
			type Cuboid struct{}

			return tx.Migrator().DropTable(&Cuboid{})
		},
	})
}
