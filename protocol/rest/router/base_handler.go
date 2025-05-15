package router

import (
	"encoding/json"

	"github.com/enesanbar/go-service/core/errors"
	"github.com/enesanbar/go-service/core/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type BaseHandler struct {
	logger log.Factory
}

func NewBaseHandler(logger log.Factory) BaseHandler {
	return BaseHandler{logger: logger}
}

func (bh BaseHandler) DecodeRequest(c echo.Context, requestObject interface{}) error {
	err := json.NewDecoder(c.Request().Body).Decode(&requestObject)
	defer func() {
		err = c.Request().Body.Close()
		if err != nil {
			bh.logger.For(c.Request().Context()).
				With(zap.Error(err)).
				Info("Unable to close request bidy")
		}
	}()
	if err != nil {
		return errors.Error{
			Code:    errors.EINVALID,
			Message: "unable to serialize JSON body.",
			Op:      "DecodeRequest",
			Err:     err,
		}
	}

	return nil
}

func (bh *BaseHandler) NewSuccess(c echo.Context, responseObject interface{}, status int) error {
	return c.JSON(status, NewApiResponse(status, responseObject, nil))
}

func (bh *BaseHandler) NewError(c echo.Context, err error) error {
	routeError := errors.Error{
		Op:  "RouteHandler",
		Err: err,
	}
	bh.logger.For(c.Request().Context()).Error("", zap.Error(err))

	var response ApiResponse
	response = NewApiResponse(ErrorStatus(err), errors.ErrorData(err), routeError)
	return c.JSON(response.Status, response)
}
