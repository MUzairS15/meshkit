package entity

import "github.com/google/uuid"

type EntityType string

const (
	ComponentDefinition    EntityType = "component"
	PolicyDefinition       EntityType = "policy"
	RelationshipDefinition EntityType = "relationship"
	Model                  EntityType = "model"
)

// Each entity will have it's own Filter implementation via which it exposes the nobs and dials to fetch entities
type Filter interface {
	Create(map[string]interface{})
}

type Entity interface {
	Create() (entityID uuid.UUID, err error)
	Get() (entityID uuid.UUID, err error)
}
