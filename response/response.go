package response

import "cybercampus_module/models"

type Response struct {
	Status  int         `json:"STATUS"`
	Message string      `json:"MESSAGE"`
	Data    interface{} `json:"DATA"`
}

type LoginResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    models.UserLogin `json:"data"`
	Token   string           `json:"token"`
}

type ValidateErrorResponse struct {
	Error       bool        `json:"ERROR"`
	FailedField string      `json:"FAILED_FIELD"`
	Tag         string      `json:"TAG"`
	Value       interface{} `json:"VALUE"`
}