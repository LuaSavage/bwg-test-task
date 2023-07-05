package account

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"github.com/LuaSavage/bwg-test-task/service-a/internal/adapter/api/dto"
	mware "github.com/LuaSavage/bwg-test-task/service-a/internal/middleware"
)

const (
	transferHandlerUrl = "/account/transfer"
)

type service interface {
	Transfer(ctx context.Context, requestDTO *dto.TransferRequestDTO) error
}

type Handler struct {
	service   service
	rateLimit rate.Limit
}

func NewHandler(service service, rateLimit int) *Handler {
	return &Handler{
		service:   service,
		rateLimit: rate.Limit(rateLimit),
	}
}

func (h *Handler) Register(e *echo.Echo) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Adding rate limiter to all routes
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(h.rateLimit)))

	e.POST(transferHandlerUrl, h.Transfer, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Logger().Infof("In %s", transferHandlerUrl)
			return mware.FilterTransactionID(ctx)
		}
	})
}

func (h *Handler) Transfer(ctx echo.Context) error {
	// Doing money transfer
	var reqDto dto.TransferRequestDTO
	ctx.Bind(reqDto)
	err := h.service.Transfer(ctx.Request().Context(), &reqDto)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "can't provide transfer transaction")
	}

	return ctx.JSON(http.StatusOK, "All done")
}
