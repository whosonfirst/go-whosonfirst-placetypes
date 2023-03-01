package placetypes

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/dominikbraun/graph"
)

//go:embed placetypes.json
var FS embed.FS

type WOFPlacetypeSpecification struct {
	catalog                map[string]WOFPlacetype
	mu                     *sync.RWMutex
	relationships          *sync.Map
	indexing_relationships int32
	indexing_channel       chan bool
}

func DefaultWOFPlacetypeSpecification() (*WOFPlacetypeSpecification, error) {

	r, err := FS.Open("placetypes.json")

	if err != nil {
		return nil, fmt.Errorf("Failed to open placetypes, %w", err)
	}

	defer r.Close()

	spec, err := NewWOFPlacetypeSpecificationWithReader(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to load placetypes, %w", err)
	}

	go spec.indexRelationships()

	return spec, nil
}

func NewWOFPlacetypeSpecificationWithReader(r io.Reader) (*WOFPlacetypeSpecification, error) {

	var catalog map[string]WOFPlacetype

	dec := json.NewDecoder(r)
	err := dec.Decode(&catalog)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode reader, %w", err)
	}

	mu := new(sync.RWMutex)

	relationships := new(sync.Map)

	indexing_channel := make(chan bool)

	go func() {
		indexing_channel <- true
	}()

	spec := &WOFPlacetypeSpecification{
		catalog:          catalog,
		mu:               mu,
		relationships:    relationships,
		indexing_channel: indexing_channel,
	}

	return spec, nil
}

func NewWOFPlacetypeSpecification(body []byte) (*WOFPlacetypeSpecification, error) {

	r := bytes.NewReader(body)
	return NewWOFPlacetypeSpecificationWithReader(r)
}

func (spec *WOFPlacetypeSpecification) isIndexingRelationships() bool {

	if atomic.LoadInt32(&spec.indexing_relationships) > 0 {
		return true
	}

	return false
}

func (spec *WOFPlacetypeSpecification) indexRelationships() {

	<-spec.indexing_channel

	atomic.AddInt32(&spec.indexing_relationships, 1)

	defer func() {
		atomic.AddInt32(&spec.indexing_relationships, -1)
		spec.indexing_channel <- true
	}()

	spec.relationships = new(sync.Map)

	roles := []string{
		COMMON_ROLE,
		OPTIONAL_ROLE,
		COMMON_OPTIONAL_ROLE,
		CUSTOM_ROLE,
	}

	count_roles := len(roles)

	for i := 0; i < count_roles; i++ {

		pt_roles := roles[0:i]

		for _, pt := range spec.catalog {
			go spec.Children(&pt)
			go spec.DescendantsForRoles(&pt, pt_roles)
			go spec.AncestorsForRoles(&pt, pt_roles)
		}
	}

}

func (spec *WOFPlacetypeSpecification) PlacetypesForRoles(roles []string) ([]*WOFPlacetype, error) {

	pl, err := spec.GetPlacetypeByName("planet")

	if err != nil {
		return nil, fmt.Errorf("Failed to load 'planet' placetype, %w", err)
	}

	pt_list := spec.DescendantsForRoles(pl, roles)

	pt_list = append([]*WOFPlacetype{pl}, pt_list...)
	return pt_list, nil
}

func (spec *WOFPlacetypeSpecification) GetPlacetypeByName(name string) (*WOFPlacetype, error) {

	// spec.mu.RLock()
	// defer spec.mu.RUnlock()

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

	return nil, fmt.Errorf("Invalid placetype")
}

func (spec *WOFPlacetypeSpecification) GetPlacetypeById(id int64) (*WOFPlacetype, error) {

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

	return nil, fmt.Errorf("Invalid placetype")
}

func (spec *WOFPlacetypeSpecification) AppendPlacetypeSpecification(other_spec *WOFPlacetypeSpecification) error {

	if spec.isIndexingRelationships() {

		<-spec.indexing_channel

		go func() {
			spec.indexing_channel <- true
		}()
	}

	for _, pt := range other_spec.Catalog() {

		err := spec.AppendPlacetype(pt)

		if err != nil {
			return fmt.Errorf("Failed to append placetype %v, %w", pt, err)
		}
	}

	// Note the way we are not doing this in a Go routine; we want to block
	// until things the relationships between the two specifications have been
	// updated.

	spec.indexRelationships()

	return nil
}

func (spec *WOFPlacetypeSpecification) AppendPlacetype(pt WOFPlacetype) error {

	spec.mu.Lock()
	defer spec.mu.Unlock()

	existing_pt, _ := spec.GetPlacetypeById(pt.Id)

	if existing_pt != nil {
		return fmt.Errorf("Placetype ID %d (%s) already registered", pt.Id, pt.Name)
	}

	existing_pt, _ = spec.GetPlacetypeByName(pt.Name)

	if existing_pt != nil {
		return fmt.Errorf("Placetype name '%s' (%d) already registered", pt.Name, pt.Id)
	}

	/*
		for _, pid := range pt.Parent {

			_, err := spec.GetPlacetypeById(pid)

			if err != nil {
				return fmt.Errorf("Failed to get placetype by ID %d, %w", pid, err)
			}
		}
	*/

	str_id := strconv.FormatInt(pt.Id, 10)
	spec.catalog[str_id] = pt
	return nil
}

func (spec *WOFPlacetypeSpecification) Catalog() map[string]WOFPlacetype {
	return spec.catalog
}

// IsValidPlacetypeId returns a boolean value indicating whether 'name' is a known and valid placetype name.
func (spec *WOFPlacetypeSpecification) IsValidPlacetype(name string) bool {

	for _, pt := range spec.catalog {

		if pt.Name == name {
			return true
		}
	}

	return false
}

// IsValidPlacetypeId returns a boolean value indicating whether 'id' is a known and valid placetype ID.
func (spec *WOFPlacetypeSpecification) IsValidPlacetypeId(id int64) bool {

	for str_id, _ := range spec.catalog {

		pt_id, err := strconv.Atoi(str_id)

		if err != nil {
			continue
		}

		pt_id64 := int64(pt_id)

		if pt_id64 == id {
			return true
		}
	}

	return false
}

// Returns true is 'b' is an ancestor of 'a'.
func (spec *WOFPlacetypeSpecification) IsAncestor(a *WOFPlacetype, b *WOFPlacetype) bool {

	roles := []string{
		COMMON_ROLE,
		OPTIONAL_ROLE,
		COMMON_OPTIONAL_ROLE,
		CUSTOM_ROLE,
	}

	str_roles := strings.Join(roles, "-")
	key := fmt.Sprintf("%d_%d_%s_is_ancestor", a.Id, b.Id, str_roles)

	if !spec.isIndexingRelationships() {

		v, ok := spec.relationships.Load(key)

		if ok {
			return v.(bool)
		}
	}

	is_ancestor := false

	for _, ancestor := range spec.AncestorsForRoles(a, roles) {

		if ancestor.Name == b.Name {
			is_ancestor = true
			break
		}
	}

	spec.relationships.Store(key, is_ancestor)
	return is_ancestor
}

// Returns true is 'b' is a descendant of 'a'.
func (spec *WOFPlacetypeSpecification) IsDescendant(a *WOFPlacetype, b *WOFPlacetype) bool {

	roles := []string{
		COMMON_ROLE,
		OPTIONAL_ROLE,
		COMMON_OPTIONAL_ROLE,
		CUSTOM_ROLE,
	}

	str_roles := strings.Join(roles, "-")
	key := fmt.Sprintf("%d_%d_%s_is_descendant", a.Id, b.Id, str_roles)

	v, ok := spec.relationships.Load(key)

	if !spec.isIndexingRelationships() {

		if ok {
			return v.(bool)
		}
	}

	is_descendant := false

	for _, descendant := range spec.DescendantsForRoles(a, roles) {

		if descendant.Name == b.Name {
			is_descendant = true
			break
		}
	}

	spec.relationships.Store(key, is_descendant)
	return is_descendant
}

// Children returns the immediate child placetype of 'pt'.
func (spec *WOFPlacetypeSpecification) Children(pt *WOFPlacetype) []*WOFPlacetype {

	key := fmt.Sprintf("%d_children", pt.Id)

	if !spec.isIndexingRelationships() {

		v, ok := spec.relationships.Load(key)

		if ok {
			return v.([]*WOFPlacetype)
		}
	}

	children := make([]*WOFPlacetype, 0)

	for _, details := range spec.catalog {

		for _, pid := range details.Parent {

			if pid == pt.Id {
				child_pt, _ := spec.GetPlacetypeByName(details.Name)
				children = append(children, child_pt)
			}
		}
	}

	sorted := sortChildren(pt, children)

	spec.relationships.Store(key, sorted)
	return sorted
}

// Ancestors returns the ancestors of role "common" for 'pt'.
func (spec *WOFPlacetypeSpecification) Ancestors(pt *WOFPlacetype) []*WOFPlacetype {
	return spec.AncestorsForRoles(pt, []string{"common"})
}

// AncestorsForRoles returns the ancestors matching any role in 'roles' for 'pt'.
func (spec *WOFPlacetypeSpecification) AncestorsForRoles(pt *WOFPlacetype, roles []string) []*WOFPlacetype {

	str_roles := strings.Join(roles, "-")
	key := fmt.Sprintf("%d_ancestors_%s", pt.Id, str_roles)

	if !spec.isIndexingRelationships() {

		v, ok := spec.relationships.Load(key)

		if ok {
			return v.([]*WOFPlacetype)
		}
	}

	ancestors := make([]*WOFPlacetype, 0)
	ancestors = spec.fetchAncestors(pt, roles, ancestors)

	spec.relationships.Store(key, ancestors)
	return ancestors
}

// Descendants returns the descendants of role "common" for 'pt'.
func (spec *WOFPlacetypeSpecification) Descendants(pt *WOFPlacetype) []*WOFPlacetype {
	return spec.DescendantsForRoles(pt, []string{"common"})
}

// DescendantsForRoles returns the descendants matching any role in 'roles' for 'pt'.
func (spec *WOFPlacetypeSpecification) DescendantsForRoles(pt *WOFPlacetype, roles []string) []*WOFPlacetype {

	str_roles := strings.Join(roles, "-")
	key := fmt.Sprintf("%d_descendants_%s", pt.Id, str_roles)

	if !spec.isIndexingRelationships() {

		v, ok := spec.relationships.Load(key)

		if ok {
			return v.([]*WOFPlacetype)
		}
	}

	descendants := make([]*WOFPlacetype, 0)
	descendants = spec.fetchDescendants(pt, roles, descendants)

	spec.relationships.Store(key, descendants)
	return descendants
}

func (spec *WOFPlacetypeSpecification) GraphPlacetypes() (graph.Graph, error) {

	placetypeHash := func(pt *WOFPlacetype) string {
		return pt.Name
	}

	gr := graph.New(placetypeHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())

	lookup := new(sync.Map)

	for str_id, pt := range spec.catalog {

		err := gr.AddVertex(&pt)

		if err != nil {
			return nil, fmt.Errorf("Failed to add vertex for %v, %w", pt, err)
		}

		p_id := pt.Id

		if p_id == 0 {

			p_id, err = strconv.ParseInt(str_id, 10, 64)

			if err != nil {
				return nil, fmt.Errorf("Failed to parse string ID for %v (%s), %w", pt, str_id, err)
			}
		}

		lookup.Store(p_id, pt)
	}

	for _, pt := range spec.catalog {

		for _, p_id := range pt.Parent {

			v, exists := lookup.Load(p_id)

			if !exists {
				return nil, fmt.Errorf("Failed to load parent ID (%d) for %v", p_id, pt)
			}

			p_pt := v.(WOFPlacetype)

			err := gr.AddEdge(p_pt.Name, pt.Name)

			if err != nil {
				return nil, fmt.Errorf("Failed to add edge between %v and %v, %w", p_pt, pt, err)
			}
		}

	}

	return gr, nil
}

func (spec *WOFPlacetypeSpecification) fetchDescendants(pt *WOFPlacetype, roles []string, descendants []*WOFPlacetype) []*WOFPlacetype {

	grandkids := make([]*WOFPlacetype, 0)

	for _, kid := range spec.Children(pt) {

		descendants = spec.appendPlacetype(kid, roles, descendants)

		for _, grandkid := range spec.Children(kid) {
			grandkids = spec.appendPlacetype(grandkid, roles, grandkids)
		}
	}

	for _, k := range grandkids {
		descendants = spec.appendPlacetype(k, roles, descendants)
		descendants = spec.fetchDescendants(k, roles, descendants)
	}

	// fmt.Println(pt, descendants)
	return descendants
}

func (spec *WOFPlacetypeSpecification) fetchAncestors(pt *WOFPlacetype, roles []string, ancestors []*WOFPlacetype) []*WOFPlacetype {

	for _, id := range pt.Parent {

		parent, _ := spec.GetPlacetypeById(id)

		role_ok := false

		for _, r := range roles {

			if r == parent.Role {
				role_ok = true
				break
			}
		}

		if !role_ok {
			continue
		}

		append_ok := true

		for _, a := range ancestors {

			if a.Id == parent.Id {
				append_ok = false
				break
			}
		}

		if append_ok {

			has_grandparent := false
			offset := -1

			for _, gpid := range parent.Parent {

				for idx, a := range ancestors {

					if a.Id == gpid {
						offset = idx
						has_grandparent = true
						break
					}
				}

				if has_grandparent {
					break
				}
			}

			// log.Printf("APPEND %s < %s GP: %t (%d)\n", parent.Name, pt.Name, has_grandparent, offset)

			if has_grandparent {

				// log.Println("WTF 1", len(ancestors))

				tail := ancestors[offset+1:]
				ancestors = ancestors[0:offset]

				ancestors = append(ancestors, parent)

				for _, a := range tail {
					ancestors = append(ancestors, a)
				}

			} else {
				ancestors = append(ancestors, parent)
			}
		}

		ancestors = spec.fetchAncestors(parent, roles, ancestors)
	}

	return ancestors
}

func (spec *WOFPlacetypeSpecification) appendPlacetype(pt *WOFPlacetype, roles []string, others []*WOFPlacetype) []*WOFPlacetype {

	do_append := true

	for _, o := range others {

		if pt.Id == o.Id {
			do_append = false
			break
		}
	}

	if !do_append {
		return others
	}

	has_role := false

	for _, r := range roles {

		if pt.Role == r {
			has_role = true
			break
		}
	}

	if !has_role {
		return others
	}

	others = append(others, pt)
	return others
}

func sortChildren(pt *WOFPlacetype, all []*WOFPlacetype) []*WOFPlacetype {

	kids := make([]*WOFPlacetype, 0)
	grandkids := make([]*WOFPlacetype, 0)

	for _, other := range all {

		is_grandkid := false

		for _, pid := range other.Parent {

			for _, p := range all {

				if pid == p.Id {
					is_grandkid = true
					break
				}
			}

			if is_grandkid {
				break
			}
		}

		if is_grandkid {
			grandkids = append(grandkids, other)
		} else {
			kids = append(kids, other)
		}
	}

	if len(grandkids) > 0 {
		grandkids = sortChildren(pt, grandkids)
	}

	for _, k := range grandkids {
		kids = append(kids, k)
	}

	return kids
}
