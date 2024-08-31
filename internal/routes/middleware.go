package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ErrorLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		for _, ginErr := range ctx.Errors {
			logrus.WithError(ginErr).Error("error handling http request")
		}
	}
}