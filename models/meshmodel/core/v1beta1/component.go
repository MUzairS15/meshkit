package v1beta1

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"github.com/layer5io/meshkit/utils"

	"github.com/google/uuid"
)

type VersionMeta struct {
	SchemaVersion string `json:"schemaVersion,omitempty" yaml:"schemaVersion"`
	Version       string `json:"version,omitempty" yaml:"version"`
}

type TypeMeta struct {
	Kind    string `json:"kind,omitempty" yaml:"kind"`
	Version string `json:"version,omitempty" yaml:"version"`
}

type ComponentFormat string

const (
	JSON ComponentFormat = "JSON"
	YAML ComponentFormat = "YAML"
	CUE  ComponentFormat = "CUE"
)

type component struct {
	TypeMeta
	Schema string `json:"schema,omitempty" yaml:"schema"`
}

// swagger:response ComponentDefinition
// use NewComponent function for instantiating
type ComponentDefinition struct {
	ID uuid.UUID `json:"id,omitempty"`
	VersionMeta
	DisplayName string                 `json:"displayName" gorm:"displayName"`
	Description string                 `json:"description" gorm:"description"`
	Format      ComponentFormat        `json:"format" yaml:"format"`
	Model       Model                  `json:"model"`
	Metadata    map[string]interface{} `json:"metadata" yaml:"metadata"`
	// component corresponds to the specifications of underlying entity eg: Pod/Deployment....
	Component component `json:"component,omitempty" yaml:"component"`
}

type ComponentDefinitionDB struct {
	ID uuid.UUID `json:"id"`
	VersionMeta
	DisplayName string          `json:"displayName" gorm:"displayName"`
	Description string          `json:"description" gorm:"description"`
	Format      ComponentFormat `json:"format" yaml:"format"`
	ModelID     uuid.UUID       `json:"-" gorm:"index:idx_component_definition_dbs_model_id,column:modelID"`
	Metadata    []byte          `json:"metadata" yaml:"metadata"`
	Component   component       `json:"component,omitempty" yaml:"component" gorm:"component"`
}

func (c ComponentDefinition) Type() entity.EntityType {
	return entity.ComponentDefinition
}

func (c ComponentDefinition) GetID() uuid.UUID {
	return c.ID
}

func (c *ComponentDefinition) GetEntityDetail() string {
	return fmt.Sprintf("type: %s, definition version: %s, name: %s, model: %s, version: %s", c.Type(), c.Version, c.DisplayName, c.Model.Name, c.Model.Version)
}

func (c *ComponentDefinition) Create(db *database.Handler, hostID uuid.UUID) (uuid.UUID, error) {
	c.ID = uuid.New()

	isAnnotation, _ := c.Metadata["isAnnotation"].(bool)

	if c.Component.Schema == "" && !isAnnotation { //For components which has an empty schema and is not an annotation, return error
		// return ErrEmptySchema()
		return uuid.Nil, nil
	}

	mid, err := c.Model.Create(db, hostID)
	if err != nil {
		return uuid.UUID{}, err
	}

	if !utils.IsSchemaEmpty(c.Component.Schema) {
		c.Metadata["hasInvalidSchema"] = true
	}
	cdb := c.GetComponentDefinitionDB()
	cdb.ModelID = mid
	err = db.Create(&cdb).Error
	return c.ID, err
}

func (m *ComponentDefinition) UpdateStatus(db *database.Handler, status entity.EntityStatus) error {
	return nil
}

func (c *ComponentDefinition) GetComponentDefinitionDB() (cmd ComponentDefinitionDB) {
	// cmd.ID = c.ID id will be assigned by the database itself don't use this, as it will be always uuid.nil, because id is not known when comp gets generated.
	// While database creates an entry with valid primary key but to avoid confusion, it is disabled and accidental assignment of custom id.
	cmd.VersionMeta = c.VersionMeta
	cmd.DisplayName = c.DisplayName
	cmd.Description = c.Description
	cmd.Format = c.Format
	cmd.ModelID = c.Model.ID
	cmd.Metadata, _ = json.Marshal(c.Metadata)
	cmd.Component = c.Component
	return
}

func (c ComponentDefinition) WriteComponentDefinition(componentDirPath string) error {
	componentPath := filepath.Join(componentDirPath, c.Component.Kind+".json")
	err := utils.WriteJSONToFile[ComponentDefinition](componentPath, c)
	return err
}

func (cmd *ComponentDefinitionDB) GetComponentDefinition(model Model) (c ComponentDefinition) {
	c.ID = cmd.ID
	c.VersionMeta = cmd.VersionMeta
	c.DisplayName = cmd.DisplayName
	c.Description = cmd.Description
	c.Format = cmd.Format
	c.Model = model
	if c.Metadata == nil {
		c.Metadata = make(map[string]interface{})
	}
	_ = json.Unmarshal(cmd.Metadata, &c.Metadata)
	c.Component = cmd.Component
	return
}
