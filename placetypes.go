package placetypes

import (
	"encoding/json"
	"errors"
	"github.com/whosonfirst/go-whosonfirst-placetypes/placetypes"
	"strconv"
)

type WOFPlacetypeName struct {
	Lang string `json:"language"`
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type WOFPlacetypeAltNames map[string][]string

type WOFPlacetype struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Role   string  `json:"role"`
	Parent []int64 `json:"parent"`
	// AltNames []WOFPlacetypeAltNames		`json:"names"`
}

type WOFPlacetypeSpecification map[string]WOFPlacetype

var specification *WOFPlacetypeSpecification

func init() {

	var err error

	specification, err = Spec()

	if err != nil {
		log.Fatal("Failed to parse specification", err)
	}
}

func Spec() (*WOFPlacetypeSpecification, error) {

	places := placetypes.Spec()

	var spec WOFPlacetypeSpecification
	err := json.Unmarshal([]byte(placetypes.Specification), &spec)

	if err != nil {
		return nil, err
	}

	return &spec, nil
}

func GetPlacetypeByName(name string) (*WOFPlacetype, error) {

	for str_id, pt := range *specification {

		if pt.Name == name {

			id, _ := strconv.Atoi(str_id)
			pt.Id = int64(id)

			return &pt, nil
		}
	}

	return nil, errors.New("Invalid placetype")
}
