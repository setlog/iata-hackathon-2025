package service

import (
	"com.setlog/internal/configuration"
	"com.setlog/internal/model"
	"com.setlog/internal/model/iata"
	"encoding/json"
	"log/slog"
)

type IataService struct {
	config *configuration.Config
	token  *TokenService
}

func NewIataService(config *configuration.Config, token *TokenService) *IataService {
	return &IataService{config: config, token: token}
}
func (service *IataService) CreateIataData(data *model.EntityCollection) error {
	err, itemLocations := service.createItemData(data.Items)
	if err != nil {
		return err
	}
	err, pieceLocations := service.createPieceData(data.Pieces, itemLocations)
	if err != nil {
		return err
	}
	err, orgLocations := service.createOrganisationData(data.Organizations)
	if err != nil {
		return err
	}
	err, shipLocations := service.createShipmentData(data.Shipments, pieceLocations, orgLocations)
	if err != nil {
		return err
	}
	err = service.createHwbData(data.Hwbs, shipLocations[0], orgLocations)
	if err != nil {
		return err
	}

	return nil
}

func (service *IataService) createShipmentData(data []iata.Shipment, pieces []string, orga map[string]string) (error, []string) {
	var shipLoc []string
	for _, ship := range data {
		for _, piece := range pieces {
			obj := iata.ObjectLink{Id: piece}
			ship.ShipmentOfPieces = append(ship.ShipmentOfPieces, obj)
		}
		for key, value := range orga {
			party := iata.InvolvedParty{}
			party.Type = "cargo:Party"
			party.Role = key
			party.Organization = iata.ObjectLink{Id: value}
			ship.InvolvedParties = append(ship.InvolvedParties, party)
		}

		payload, err := json.Marshal(ship)
		if err != nil {
			return err, nil
		}

		err, _, loc := service.token.RequestData("POST", service.config.IataServiceUrl, payload)
		if err != nil {
			return err, nil
		}
		shipLoc = append(shipLoc, loc)
	}

	slog.Info("Shipment imported in IATA OneRecord")
	return nil, shipLoc
}
func (service *IataService) createProductData(data *iata.Product) (error, string) {
	payload, err := json.Marshal(data)
	if err != nil {
		return err, ""
	}
	err, _, location := service.token.RequestData("POST", service.config.IataServiceUrl, payload)
	if err != nil {
		return err, ""
	}

	return nil, location
}
func (service *IataService) createHwbData(data []iata.Hwb, shipLoc string, orga map[string]string) error {
	hwb := data[0]
	for key, value := range orga {
		party := iata.InvolvedParty{}
		party.Type = "cargo:Party"
		party.Role = key
		party.Organization = iata.ObjectLink{Id: value}
		hwb.InvolvedParties = append(hwb.InvolvedParties, party)
	}
	hwb.Shipment = iata.ObjectLink{Id: shipLoc}
	payload, err := json.Marshal(hwb)
	if err != nil {
		return err
	}
	err, _, _ = service.token.RequestData("POST", service.config.IataServiceUrl, payload)
	if err != nil {
		return err
	}

	slog.Info("HWB imported in IATA OneRecord")
	return nil
}
func (service *IataService) createOrganisationData(data []iata.Organization) (error, map[string]string) {
	orgLoc := make(map[string]string)
	for _, org := range data {
		payload, err := json.Marshal(org)
		if err != nil {
			return err, nil
		}
		err, _, loc := service.token.RequestData("POST", service.config.IataServiceUrl, payload)
		if err != nil {
			return err, nil
		}
		orgLoc[org.Type] = loc
	}
	return nil, orgLoc
}

func (service *IataService) createPieceData(pieces []iata.Piece, itemLocations []string) (error, []string) {
	var loc []string
	for _, piece := range pieces {
		for _, itemLocation := range itemLocations {
			piece.ContainedItems = append(piece.ContainedItems, iata.ObjectLink{Id: itemLocation})
		}
		payload, err := json.Marshal(piece)
		if err != nil {
			return err, nil
		}
		err, _, location := service.token.RequestData("POST", service.config.IataServiceUrl, payload)
		if err != nil {
			return err, nil
		}
		loc = append(loc, location)
	}

	slog.Info("Pieces imported in IATA OneRecord")
	return nil, loc
}

func (service *IataService) createItemData(items []iata.Item) (error, []string) {
	var locations []string
	for _, item := range items {
		product := item.RawProduct
		err, productLocation := service.createProductData(&product)
		if err != nil {
			return err, nil
		}

		item.DescribedByProduct.Id = productLocation
		payload, err := json.Marshal(item)
		if err != nil {
			return err, nil
		}
		err, _, location := service.token.RequestData("POST", service.config.IataServiceUrl, payload)
		if err != nil {
			return err, nil
		}
		locations = append(locations, location)
	}

	slog.Info("Products imported in IATA OneRecord")
	slog.Info("Items imported in IATA OneRecord")
	return nil, locations
}
