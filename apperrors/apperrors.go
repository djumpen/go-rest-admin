package apperrors

import (
	"regexp"
)

type SimpleError struct {
	err error
}

func (e *SimpleError) Error() string {
	return e.err.Error()
}

type DuplicateEntry struct {
	s     string
	entry string
	key   string
}

type Unauthorized struct{ SimpleError }

type NoRows struct{ SimpleError }

type Conflict struct{ SimpleError }

type NotFound struct{ SimpleError }

type Notification struct{ SimpleError }

type BadRequest struct{ SimpleError }

// ----------------------
func NewUnauthorized(err error) *Unauthorized {
	return &Unauthorized{SimpleError{err}}
}

// ----------------------
func NewDuplicateEntry(e error) *DuplicateEntry {
	re := regexp.MustCompile(`Duplicate entry '(\w*)' for key '(\w*)'`)
	match := re.FindStringSubmatch(e.Error())
	var entry string
	var key string
	if len(match) == 3 {
		entry = match[1]
		key = match[2]
	}
	return &DuplicateEntry{
		s:     e.Error(),
		entry: entry,
		key:   key,
	}
}
func (e *DuplicateEntry) Error() string {
	return e.s
}
func (e *DuplicateEntry) Entry() string {
	return e.entry
}
func (e *DuplicateEntry) Key() string {
	return e.key
}

// ----------------------
func NewNoRows(err error) *NoRows {
	return &NoRows{SimpleError{err}}
}

// ----------------------
func NewConflict(err error) *Conflict {
	return &Conflict{SimpleError{err}}
}

// ----------------------
func NewNotFound(err error) *NotFound {
	return &NotFound{SimpleError{err}}
}

// ----------------------
func NewNotification(err error) *Notification {
	return &Notification{SimpleError{err}}
}

func NewBadRequest(err error) *BadRequest {
	return &BadRequest{SimpleError{err}}
}
