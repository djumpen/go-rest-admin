package services

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/pkg/errors"
)

// Run storage methods within transaction
func withTransaction(db *gorm.DB, f func(tx *gorm.DB) error) error {
	tx := db.Begin()
	defer tx.Rollback()
	fmt.Println("run transaction func")
	if err := f(tx); err != nil {
		return err
	}
	fmt.Println("COMMIT transaction")
	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	return nil
}
