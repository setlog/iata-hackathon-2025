package service

import (
	"com.setlog/internal/configuration"
	"com.setlog/internal/model"
)

type IataService struct {
	config *configuration.Config
	token  *TokenService
}

func NewIataService(config *configuration.Config, token *TokenService) *IataService {
	return &IataService{config: config, token: token}
}

func (service *IataService) CreateIataData(data *model.InspectionTestValidationResponse) error {

	return nil
}
