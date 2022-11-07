package v1alpha1

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/models/meshmodel/core/types"
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

// use NewComponent function for instantiating
type ComponentDefinition struct {
	ID        uuid.UUID `json:"-"`
	TypeMeta  `gorm:"embedded" yaml:"typemeta"`
	Format    ComponentFormat   `gorm:"format"`
	Metadata  ComponentMetadata `gorm:"-"`
	Schema    string            `gorm:"embedded" yaml:"schema"`
	CreatedAt time.Time         `json:"-"`
	UpdatedAt time.Time         `json:"-"`
}

func (c ComponentDefinition) Type() types.CapabilityType {
	return types.ComponentDefinition
}

func CreateComponent(db *database.Handler, c ComponentDefinition) (uuid.UUID, error) {
	c.ID = uuid.New()
	c.Metadata.ID = uuid.New()
	compMeta := c.Metadata
	err := db.Create(&compMeta).Error
	if err != nil {
		return uuid.UUID{}, err
	}
	err = db.Create(&c).Error
	return c.ID, err
}
func GetComponents(db *database.Handler, f ComponentFilter) (c []ComponentDefinition) {
	if f.ModelName != "" {
		var metas []ComponentMetadata
		_ = db.Where("model = ?", f.ModelName).Find(&metas).Error
		var ids []uuid.UUID
		mapIDsToComponentsMetadata := make(map[uuid.UUID]ComponentMetadata)
		for _, m := range metas {
			ids = append(ids, m.ComponentID)
			mapIDsToComponentsMetadata[m.ComponentID] = m
		}
		var ctemp []ComponentDefinition
		_ = db.Where("id IN ?", ids).Where("name = ?", f.Name).Find(&ctemp).Error
		for _, comp := range ctemp {
			comp.Metadata = mapIDsToComponentsMetadata[comp.ID]
			c = append(c, comp)
		}
	}

	return
}

type ComponentFilter struct {
	Name      string
	ModelName string
}

// Create the filter from map[string]interface{}
func (cf *ComponentFilter) Create(m map[string]interface{}) {
	if m == nil {
		return
	}
	cf.Name = m["name"].(string)
}

type ComponentMetadata struct {
	ID          uuid.UUID `json:"-"`
	ComponentID uuid.UUID `json:"-"`
	Model       string
	Version     string
	Category    string
	SubCategory string
	Metadata    map[string]interface{}
}

// This struct is internal to the system
type componentMetadataDB struct {
	ID          uuid.UUID
	ComponentID uuid.UUID
	Model       string
	Version     string
	Category    string
	SubCategory string
	Metadata    []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (cmd *componentMetadataDB) ToComponentMetadata() (c ComponentMetadata) {
	c.ID = cmd.ID
	c.ComponentID = cmd.ComponentID
	c.Model = cmd.Model
	c.Version = cmd.Version
	c.Category = cmd.Category
	c.SubCategory = cmd.SubCategory

	byt, _ := json.Marshal(cmd.Metadata)
	_ = json.Unmarshal(byt, &c.Metadata)
	return
}
func (cmd *componentMetadataDB) FromComponentMetadata(c ComponentMetadata) {
	cmd.ID = c.ID
	cmd.ComponentID = c.ComponentID
	cmd.Model = c.Model
	cmd.Version = c.Version
	cmd.Category = c.Category
	cmd.SubCategory = c.SubCategory

	byt, _ := json.Marshal(c.Metadata)
	_ = json.Unmarshal(byt, &cmd.Metadata)
	return
}
