package models

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/pkg/errors"
)

var DefaultMainSettings = ``

type MainSettings struct {
	ID       PKID
	Settings json.RawMessage
}

type MainSettingsVM struct {
	Settings json.RawMessage
}

func (m MainSettings) ToVM() interface{} {
	return &MainSettingsVM{
		Settings: m.Settings,
	}
}

func (m MainSettings) Key(keyPath string) (interface{}, error) {
	j, err := gabs.ParseJSON(m.Settings)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Can't provide value for path %s", keyPath))
	}
	return j.Path(keyPath).Data(), nil
}
