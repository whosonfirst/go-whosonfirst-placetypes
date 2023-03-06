package placetypes

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type Foo interface {
	Specification() *WOFPlacetypeSpecification
	Property() string
}

var foo_roster roster.Roster

// FooInitializationFunc is a function defined by individual foo package and used to create
// an instance of that foo
type FooInitializationFunc func(ctx context.Context, uri string) (Foo, error)

// RegisterFoo registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Foo` instances by the `NewFoo` method.
func RegisterFoo(ctx context.Context, scheme string, init_func FooInitializationFunc) error {

	err := ensureFooRoster()

	if err != nil {
		return err
	}

	return foo_roster.Register(ctx, scheme, init_func)
}

func ensureFooRoster() error {

	if foo_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		foo_roster = r
	}

	return nil
}

// NewFoo returns a new `Foo` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `FooInitializationFunc`
// function used to instantiate the new `Foo`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterFoo` method.
func NewFoo(ctx context.Context, uri string) (Foo, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := foo_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(FooInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureFooRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range foo_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
