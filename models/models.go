package models

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	dateLayout = "2006-01-02T15:04:05Z"
	layoutZ    = "2006-01-02T15:04:05Z"
	layoutNoZ  = "2006-01-02T15:04:05"
)

type PKID uint64

type IDEntity struct {
	ID PKID `json:"id" binding:"required"`
}

type listModels []interface{}

type ViewModeler interface {
	ToVM() interface{}
}

// Wrap every element of slise with entity name
// If model implements ViewModeler iface - it will be converted by ToVM()
// It will panic if `entities` is not slice or array
func ToListVM(entities interface{}) listModels {
	convEntities := toInterfaces(entities)
	vms := make(listModels, len(convEntities))
	for i, v := range convEntities {
		if vm, ok := v.(ViewModeler); ok {
			vms[i] = vm.ToVM()
		} else {
			modelName := trimType(v)
			vms[i] = map[string]interface{}{
				modelName: v,
			}
		}
	}
	return vms
}

func (p *PKID) String() string {
	return fmt.Sprintf("%d", p)
}

func trimType(i interface{}) string {
	t := reflect.TypeOf(i).String()
	indx := strings.Index(t, ".")
	return t[indx+1:]
}

// convert slise of entities to []interfase{}
func toInterfaces(s interface{}) []interface{} {
	v := reflect.ValueOf(s)
	// There is no need to check, we want to panic if it's not slice or array
	intf := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		intf[i] = v.Index(i).Interface()
	}
	return intf
}

//Stime converts string to time. Used by services/openings.go and api/openings.go
func StringToTime(s string) (time.Time, error) {
	if s[len(s)-1:] == "Z" {
		return time.Parse(layoutZ, s)
	}
	return time.Parse(layoutNoZ, s)
}

//Tstring converts time to string. Used by services/openings.go and api/openings.go
func TimeToString(t time.Time) string {
	return t.Format(dateLayout)
}
