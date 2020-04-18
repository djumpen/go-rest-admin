package models

import (
	"math"
)

type Queryer interface {
	GetQuery(key string) (string, bool)
}

type Filter struct {
	Pager
	list  map[string]interface{}
	order string
}

type Pager struct {
	page  int
	limit int
}

func NewFilter() *Filter {
	e := make(map[string]interface{})
	return &Filter{
		list: e,
		Pager: Pager{
			page:  1,
			limit: -1,
		},
	}
}

// Set value by key
func (f *Filter) SetItem(key string, val interface{}) {
	f.list[key] = val
}

// Set value using gin Query
func (f *Filter) SetItemFromContext(c Queryer, key string) {
	if s, ok := c.GetQuery(key); ok {
		f.SetItem(key, s)
	}
}

func (f *Filter) SetItemsFromContext(c Queryer, keys ...string) {
	for _, v := range keys {
		f.SetItemFromContext(c, v)
	}
}

// Get item by key
func (f Filter) Get(k string) (interface{}, bool) {
	v, ok := f.list[k]
	return v, ok
}

// Return string if it's not empty
func (f Filter) MustGetString(k string) (string, bool) {
	v, ok := f.list[k]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	if len(s) == 0 {
		return "", false
	}
	return s, true
}

// Set page and limit
func (p *Pager) SetPager(page, limit int) {
	p.page = page
	p.limit = limit
}

// Get limit
func (p Pager) Limit() int {
	return p.limit
}

// Calculate and return offset
func (p Pager) Offset() int {
	page := int(math.Abs(float64(p.page)))
	limit := int(math.Abs(float64(p.limit)))
	return limit * (page - 1)
}
