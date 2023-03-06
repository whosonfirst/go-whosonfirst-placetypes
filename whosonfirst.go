package placetypes

import (
	"context"
	"fmt"
)

type WhosOnFirstFoo struct {
	Foo
	spec *WOFPlacetypeSpecification
	prop string
}

func NewWhosOnFirstFoo(ctx context.Context, uri string) (Foo, error) {

	spec, err := DefaultWOFPlacetypeSpecification()

	if err != nil {
		return nil, fmt.Errorf("Failed to create default WOF placetype specification, %w", err)
	}

	s := &WhosOnFirstFoo{
		spec: spec,
		prop: "wof:placetype",
	}

	return s, nil
}

func (s *WhosOnFirstFoo) Specification() *WOFPlacetypeSpecification {
	return s.spec
}

func (s *WhosOnFirstFoo) Property() string {
	return s.prop
}
