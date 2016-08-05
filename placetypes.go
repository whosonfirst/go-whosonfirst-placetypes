package placetypes

import (
	"encoding/json"
	"errors"
	"github.com/whosonfirst/go-whosonfirst-placetypes/spec"
)

type WOFPlacetypeName struct {
	Lang string
	Kind string
	Name string
}

type WOFPlacetypeAltNames struct {
	Lang string
	Name string
}

type WOFPlacetype struct {
	Id       int64
	Name     string
	Role     string
	Parent   []string
	AltNames []WOFPlacetypeAltNames
}

func Init() (interface{}, error) {

	spec := placetypes.Spec()

	var d interface{}
	err := json.Unmarshal([]byte(spec), &d)

	if err != nil {
		return nil, err
	}

	return d, nil
}

func NewPlacetypeName(name string) (*WOFPlacetypeName, error) {

	return nil, errors.New("Please write me...")
}

func NewPlacetype(placetype string) (*WOFPlacetype, error) {

	return nil, errors.New("Please write me...")
}

func IsValidPlacetype(placetype string) bool {
	return false
}

func Common() []string {
	return WithRole("common")
}

func CommonOptional() []string {
	return WithRole("common_optional")
}

func Optional() []string {
	return WithRole("optional")
}

func WithRole(role string) []string {

	places := make([]string, 0)
	return places
}

func WithRoles(roles []string) []string {

	places := make([]string, 0)
	return places
}

func IsValidRole(role string) bool {

	return false
}

func (p WOFPlacetype) Names() []*WOFPlacetypeName {

	names := make([]*WOFPlacetypeName, 0)
	return names
}

func (p WOFPlacetype) Parents() []*WOFPlacetype {

	places := make([]*WOFPlacetype, 0)
	return places
}

func (p WOFPlacetype) Ancestors() []string {

	places := make([]string, 0)
	return places
}

func (p WOFPlacetype) AncestorsWithRoles([]string) []string {

	places := make([]string, 0)
	return places
}
