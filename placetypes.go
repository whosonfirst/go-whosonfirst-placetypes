package placetypes

import (
	"encoding/json"
	"github.com/whosonfirst/go-whosonfirst-placetypes/spec"
)

type WOFPlacetypeSpec map[string]WOFPlacetype

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

func Init() (*WOFPlacetypeSpec, error) {

	places := placetypes.Spec()

	var spec WOFPlacetypeSpec
	err := json.Unmarshal([]byte(places), &spec)

	if err != nil {
		return nil, err
	}

	return &spec, nil
}

/*
func (sp *WOFPlacetypeSpec) Common() []string {
	return sp.WithRole("common")
}

func (sp *WOFPlacetypeSpec) CommonOptional() []string {
	return sp.WithRole("common_optional")
}

func (sp *WOFPlacetypeSpec) Optional() []string {
	return sp.WithRole("optional")
}

func (sp *WOFPlacetypeSpec) WithRole(role string) []string {

	places := make([]string, 0)

	for id, placetype := range sp {
		if placetype.Role != role {
		   continue
		}

		places = append(places, role)
	}

	return places
}

func (sp *WOFPlacetypeSpec) WithRoles(roles []string) []string {

	places := make([]string, 0)

	for _, role := range roles {
	    places = append(places, sp.WithRole(role))
	}

	return places
}

func IsValidRole(role string) bool {

	return false
}
*/
