package model

import "com.setlog/internal/model/iata"

type EntityCollection struct {
	Hwbs              []iata.Hwb
	Items             []iata.Item
	Organizations     []iata.Organization
	Persons           []iata.Person
	Pieces            []iata.Piece
	Products          []iata.Product
	Shipments         []iata.Shipment
	TransportMeans    []iata.TransportMeans
	TransportMovement []iata.TransportMovement
}
