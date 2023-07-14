// NewPatternFileFromCytoscapeJSJSON takes in CytoscapeJS JSON
// and creates a PatternFile from it.
// This function always returns meshkit error
package cytoscape

import (
	"encoding/json"
	"fmt"
	"strings"

	cytoscapejs "gonum.org/v1/gonum/graph/formats/cytoscapejs"
	mathrand "math/rand"

	"github.com/gofrs/uuid"
)

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


func NewPatternFileFromCytoscapeJSJSON(name string, byt []byte) (Pattern, error) {
	// Unmarshal data into cytoscape struct
	var cy cytoscapejs.GraphElem
	if err := json.Unmarshal(byt, &cy); err != nil {
		return Pattern{}, ErrPatternFromCytoscape(err)
	}
	if name == "" {
		name = "MesheryGeneratedPattern"
	}
	// Convert cytoscape struct to patternfile
	pf := Pattern{
		Name:     name,
		Services: make(map[string]*Service),
	}
	dependsOnMap := make(map[string][]string, 0) //used to figure out dependencies from traits.meshmap.parent
	eleToSvc := make(map[string]string)          //used to map cyto element ID uniquely to the name of the service created.
	countDuplicates := make(map[string]int)
	//store the names of services and their count
	err := processCytoElementsWithPattern(cy.Elements, func(svc Service, ele cytoscapejs.Element) error {
		name := svc.Name
		countDuplicates[name]++
		return nil
	})
	if err != nil {
		return pf, ErrPatternFromCytoscape(err)
	}

	//Populate the dependsOn field with appropriate unique service names
	err = processCytoElementsWithPattern(cy.Elements, func(svc Service, ele cytoscapejs.Element) error {
		//Extract parents, if present
		m, ok := svc.Traits["meshmap"].(map[string]interface{})
		if ok {
			parentID, ok := m["parent"].(string)
			if ok { //If it does not have a parent then we can skip and we dont make it depend on anything
				elementID, ok := m["id"].(string)
				if !ok {
					return fmt.Errorf("required meshmap trait field: \"id\" missing")
				}
				dependsOnMap[elementID] = append(dependsOnMap[elementID], parentID)
			}
		}

		//Only make the name unique when duplicates are encountered. This allows clients to preserve and propagate the unique name they want to give to their workload
		uniqueName := svc.Name
		if countDuplicates[uniqueName] > 1 {
			//set appropriate unique service name
			uniqueName = strings.ToLower(svc.Name)
			uniqueName += "-" + GetRandomAlphabetsOfDigit(5)
		}
		eleToSvc[ele.Data.ID] = uniqueName //will be used while adding depends-on
		pf.Services[uniqueName] = &svc
		return nil
	})
	if err != nil {
		return pf, ErrPatternFromCytoscape(err)
	}
	//add depends-on field
	for child, parents := range dependsOnMap {
		childSvc := eleToSvc[child]
		if childSvc != "" {
			for _, parent := range parents {
				if eleToSvc[parent] != "" {
					pf.Services[childSvc].DependsOn = append(pf.Services[childSvc].DependsOn, eleToSvc[parent])
				}
			}
		}
	}
	return pf, nil
}


// processCytoElementsWithPattern iterates over all the cyto elements, convert each into a patternfile service and exposes a callback to handle that service
func processCytoElementsWithPattern(eles []cytoscapejs.Element, callback func(svc Service, ele cytoscapejs.Element) error) error {
	for _, elem := range eles {
		// Try to create Service object from the elem.scratch's _data field
		// if this fails then immediately fail the process and return an error
		castedScratch, ok := elem.Scratch.(map[string]interface{})
		if !ok {
			return fmt.Errorf("empty scratch field is not allowed, must contain \"_data\" field holding metadata")
		}

		data, ok := castedScratch["_data"]
		if !ok {
			return fmt.Errorf("\"_data\" cannot be empty")
		}

		// Convert data to JSON for easy serialization
		svcByt, err := json.Marshal(&data)
		if err != nil {
			return fmt.Errorf("failed to serialize service from the metadata in the scratch")
		}

		// Unmarshal the JSON into a service
		svc := Service{
			Settings: map[string]interface{}{},
			Traits:   map[string]interface{}{},
		}

		// Add meshmap position
		svc.Traits["meshmap"] = map[string]interface{}{
			"position": map[string]float64{
				"posX": elem.Position.X,
				"posY": elem.Position.Y,
			},
		}

		if err := json.Unmarshal(svcByt, &svc); err != nil {
			return fmt.Errorf("failed to create service from the metadata in the scratch")
		}
		if svc.Name == "" {
			return fmt.Errorf("cannot save service with empty name")
		}
		err = callback(svc, elem)
		if err != nil {
			return err
		}
	}
	return nil
}


func GetRandomAlphabetsOfDigit(length int) (s string) {
	charSet := "abcdedfghijklmnopqrstuvwxyz"
	for i := 0; i < length; i++ {
		random := mathrand.Intn(len(charSet))
		randomChar := charSet[random]
		s += string(randomChar)
	}
	return
}
