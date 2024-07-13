package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/chaiyeole/todo/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type taskHandler struct {
	tasks domain.ITask
}

func NewServer() *gin.Engine {
	return gin.Default()
}

func RunServer(svr *gin.Engine, host string) error {
	addr := host

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
	// REST API verbs are self-explainatory; verbs need not be present in the path.
	// path should only convey the resource (URL) one or many ; if one, which one,
	svr.GET("/tasks", h.getAllTasks)

	svr.GET("/tasks/:id", h.getTask)

	svr.POST("/tasks", h.setTask)
}

// paths should be self-explainable about what they do and what they will give back

func (h *taskHandler) getAllTasks(ctx *gin.Context) {
	getAllTasksRes, err := h.tasks.GetAllTasks(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)

		return
	}

	ctx.JSON(http.StatusOK, getAllTasksRes.Task)
}

func (h *taskHandler) getTask(ctx *gin.Context) {
	getTaskReq := new(domain.GetTaskReq)

	// layer wise validation: based on design; try to validate on the earliest level.

	id := ctx.Param("id")
	UUID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("Error while parsing string as UUID in API", "err", err)

		return
	}

	getTaskReq.Id = UUID

	getTaskRes, customErr := h.tasks.GetTask(ctx, getTaskReq)
	if customErr != nil {
		ctx.AbortWithError(http.StatusBadRequest, customErr)

		return
	}

	ctx.JSON(http.StatusOK, &getTaskRes)
}

func (h *taskHandler) setTask(ctx *gin.Context) {
	setTaskReq := new(domain.SetTaskReq)

	err := ctx.BindJSON(&setTaskReq.Task)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)

		return
	}

	if setTaskReq.Task.Id.String() == "" {
		err := errors.New("error UUID")
		slog.Error("Error while parsing UUID in API", "err", err)

		return
	}

	setTaskRes, customErr := h.tasks.SetTask(ctx, setTaskReq)
	if customErr != nil {
		ctx.AbortWithError(http.StatusBadRequest, customErr)

		return
	}

	ctx.JSON(http.StatusOK, setTaskRes)
}
