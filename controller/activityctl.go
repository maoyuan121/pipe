package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func showActivitiesAction(c *gin.Context) {
	dm, _ := c.Get("dataModel")
	dataModel := *(dm.(*DataModel))
	c.HTML(http.StatusOK, getTheme(c)+"/activities.html", dataModel)
}
