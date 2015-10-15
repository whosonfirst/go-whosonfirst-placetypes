package placetypes

import (

)

type WOFPlacetypeName struct {
     Lang string
     Kind string
     Name string
}

type WOFPlacetype struct {
     Id int
     Name string
     Role string
     Parent []string
}

func NewPlacetypeName(name string) *WOFPlacetypeName, error {

}

func NewPlacetype(placetype string) *WOFPlacetype, error {

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

}

func WithRoles(roles []string) []string {

}

func IsValidRole(role string) bool {

}

func (p WOFPlacetype) Names() []*WOFPlacetypeName {

}

func (p WOFPlacetype) Parents() []*WOFPlacetype {

}

func (p WOFPlacetype) Ancestors() []string {

}

func (p WOFPlacetype) AncestorsWithRoles([]string) []string {

}