package controller

import (
	"github.com/NorbertKa/LambdaCMS/models"
)

type Handler struct {
	*db.DB
}

func NewHandler() *Handler {
	return &Handler{}
}
