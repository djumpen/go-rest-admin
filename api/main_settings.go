package api

import (
	"encoding/json"

	"github.com/djumpen/go-rest-admin/models"
	"github.com/gin-gonic/gin"
)

type MainSettingsReq struct {
	Settings json.RawMessage `binding:"required"`
}

// ----------------------------------

type mainSettingsService interface {
	Read() (*models.MainSettings, error)
	Update(settings *models.MainSettings) (*models.MainSettings, error)
}

type mainSettingsResource struct {
	svc  mainSettingsService
	resp SimpleResponder
}

func NewMainSettingsResource(svc mainSettingsService, resp SimpleResponder) *mainSettingsResource {
	return &mainSettingsResource{
		svc:  svc,
		resp: resp,
	}
}

func (r *mainSettingsResource) Read(c *gin.Context) {
	settings, err := r.svc.Read()
	if err != nil {
		c.Error(err)
		return
	}
	r.resp.OK(c, settings.ToVM())
}

func (r *mainSettingsResource) Update(c *gin.Context) {
	var req MainSettingsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	j, err := json.Marshal(req.Settings)
	if err != nil {
		c.Error(err)
		return
	}
	settings := &models.MainSettings{
		Settings: json.RawMessage(j),
	}
	settings, err = r.svc.Update(settings)
	if err != nil {
		c.Error(err)
		return
	}
	r.resp.OK(c, settings.ToVM())
}
