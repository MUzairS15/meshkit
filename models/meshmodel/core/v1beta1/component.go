package v1beta1

import (
	"encoding/json"

	"github.com/layer5io/meshkit/database"
	types "github.com/layer5io/meshkit/models/meshmodel/entity"
	"github.com/layer5io/meshkit/utils"

	"github.com/google/uuid"
)

type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion"`
}
type ComponentFormat string

const (
	JSON ComponentFormat = "JSON"
	YAML ComponentFormat = "YAML"
	CUE  ComponentFormat = "CUE"
)

type entity struct {
	TypeMeta
	Type    string `json:"type,omitempty" yaml:"type"`
	SubType string `json:"subType,omitempty" yaml:"subType"`
	Schema  string `json:"schema,omitempty" yaml:"schema"`
}

// swagger:response ComponentDefinition
// use NewComponent function for instantiating
type ComponentDefinition struct {
	ID                uuid.UUID              `json:"id,omitempty"`
	SchemaVersion     string                 `json:"schemaVersion,omitempty" yaml:"schemaVersion"`
	DefinitionVersion string                 `json:"definitionVersion,omitempty" yaml:"definitionVersion"`
	DisplayName       string                 `json:"displayName" gorm:"displayName"`
	Format            ComponentFormat        `json:"format" yaml:"format"`
	Model             Model                  `json:"model"`
	Metadata          map[string]interface{} `json:"metadata" yaml:"metadata"`
	// Rename this here and in schema as well, it creates confustion b/w entities like comps. relationship. model, here entity means Pod/Deployment....
	Entity entity `json:"entity,omitempty" yaml:"entity"`
}

type ComponentDefinitionDB struct {
	ID                uuid.UUID       `json:"id"`
	SchemaVersion     string          `json:"schemaVersion,omitempty" yaml:"schemaVersion"`
	DefinitionVersion string          `json:"definitionVersion,omitempty" yaml:"definitionVersion"`
	DisplayName       string          `json:"displayName" gorm:"displayName"`
	Format            ComponentFormat `json:"format" yaml:"format"`
	ModelID           uuid.UUID       `json:"-" gorm:"index:idx_component_definition_dbs_model_id,column:modelID"`
	Metadata          []byte          `json:"metadata" yaml:"metadata"`
	Entity entity `json:"entity,omitempty" yaml:"entity" gorm:"entity"`
}

func (c ComponentDefinition) Type() types.EntityType {
	return types.ComponentDefinition
}

func (c ComponentDefinition) GetID() uuid.UUID {
	return c.ID
}

func (c *ComponentDefinition) Create(db *database.Handler) (uuid.UUID, error) {
	c.ID = uuid.New()
	mid, err := c.Model.Create(db)
	if err != nil {
		return uuid.UUID{}, err
	}

	if !utils.IsSchemaEmpty(c.Entity.Schema) {
		c.Metadata["hasInvalidSchema"] = true
	}
	cdb := c.GetComponentDefinitionDB()
	cdb.ModelID = mid
	err = db.Create(&cdb).Error
	return c.ID, err
}

func (c *ComponentDefinition) GetComponentDefinitionDB() (cmd ComponentDefinitionDB) {
	cmd.ID = c.ID
	cmd.Entity.TypeMeta = c.Entity.TypeMeta
	cmd.Format = c.Format
	cmd.Metadata, _ = json.Marshal(c.Metadata)
	cmd.DisplayName = c.DisplayName
	cmd.Entity.Schema = c.Entity.Schema
	return
}

func (cmd *ComponentDefinitionDB) GetComponentDefinition(model Model) (c ComponentDefinition) {
	c.ID = cmd.ID
	c.Entity.TypeMeta = cmd.Entity.TypeMeta
	c.Format = cmd.Format
	c.DisplayName = cmd.DisplayName
	if c.Metadata == nil {
		c.Metadata = make(map[string]interface{})
	}
	_ = json.Unmarshal(cmd.Metadata, &c.Metadata)
	c.Entity.Schema = cmd.Entity.Schema
	c.Model = model
	return
}
