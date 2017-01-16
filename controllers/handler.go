package controller

import (
	"github.com/NorbertKa/LambdaCMS/config"
	"github.com/NorbertKa/LambdaCMS/models"
)

type Handler struct {
	*db.DB
	Conf *config.Config
}

func NewHandler() *Handler {
	return &Handler{}
}
