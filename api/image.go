package api

import (
	"bytes"

	"github.com/djumpen/go-rest-admin/util/img"
	"github.com/djumpen/go-rest-admin/util/uploader"
	"github.com/gin-gonic/gin"
)

type ImageReq struct {
	Image string `json:"image" binding:"required"`
}

type ImageResp struct {
	ImageURL string `json:"image_url"`
}

type imageResource struct {
	resp SimpleResponder
}

func NewImageResource(resp SimpleResponder) *imageResource {
	return &imageResource{
		resp: resp,
	}
}

func (r *imageResource) Upload(c *gin.Context) {
	var req ImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	converter := img.GetConverter(req.Image)
	content, err := converter.BytesBufferFromString(req.Image)
	if err != nil {
		c.Error(err)
		return
	}
	url, err := uploader.NewS3Uploader(converter.GetExt()).Upload(bytes.NewReader(content.Bytes()))
	if err != nil {
		c.Error(err)
		return
	}
	r.resp.OK(c, ImageResp{
		ImageURL: url,
	})
}
