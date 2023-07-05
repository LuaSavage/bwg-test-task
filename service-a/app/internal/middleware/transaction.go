package middleware

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/LuaSavage/bwg-test-task/service-a/internal/adapter/api/dto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var requestIDstore = NewRequestIdStore()

type RequestIdStore struct {
	requestIds map[uuid.UUID]uuid.UUID
	mtx        *sync.Mutex
}

func NewRequestIdStore() *RequestIdStore {
	return &RequestIdStore{
		requestIds: make(map[uuid.UUID]uuid.UUID),
		mtx:        &sync.Mutex{},
	}
}

func (r *RequestIdStore) Find(reqDTO *dto.TransferRequestDTO) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if storedReqId, ok := r.requestIds[reqDTO.AccountID]; ok {
		if reqDTO.TransactionId == storedReqId {
			return nil
		}
	}

	return fmt.Errorf("couldn't find pair: %s, %s", reqDTO.AccountID, reqDTO.TransactionId)
}

func (r *RequestIdStore) Register(reqDTO *dto.TransferRequestDTO) (err error) {
	err = r.Find(reqDTO)
	if err != nil {
		return
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.requestIds[reqDTO.AccountID] = reqDTO.TransactionId
	return
}

func FilterTransactionID(ctx echo.Context) error {
	// Checking request body in general
	var requestBody dto.TransferRequestDTO
	err := ctx.Bind(&requestBody)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid request parameters")
	}

	// Is this reqID unique?
	err = requestIDstore.Register(&requestBody)
	if err != nil {
		return ctx.JSON(http.StatusConflict, err.Error())
	}

	return nil
}
