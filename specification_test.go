package placetypes

import (
	"testing"
)

func TestAppendSpecification(t *testing.T) {

	spec, err := DefaultWOFPlacetypeSpecification()

	if err != nil {
		t.Fatal(err)
	}

	parents := []int64{
		102312307, // country
	}

	pt := WOFPlacetype{
		Id:     1,
		Name:   "map",
		Role:   "optional",
		Parent: parents,
	}

	err = spec.AppendPlacetype(pt)

	if err != nil {
		t.Fatal(err)
	}

	_, err = spec.GetPlacetypeById(pt.Id)

	if err != nil {
		t.Fatal(err)
	}

	_, err = spec.GetPlacetypeByName(pt.Name)

	if err != nil {
		t.Fatal(err)
	}

}
