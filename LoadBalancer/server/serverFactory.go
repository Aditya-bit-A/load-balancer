package server

import (
	"loadbalancer/models"
	"log"
)

// server/factory.go
func NewServerInstance(id string, serverType string) models.ServerMetaData {
	hash := models.SH(id, 1)

	switch serverType {
	case "Type1":
		inst := &models.ServerInstance{
			ID:      id,
			HashVal: hash,
		}

		return inst
	case "Type2":
		inst := &models.ServerInstance{
			ID:      id,
			HashVal: hash,
		}

		return inst
	default:
		log.Fatalf("Unsupported server type: %s", serverType)
		return nil
	}
}
