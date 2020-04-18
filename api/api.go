package api

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/djumpen/go-rest-admin/apperrors"
	"github.com/djumpen/go-rest-admin/util/jwt"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/gin-contrib/cors"
	"github.com/pkg/errors"

	"strings"

	"github.com/djumpen/go-rest-admin/models"
	"github.com/gin-gonic/gin"
)

// TODO escalate limits to config
const (
	LIMIT_DEFAULT = 10
	LIMIT_MAX     = 50
)

var validate *validator.Validate

type SimpleResponder interface {
	OK(c *gin.Context, res interface{})
	OKList(c *gin.Context, res interface{}, p Pagination)
	Created(c *gin.Context, res interface{})
	NotFound(c *gin.Context, err error)
	OKWithCode(c *gin.Context, res interface{}, code int)
	OKListWithCode(c *gin.Context, res interface{}, pagination Pagination, code int)
	BadRequest(c *gin.Context, description string, err error)
}

type Response struct {
	Success bool               `json:"success"`
	Type    string             `json:"type"`
	Data    interface{}        `json:"data"`
	Code    int                `json:"code,omitempty"`
	Debug   *debugResponsePart `json:"debug,omitempty"`
}

type debugResponsePart struct {
	Details string `json:"details"`
	Trace   gin.H  `json:"trace,omimtempty"`
}

type Pagination struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
}

type Responder struct{}

func NewResponder() *Responder {
	return &Responder{}
}

func (r *Responder) OK(c *gin.Context, res interface{}) {
	response(c, http.StatusOK, 0, res, nil)
}

func (r *Responder) OKList(c *gin.Context, res interface{}, pagination Pagination) {
	response(c, http.StatusOK, 0, res, &pagination)
}

func (r *Responder) OKWithCode(c *gin.Context, res interface{}, code int) {
	response(c, http.StatusOK, code, res, nil)
}

func (r *Responder) OKListWithCode(c *gin.Context, res interface{}, pagination Pagination, code int) {
	response(c, http.StatusOK, code, res, &pagination)
}

func (r *Responder) Created(c *gin.Context, res interface{}) {
	response(c, http.StatusCreated, 0, res, nil)
}

func (r *Responder) BadRequest(c *gin.Context, description string, err error) {
	responseErr(c, http.StatusBadRequest, description, err, nil)
}

func (r *Responder) Unauthorized(c *gin.Context, err error) {
	responseErr(c, http.StatusUnauthorized, "", err, nil)
}

func (r *Responder) Forbidden(c *gin.Context) {
	responseErr(c, http.StatusForbidden, "", nil, nil)
}

func (r *Responder) NotFound(c *gin.Context, err error) {
	responseErr(c, http.StatusNotFound, "", err, nil)
}

func (r *Responder) NotAllowed(c *gin.Context, err error) {
	responseErr(c, http.StatusMethodNotAllowed, "", err, nil)
}

func (r *Responder) Conflict(c *gin.Context, err error, msg string) {
	responseErr(c, http.StatusConflict, "", err, nil)
}

func (r *Responder) Unprocessable(c *gin.Context, description string, err error) {
	responseErr(c, http.StatusUnprocessableEntity, description, err, nil)
}

func (r *Responder) InternalError(c *gin.Context, err error) {
	responseErr(c, http.StatusInternalServerError, "", err, nil)
}

func (r *Responder) ResponseErrWithFields(c *gin.Context, fields map[string]string) {
	responseErr(c, http.StatusUnprocessableEntity, "", nil, fields)
}

func response(c *gin.Context, httpCode, code int, res interface{}, pagination *Pagination) {
	respType := "item"
	data := gin.H{
		"item": res,
	}
	if pagination != nil {
		respType = "list"
		data = gin.H{
			"items":      res,
			"pagination": pagination,
		}
	}
	c.JSON(httpCode, Response{
		Success: true,
		Type:    respType,
		Code:    code,
		Data:    data,
	})
	c.Abort()
}

func responseErr(c *gin.Context, httpCode int, description string, err error, fields map[string]string) {
	var errorText string
	if err != nil {
		errorText = err.Error()
	}
	if httpCode == http.StatusInternalServerError {
		errorText = "Server error"
	}
	if description != "" {
		errorText = description
	}
	var debug *debugResponsePart
	if gin.Mode() == gin.DebugMode {
		debug = debugData(err)
	}
	var errData map[string]string
	if err == nil {
		errData = fields
	} else {
		errData = map[string]string{
			"error": errorText,
		}
	}
	c.JSON(httpCode, Response{
		Success: false,
		Type:    "request_error",
		Data: gin.H{
			"errors": errData,
		},
		Debug: debug,
	})
	c.Abort()
}

func debugData(err error) *debugResponsePart {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	if err == nil {
		return nil
	}
	details := err.Error()
	// if pqErr, ok := err.(*pq.Error); ok {
	// 	details = fmt.Sprintf("%s | %+v", pqErr.Error(), pqErr)
	// }
	debugResp := debugResponsePart{
		Details: details,
	}
	if errSt, ok := err.(stackTracer); ok {
		trace := make(gin.H)
		st := errSt.StackTrace()
		for i, v := range st {
			if i > 3 {
				break
			}
			trace[strconv.Itoa(i)] = fmt.Sprintf("%+v", v)
		}
		debugResp.Trace = trace
	}
	return &debugResp
}

// getPKID Get entity ID from resource URL and convert to models.PKID
func getPKID(c *gin.Context) (models.PKID, error) {
	id := c.Param("id")
	pkid, err := strconv.Atoi(id)
	if err != nil {
		return 0, apperrors.NewNotFound(errors.New("Resource not found"))
	}
	return models.PKID(pkid), err
}

//get current user id from JWT in Context
func getJWTPKID(c *gin.Context) (models.PKID, error) {
	token := c.Request.Header.Get("Authorization")
	authValue := strings.Split(token, " ")
	if len(authValue) > 1 {
		return jwt.GetID(authValue[1])
	}
	return 0, errors.New("Invalid auth header")
}

//getJWTString current user id from JWT in Context
func getJWTString(c *gin.Context) (string, error) {
	token := c.Request.Header.Get("Authorization")
	authValue := strings.Split(token, " ")
	if len(authValue) > 1 {
		return jwt.GetString(authValue[1])
	}
	return "", errors.New("Invalid auth header")
}

// Returns Page and Limit from query. Set dafaultLimit to -1 to use global default (10)
func getPageAndLimit(c *gin.Context, defaultLimit int) (p, l int) {
	if defaultLimit <= 0 {
		defaultLimit = LIMIT_DEFAULT
	}
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", strconv.Itoa(defaultLimit))

	ipage, _ := strconv.Atoi(page)
	ilimit, _ := strconv.Atoi(limit)

	if ipage <= 0 {
		ipage = 1
	}
	if ilimit > LIMIT_MAX {
		ilimit = LIMIT_MAX
	}

	return ipage, ilimit
}

// Returns cors configuration
func GetCorsConfig() cors.Config {
	return cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

var once sync.Once

// TODO: check is this code needed or we should use util/validation
// Returns validator
func Validator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		validate.SetTagName("binding")
	})
	return validate
}
