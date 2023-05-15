package libpod

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/joyrex2001/kubedock/internal/events"
	"github.com/joyrex2001/kubedock/internal/model/types"
	"github.com/joyrex2001/kubedock/internal/server/httputil"
	"github.com/joyrex2001/kubedock/internal/server/routes/common"
)

// ImagePull - pull one or more images from a container registry.
// https://docs.podman.io/en/latest/_static/api.html?version=v4.2#tag/images/operation/ImagePullLibpod
// POST "/libpod/images/pull"
func ImagePull(cr *common.ContextRouter, c *gin.Context) {
	from := c.Query("reference")
	img := &types.Image{Name: from}
	if cr.Config.Inspector {
		pts, err := cr.Backend.GetImageExposedPorts(from)
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, err)
			return
		}
		img.ExposedPorts = pts
	}

	if err := cr.DB.SaveImage(img); err != nil {
		httputil.Error(c, http.StatusInternalServerError, err)
		return
	}

	cr.Events.Publish(from, events.Image, events.Pull)

	c.JSON(http.StatusOK, gin.H{
		"Id": img.ID,
	})
}
