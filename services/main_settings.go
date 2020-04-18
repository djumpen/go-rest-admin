package services

import (
	"fmt"

	"github.com/djumpen/go-rest-admin/models"
	"github.com/djumpen/go-rest-admin/util"
	"github.com/jinzhu/gorm"

	"github.com/Jeffail/gabs"
)

var denyUpdateMainSettingsKeys = []string{
	"id",
	"modified_at",
}

type mainSettingsStorage interface {
	Read(tx *gorm.DB) (*models.MainSettings, error)
	Update(tx *gorm.DB, settings *models.MainSettings) (*models.MainSettings, error)
}

type mainSettings struct {
	st mainSettingsStorage
	db *gorm.DB
}

// New MainSettings storage
func NewMainSettings(st mainSettingsStorage, db *gorm.DB) *mainSettings {
	return &mainSettings{
		st: st,
		db: db,
	}
}

// Get current MainSettings
func (s *mainSettings) Read() (settings *models.MainSettings, err error) {
	err = withTransaction(s.db, func(tx *gorm.DB) error {
		settings, err = s.st.Read(tx)
		if err != nil {
			return fmt.Errorf("main settings service can`t read settings, err: %s", err)
		}
		return nil
	})
	return settings, err
}

// Partion Update JSON
func (s *mainSettings) Update(settings *models.MainSettings) (*models.MainSettings, error) {
	var settingsRes *models.MainSettings
	err := withTransaction(s.db, func(tx *gorm.DB) error {
		currentSettings, err := s.st.Read(tx)
		if err != nil {
			return fmt.Errorf("main settings service can`t update settings: failed to read, err: %s", err)
		}

		jCurrent, err := gabs.ParseJSON(currentSettings.Settings)
		if err != nil {
			return fmt.Errorf("main settings service can`t update settings: failed to parse json, err: %s", err)
		}
		jUpdate, err := gabs.ParseJSON(settings.Settings)
		if err != nil {
			return fmt.Errorf("main settings service can`t update settings: failed to parse json, err: %s", err)
		}
		jUpdateMap, _ := jUpdate.ChildrenMap()
		/* DEPRECATED AS OF 07/12/2019
		jUpdateMap, err := jUpdate.ChildrenMap()
		if err != nil {
			return fmt.Errorf("main settings service can`t update settings: failed to parse json, err: %s", err)
		}
		*/
		for key, container := range jUpdateMap {
			if util.StringInSlice(key, denyUpdateMainSettingsKeys) {
				continue
			}
			jCurrent.Set(container.Data(), key)
		}
		settings.Settings = jCurrent.Bytes()

		settingsRes, err = s.st.Update(tx, settings)
		if err != nil {
			return fmt.Errorf("main settings service can`t update settings, err: %s", err)
		}
		// mainSettings.SettingsFull = jCurrent.Data()
		return err
	})
	return settingsRes, err
}
