package handlers

import (
	"fmt"

	"counter/pkg/counter"

	"github.com/gin-gonic/gin"
)

func BindCounterHandler(
	route gin.IRouter,
	namespace string,
	counter counter.Counter,
) {
	const counterIdField = "counter_id"

	route.Handle("GET", fmt.Sprintf("%s/get", namespace), func(ctx *gin.Context) {
		id := ctx.Query(counterIdField)
		if id == "" {
			ctx.AbortWithStatus(400)
			return
		}

		c, err := counter.Get(ctx, id)
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		ctx.JSON(200, gin.H{
			"count": c,
		})
	})

	route.Handle("POST", fmt.Sprintf("%s/incr", namespace), func(ctx *gin.Context) {
		id := ctx.Query(counterIdField)
		if id == "" {
			ctx.AbortWithStatus(400)
			return
		}

		if err := counter.Incr(ctx, id); err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		ctx.AbortWithStatus(204)
	})
}
