package account

import (
	mware "github.com/LuaSavage/bwg-test-task/service-a/internal/middleware"
	"github.com/labstack/echo/v4"
)

func FilterTransactionIdMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Logger().Infof("In %s", transferHandlerUrl)
		return mware.FilterTransactionID(ctx)
	}
}
