package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/chaiyeole/todo/domain"
	"github.com/gin-gonic/gin"
)

type taskHandler struct {
	tasks domain.ITask
}

func NewServer() *gin.Engine {
	return gin.Default()
}

func RunServer(svr *gin.Engine, localhost string) error {
	addr := localhost

	err := svr.Run(addr)
	if err != nil {
		slog.Error("Error while running server", "err", err)

		return err
	}

	return nil
}

func errorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors[0].Err

			switch err.(type) {
			case *domain.CustomError:
				var customErr *domain.CustomError

				if errors.As(ctx.Errors[0], &customErr) {
					ctx.JSON(int(customErr.StatusCode), customErr.Error())

				}
			}
			ctx.Abort()
		}
	}
}

// speak up mf

func InitHandlers(svr *gin.Engine, tasks domain.ITask) {
	h := &taskHandler{
		tasks: tasks,
	}

	svr.Use(errorMiddleware())

	svr.GET("/getAll", h.getAll)

	svr.GET("/get", h.get)

	svr.POST("/set", h.set)
}

func (h *taskHandler) getAll(ctx *gin.Context) {
	getAllRes, err := h.tasks.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)

		return
	}

	ctx.JSON(http.StatusOK, getAllRes.Task)
}

func (h *taskHandler) get(ctx *gin.Context) {
	getReq := new(domain.GetReq)

	// layer wise validation: based on design; try to validate on the earliest level.

	err := ctx.BindJSON(&getReq)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)

		return
	}

	getRes, customErr := h.tasks.Get(ctx, getReq)
	if customErr != nil {
		ctx.AbortWithError(http.StatusBadRequest, customErr)

		return
	}

	ctx.JSON(http.StatusOK, &getRes)
}

func (h *taskHandler) set(ctx *gin.Context) {
	setReq := new(domain.SetReq)

	err := ctx.BindJSON(&setReq.Task)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)

		return
	}

	setRes, customErr := h.tasks.Set(ctx, setReq)
	if customErr != nil {
		ctx.AbortWithError(http.StatusBadRequest, customErr)

		return
	}

	ctx.JSON(http.StatusOK, setRes)
}
