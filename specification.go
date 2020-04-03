package placetypes

import (
	"strconv"
	"encoding/json"
	"github.com/whosonfirst/go-whosonfirst-placetypes/placetypes"
	"sync"
	"errors"
)

type WOFPlacetypeSpecification struct {
	catalog map[string]WOFPlacetype
	mu *sync.RWMutex
}

func Spec() (*WOFPlacetypeSpecification, error) {

	var catalog map[string]WOFPlacetype
	err := json.Unmarshal([]byte(placetypes.Specification), &catalog)

	if err != nil {
		return nil, err
	}

	mu := new(sync.RWMutex)
	
	spec := &WOFPlacetypeSpecification{
		catalog: catalog,
		mu: mu,
	}
	
	return spec, nil
}

func (spec *WOFPlacetypeSpecification) GetPlacetypeByName(name string) (*WOFPlacetype, error) {

	spec.mu.RLock()
	defer spec.mu.Unlock()
	
	for str_id, pt := range spec.catalog {

		if pt.Name == name {

			pt_id, err := strconv.Atoi(str_id)

			if err != nil {
				continue
			}

			pt_id64 := int64(pt_id)

			pt.Id = pt_id64
			return &pt, nil
		}
	}

	return nil, errors.New("Invalid placetype")
}

func (spec *WOFPlacetypeSpecification) GetPlacetypeById(id int64) (*WOFPlacetype, error) {

	spec.mu.RLock()
	defer spec.mu.Unlock()
	
	for str_id, pt := range spec.catalog {

		pt_id, err := strconv.Atoi(str_id)

		if err != nil {
			continue
		}

		pt_id64 := int64(pt_id)

		if pt_id64 == id {
			pt.Id = pt_id64
			return &pt, nil
		}
	}

	return nil, errors.New("Invalid placetype")
}

func (spec *WOFPlacetypeSpecification) AppendPlacetype(pt WOFPlacetype) error {

	spec.mu.Lock()
	defer spec.mu.Unlock()
	
	str_id := strconv.FormatInt(pt.Id, 10)
	
	_, exists := spec.catalog[str_id]

	if exists {
		return errors.New("Placetype ID already registered")
	}

	for _, pid := range pt.Parent {

		str_pid := strconv.FormatInt(pid, 10)

		_, exists := spec.catalog[str_pid]

		if !exists {
			return errors.New("Missing parent ID")
		}
	}

	spec.catalog[str_id] = pt
	return nil
}

func (spec *WOFPlacetypeSpecification) Catalog() map[string]WOFPlacetype {
	return spec.catalog
}
