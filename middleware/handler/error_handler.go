package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/djumpen/go-rest-admin/apperrors"
	"github.com/djumpen/go-rest-admin/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	validator "gopkg.in/go-playground/validator.v9"
)

type Responder interface {
	BadRequest(c *gin.Context, description string, err error)
	Unauthorized(c *gin.Context, err error)
	NotFound(c *gin.Context, err error)
	Conflict(c *gin.Context, err error, msg string)
	ResponseErrWithFields(c *gin.Context, fields map[string]string)
	InternalError(c *gin.Context, err error)
}

func ErrorHandler(r Responder, reqs map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// handle errors from previous middleware
		if len(c.Errors) > 0 {
			processError(c, c.Errors[0].Err, r, reqs)
			return
		}
		c.Next()
		// handle errors from main handler
		if len(c.Errors) > 0 {
			processError(c, c.Errors[0].Err, r, reqs)
			return
		}
	}
}

func processError(c *gin.Context, err error, r Responder, reqs map[string]interface{}) {
	if err == nil {
		return
	}

	log.Printf("ERROR, %v", err)

	switch ve := errors.Cause(err).(type) {
	case validator.ValidationErrors:
		fields := make(map[string]string)
		for _, v := range ve {
			fields[getFieldName(v.StructNamespace(), reqs)] = validationErrorToText(v)
		}
		r.ResponseErrWithFields(c, fields)
		return
	case *apperrors.Validation:
		fields := make(map[string]string)
		fields[getFieldName(ve.Namespace(), reqs)] = ve.Error()
		r.ResponseErrWithFields(c, fields)
		return
	case *json.UnmarshalTypeError:
		err := err.(*json.UnmarshalTypeError)
		validationError := unmarshalTypeErrorToValidation(err)
		r.ResponseErrWithFields(c, validationError)
		return
	case *apperrors.NoRows, *apperrors.NotFound:
		r.NotFound(c, ve)
		return
	case *apperrors.Unauthorized:
		r.Unauthorized(c, err)
		return
	case *apperrors.DuplicateEntry, *apperrors.Conflict:
		r.Conflict(c, err, ve.Error())
		return
	case *apperrors.Notification:
		r.Conflict(c, err, ve.Error())
		return
	case *apperrors.BadRequest:
		r.BadRequest(c, ve.Error(), ve)
		return
	}

	r.InternalError(c, err)
}

func validationErrorToText(e validator.FieldError) string {
	word := split(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", word)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", word, e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", word, e.Param())
	case "numeric":
		return fmt.Sprintf("%s must be numeric", word)
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", word, e.Param())
	case "url":
		return fmt.Sprintf("Invalid url format")
	}
	return fmt.Sprintf("%s is not valid", word)
}

func unmarshalTypeErrorToValidation(err *json.UnmarshalTypeError) map[string]string {
	resError := make(map[string]string)
	kind := err.Type.Kind().String()
	if util.StringInSlice(kind, []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64"}) {
		kind = "int"
	}
	resError[err.Field] = fmt.Sprintf("%s has type '%s', but '%s' required", err.Field, err.Value, kind)
	return resError
}
