package patterns

import "github.com/gofrs/uuid"

type Pattern struct {
	// Name is the human-readable, display-friendly descriptor of the pattern
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
	//Vars will be used to configure the pattern when it is imported from other patterns.
	Vars map[string]interface{} `yaml:"vars,omitempty" json:"vars,omitempty"`
	// PatternID is the moniker use to uniquely identify any given pattern
	// Convention: SMP-###-v#.#.#
	PatternID string              `yaml:"patternID,omitempty" json:"patternID,omitempty"`
	Services  map[string]*Service `yaml:"services,omitempty" json:"services,omitempty"`
}

// Service represents the services defined within the appfile
type Service struct {
	// ID is the id of the service and is completely internal to
	// Meshery Server and meshery providers
	ID *uuid.UUID `yaml:"id,omitempty" json:"id,omitempty"`
	// Name is the name of the service and is an optional parameter
	// If given then this supercedes the name of the service inherited
	// from the parent
	Name         string            `yaml:"name,omitempty" json:"name,omitempty"`
	Type         string            `yaml:"type,omitempty" json:"type,omitempty"`
	APIVersion   string            `yaml:"apiVersion,omitempty" json:"apiVersion,omitempty"`
	Namespace    string            `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Version      string            `yaml:"version,omitempty" json:"version,omitempty"`
	Model        string            `yaml:"model,omitempty" json:"model,omitempty"`
	IsAnnotation bool              `yaml:"isAnnotation,omitempty" json:"isAnnotation,omitempty"`
	Labels       map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations  map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	// DependsOn correlates one or more objects as a required dependency of this service
	// DependsOn is used to determine sequence of operations
	DependsOn []string `yaml:"dependsOn,omitempty" json:"dependsOn,omitempty"`

	Settings map[string]interface{} `yaml:"settings,omitempty" json:"settings,omitempty"`
	Traits   map[string]interface{} `yaml:"traits,omitempty" json:"traits,omitempty"`
}
