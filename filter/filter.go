package filter

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/warning"
)

type PlacetypesFilter struct {
	required  map[string]*placetypes.WOFPlacetype
	forbidden map[string]*placetypes.WOFPlacetype
}

func NewPlacetypesFilter(include []string, include_roles []string, exclude []string) (*PlacetypesFilter, error) {

	required := make(map[string]*placetypes.WOFPlacetype)
	forbidden := make(map[string]*placetypes.WOFPlacetype)

	for _, p := range include {

		_, ok := required[p]

		if ok {
			continue
		}

		pt, err := placetypes.GetPlacetypeByName(p)

		if err != nil {
			return nil, fmt.Errorf("Failed to get placetype by name for inclusion %s, %w", p, err)
		}

		required[p] = pt
	}

	if len(include_roles) > 0 {
		return nil, fmt.Errorf("included roles are not supported yet")
	}

	for _, p := range exclude {

		_, ok := forbidden[p]

		if ok {
			continue
		}

		pt, err := placetypes.GetPlacetypeByName(p)

		if err != nil {
			return nil, fmt.Errorf("Failed to get placetype by name for exclusion %s, %w", p, err)
		}

		forbidden[p] = pt
	}

	f := PlacetypesFilter{
		required:  required,
		forbidden: forbidden,
	}

	return &f, nil
}

func (f *PlacetypesFilter) AllowFromString(pt_str string) (bool, error) {

	pt, err := placetypes.GetPlacetypeByName(pt_str)

	if err != nil {
		e := fmt.Errorf("Failed to get placetype by name %s, %w", pt_str, err)
		return true, warning.New(e.Error())
	}

	return f.Allow(pt)
}

func (f *PlacetypesFilter) Allow(pt *placetypes.WOFPlacetype) (bool, error) {

	if len(f.forbidden) > 0 {

		_, ok := f.forbidden[pt.Name]

		if ok {
			return false, nil
		}

	}

	if len(f.required) > 0 {

		_, ok := f.required[pt.Name]

		if !ok {
			return false, nil
		}
	}

	return true, nil
}
