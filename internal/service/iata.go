package service

import (
	"com.setlog/internal/configuration"
	"com.setlog/internal/model"
	"com.setlog/internal/model/iata"
	"encoding/json"
	"fmt"
)

type IataService struct {
	config *configuration.Config
	token  *TokenService
}

func NewIataService(config *configuration.Config, token *TokenService) *IataService {
	return &IataService{config: config, token: token}
}
func (service *IataService) CreateIataData(data *model.EntityCollection) error {

	return nil
}

func (service *IataService) createShipmentData(data *iata.Shipment) error {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err, _ = service.token.RequestData("POST", url, payload)
	if err != nil {
		return err
	}
	return nil
}
func (service *IataService) createProductData(data *iata.Product) error {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err, _ = service.token.RequestData("POST", url, payload)
	if err != nil {
		return err
	}
	return nil
}
func (service *IataService) createHwbData(data *iata.Hwb) error {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err, _ = service.token.RequestData("POST", url, payload)
	if err != nil {
		return err
	}
	return nil
}
func (service *IataService) createOrganisationData(data *iata.Organization) error {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err, _ = service.token.RequestData("POST", url, payload)
	if err != nil {
		return err
	}
	return nil
}

func (service *IataService) createPersonData(data *iata.Person) error {

	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err, _ = service.token.RequestData("POST", url, payload)
	if err != nil {
		return err
	}
	return nil
}

func (service *IataService) createPieceData(data *iata.Piece) error {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err, _ = service.token.RequestData("POST", url, payload)
	if err != nil {
		return err
	}
	return nil
}

func (service *IataService) createItemData(data *iata.Item) error {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err, _ = service.token.RequestData("POST", url, payload)
	if err != nil {
		return err
	}
	return nil
}

func (service *IataService) getShipmentData(data *model.InspectionTestValidationResponse) (error, *iata.Shipment) {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	err, body := service.token.RequestData("GET", url, nil)
	if err != nil {
		return err, nil
	}
	var ship *iata.Shipment
	err = json.Unmarshal(body, &ship)
	if err != nil {
		return err, nil
	}

	return nil, ship
}
func (service *IataService) getProductData(data *model.InspectionTestValidationResponse) (error, *iata.Product) {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	err, body := service.token.RequestData("GET", url, nil)
	if err != nil {
		return err, nil
	}
	var product *iata.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		return err, nil
	}

	return nil, product
}
func (service *IataService) getHwbData(data *model.InspectionTestValidationResponse) (error, *iata.Hwb) {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	err, body := service.token.RequestData("GET", url, nil)
	if err != nil {
		return err, nil
	}
	var hwb *iata.Hwb
	err = json.Unmarshal(body, &hwb)
	if err != nil {
		return err, nil
	}

	return nil, hwb
}
func (service *IataService) getOrganisationData(data *model.InspectionTestValidationResponse) (error, *iata.Organization) {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	err, body := service.token.RequestData("GET", url, nil)
	if err != nil {
		return err, nil
	}
	var org *iata.Organization
	err = json.Unmarshal(body, &org)
	if err != nil {
		return err, nil
	}

	return nil, org
}

func (service *IataService) getPersonData(data *model.InspectionTestValidationResponse) (error, *iata.Person) {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	err, body := service.token.RequestData("GET", url, nil)
	if err != nil {
		return err, nil
	}
	var person *iata.Person
	err = json.Unmarshal(body, &person)
	if err != nil {
		return err, nil
	}

	return nil, person
}

func (service *IataService) getPieceData(data *model.InspectionTestValidationResponse) (error, *iata.Piece) {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	err, body := service.token.RequestData("GET", url, nil)
	if err != nil {
		return err, nil
	}
	var piece *iata.Piece
	err = json.Unmarshal(body, &piece)
	if err != nil {
		return err, nil
	}

	return nil, piece
}

func (service *IataService) getItemData(data *model.InspectionTestValidationResponse) (error, *iata.Item) {
	url := fmt.Sprintf("%s/on-carriages", service.config.IataServiceUrl)
	err, body := service.token.RequestData("GET", url, nil)
	if err != nil {
		return err, nil
	}
	var item *iata.Item
	err = json.Unmarshal(body, &item)
	if err != nil {
		return err, nil
	}

	return nil, item
}
