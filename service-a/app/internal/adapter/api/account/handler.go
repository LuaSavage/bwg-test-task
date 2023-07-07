package account

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	accservice "github.com/LuaSavage/bwg-test-task/service-a/internal/domain/account"
)

const (
	transferHandlerUrl = "/account/transfer"
)

type service interface {
	Transfer(ctx context.Context, request *accservice.TransferRequest) error
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

	e.POST(transferHandlerUrl, h.Transfer, FilterTransactionIdMiddleware)
}

func (h *Handler) Transfer(ctx echo.Context) error {
	// Doing money transfer
	var req accservice.TransferRequest
	ctx.Bind(req)
	err := h.service.Transfer(ctx.Request().Context(), &req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "can't provide transfer transaction")
	}

	return ctx.JSON(http.StatusOK, "All done")
}
