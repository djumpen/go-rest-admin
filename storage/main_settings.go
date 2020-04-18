package storage

import (
	"encoding/json"

	"github.com/jinzhu/gorm"

	"github.com/djumpen/go-rest-admin/models"
)

type mainSettingsStorage struct{}

func NewMainSettingsStorage() *mainSettingsStorage {
	return &mainSettingsStorage{}
}

func (s *mainSettingsStorage) Update(tx *gorm.DB, settings *models.MainSettings) (*models.MainSettings, error) {
	var settingsRes models.MainSettings
	err := tx.Assign(settings).FirstOrCreate(&settingsRes).Error
	if err != nil {
		return nil, wrapErr(err)
	}
	return &settingsRes, nil
}

func (s *mainSettingsStorage) Read(tx *gorm.DB) (*models.MainSettings, error) {
	var settings models.MainSettings
	err := tx.First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			emptySettings := &models.MainSettings{
				Settings: json.RawMessage([]byte(models.DefaultMainSettings)),
			}
			settingsRes, err := s.Update(tx, emptySettings)
			if err != nil {
				return nil, wrapErr(err)
			}
			settings = *settingsRes
		} else {
			return nil, wrapErr(err)
		}
	}
	return &settings, nil
}
