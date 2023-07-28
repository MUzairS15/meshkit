package schema

// import (
// 	"encoding/json"
// 	"errors"
// 	"cuelang.org/go/cue"
// 	"github.com/layer5io/meshkit/utils"
// 	"github.com/layer5io/meshkit/utils/manifests"
// )

// // all paths should be a valid CUE expression
// type CuePathConfig struct {
// 	NamePath       string
// 	GroupPath      string
// 	VersionPath    string
// 	SpecPath       string
// 	ScopePath      string
// 	PropertiesPath string
// 	// identifiers are the values that uniquely identify a CRD (in most of the cases, it is the 'Name' field)
// 	IdentifierPath string
// }

// // Remove the fields which is either not required by end user (like status) or is prefilled by system (like apiVersion, kind and metadata)
// var fieldsToDelete = [4]string{"apiVersion", "kind", "status", "metadata"}

// // extracts the JSONSCHEMA of the CRD and outputs the json encoded string of the schema
// func GetSchema(parsedCrd cue.Value, pathConf CuePathConfig) (string, error) {
// 	schema := map[string]interface{}{}
// 	specCueVal, err := utils.Lookup(parsedCrd, pathConf.SpecPath)
// 	if err != nil {
// 		return "", err
// 	}
// 	marshalledJson, err := specCueVal.MarshalJSON()
// 	if err != nil {
// 		return "", ErrGetSchema(err)
// 	}
// 	err = json.Unmarshal(marshalledJson, &schema)
// 	if err != nil {
// 		return "", ErrGetSchema(err)
// 	}
// 	resourceId, err := extractCueValueFromPath(parsedCrd, pathConf.IdentifierPath)
// 	if err != nil {
// 		return "", ErrGetSchema(err)
// 	}

// 	updatedProps, err := UpdateProperties(specCueVal, cue.ParsePath(pathConf.PropertiesPath), resourceId)

// 	if err != nil {
// 		return "", err
// 	}

// 	schema = updatedProps
// 	DeleteFields(schema)

// 	(schema)["title"] = manifests.FormatToReadableString(resourceId)
// 	var output []byte
// 	output, err = json.MarshalIndent(schema, "", " ")
// 	if err != nil {
// 		return "", ErrGetSchema(err)
// 	}
// 	return string(output), nil
// }

// func extractCueValueFromPath(crd cue.Value, pathConf string) (string, error) {
// 	cueRes, err := utils.Lookup(crd, pathConf)
// 	if err != nil {
// 		return "", err
// 	}
// 	res, err := cueRes.String()
// 	if err != nil {
// 		return "", err
// 	}
// 	return res, nil
// }

// // function to remove fields that are not required or prefilled
// func DeleteFields(m map[string]interface{}) {
// 	key := "properties"
// 	if m[key] == nil {
// 		return
// 	}
// 	if prop, ok := m[key].(map[string]interface{}); ok && prop != nil {
// 		for _, f := range fieldsToDelete {
// 			delete(prop, f)
// 		}
// 		m[key] = prop
// 	}
// }

// /*
// Find and modify specific schema properties.
// 1. Identify interesting properties by walking entire schema.
// 2. Store path to interesting properties. Finish walk.
// 3. Iterate all paths and modify properties.
// 5. If error occurs, return nil and skip modifications.
// */
// func UpdateProperties(fieldVal cue.Value, cuePath cue.Path, group string) (map[string]interface{}, error) {
// 	rootPath := fieldVal.Path().Selectors()

// 	compProperties := fieldVal.LookupPath(cuePath)
// 	crd, err := fieldVal.MarshalJSON()
// 	if err != nil {
// 		return nil, ErrUpdateSchema(err, group)
// 	}

// 	modified := make(map[string]interface{})
// 	pathSelectors := [][]cue.Selector{}

// 	err = json.Unmarshal(crd, &modified)
// 	if err != nil {
// 		return nil, ErrUpdateSchema(err, group)
// 	}

// 	compProperties.Walk(func(c cue.Value) bool {
// 		return true
// 	}, func(c cue.Value) {
// 		val := c.LookupPath(cue.ParsePath(`"x-kubernetes-preserve-unknown-fields"`))
// 		if val.Exists() {
// 			child := val.Path().Selectors()
// 			childM := child[len(rootPath):(len(child) - 1)]
// 			pathSelectors = append(pathSelectors, childM)
// 		}
// 	})

// 	// "pathSelectors" contains all the paths from root to the property which needs to be modified.
// 	for _, selectors := range pathSelectors {
// 		var m interface{}
// 		m = modified
// 		index := 0

// 		for index < len(selectors) {
// 			selector := selectors[index]
// 			selectorType := selector.Type()
// 			s := selector.String()
// 			if selectorType == cue.IndexLabel {
// 				t, ok := m.([]interface{})
// 				if !ok {
// 					return nil, ErrUpdateSchema(errors.New("error converting to []interface{}"), group)
// 				}
// 				token := selector.Index()
// 				m, ok = t[token].(map[string]interface{})
// 				if !ok {
// 					return nil, ErrUpdateSchema(errors.New("error converting to map[string]interface{}"), group)
// 				}
// 			} else {
// 				if selectorType == cue.StringLabel {
// 					s = selector.Unquoted()
// 				}
// 				t, ok := m.(map[string]interface{})
// 				if !ok {
// 					return nil, ErrUpdateSchema(errors.New("error converting to map[string]interface{}"), group)
// 				}
// 				m = t[s]
// 			}
// 			index++
// 		}

// 		t, ok := m.(map[string]interface{})
// 		if !ok {
// 			return nil, ErrUpdateSchema(errors.New("error converting to map[string]interface{}"), group)
// 		}
// 		delete(t, "x-kubernetes-preserve-unknown-fields")
// 		if m == nil {
// 			m = make(map[string]interface{})
// 		}
// 		t["type"] = "string"
// 		t["format"] = "textarea"
// 	}
// 	return modified, nil
// }