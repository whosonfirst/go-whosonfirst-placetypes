package placetypes

import (
	"log"
)

type WOFPlacetypeName struct {
	// Lang is the RFC 5646 (BCP-47) language tag for the placetype name
	Lang string `json:"language"`
	Kind string `json:"kind"`
	// Name is the name of the placetype (in the language defined by `Lang`)
	Name string `json:"name"`
}

type WOFPlacetypeAltNames map[string][]string

// Type WOFPlacetype defines an individual placetype encoded in a `WOFPlacetypeSpecification`
// instance. The choice of naming this "WOFPlacetype" is unfortunate because since it is easily
// confused with the actual JSON definition files for placetypes. However, we're stuck with it
// for now in order to preserve backwards compatibility. Womp womp...
type WOFPlacetype struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Role   string  `json:"role"`
	Parent []int64 `json:"parent"`
	// AltNames []WOFPlacetypeAltNames		`json:"names"`
}

func (pt *WOFPlacetype) String() string {
	return pt.Name
}

var specification *WOFPlacetypeSpecification

func init() {

	s, err := DefaultWOFPlacetypeSpecification()

	if err != nil {
		log.Fatal("Failed to load default WOF specification", err)
	}

	specification = s
}

func GetPlacetypeByName(name string) (*WOFPlacetype, error) {
	return specification.GetPlacetypeByName(name)
}

func GetPlacetypeById(id int64) (*WOFPlacetype, error) {
	return specification.GetPlacetypeById(id)
}

func AppendPlacetype(pt WOFPlacetype) error {
	return specification.AppendPlacetype(pt)
}

func AppendPlacetypeSpecification(spec *WOFPlacetypeSpecification) error {
	return specification.AppendPlacetypeSpecification(spec)
}

// Placetypes returns all the known placetypes for the 'common', 'optional' and 'common_optional' roles.
func Placetypes() ([]*WOFPlacetype, error) {

	roles := []string{
		COMMON_ROLE,
		OPTIONAL_ROLE,
		COMMON_OPTIONAL_ROLE,
		CUSTOM_ROLE,
	}

	return PlacetypesForRoles(roles)
}

func PlacetypesForRoles(roles []string) ([]*WOFPlacetype, error) {
	return specification.PlacetypesForRoles(roles)
}

// IsValidPlacetypeId returns a boolean value indicating whether 'name' is a known and valid placetype name.
func IsValidPlacetype(name string) bool {
	return specification.IsValidPlacetype(name)
}

// IsValidPlacetypeId returns a boolean value indicating whether 'id' is a known and valid placetype ID.
func IsValidPlacetypeId(id int64) bool {
	return specification.IsValidPlacetypeId(id)
}

// Returns true is 'b' is an ancestor of 'a'.
func IsAncestor(a *WOFPlacetype, b *WOFPlacetype) bool {
	return specification.IsAncestor(a, b)
}

// Returns true is 'b' is a descendant of 'a'.
func IsDescendant(a *WOFPlacetype, b *WOFPlacetype) bool {
	return specification.IsDescendant(a, b)
}

// Children returns the immediate child placetype of 'pt'.
func Children(pt *WOFPlacetype) []*WOFPlacetype {
	return specification.Children(pt)
}

// Descendants returns the descendants of role "common" for 'pt'.
func Descendants(pt *WOFPlacetype) []*WOFPlacetype {
	return specification.Descendants(pt)
}

// DescendantsForRoles returns the descendants matching any role in 'roles' for 'pt'.
func DescendantsForRoles(pt *WOFPlacetype, roles []string) []*WOFPlacetype {
	return specification.DescendantsForRoles(pt, roles)
}

// Ancestors returns the ancestors of role "common" for 'pt'.
func Ancestors(pt *WOFPlacetype) []*WOFPlacetype {
	return AncestorsForRoles(pt, []string{"common"})
}

// AncestorsForRoles returns the ancestors matching any role in 'roles' for 'pt'.
func AncestorsForRoles(pt *WOFPlacetype, roles []string) []*WOFPlacetype {
	return specification.AncestorsForRoles(pt, roles)
}
