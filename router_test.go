package gungnir

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteMatch(t *testing.T) {
	paths := []string{"aa", "bb", "cc"}
	r := newRoute(paths[0])
	g := r.Group(paths[1])
	g.GET(paths[2], func(ctx *Ctx) {})
	g.GET(filepath.Join(paths[2], paths[3]), nil)

	matchedRoute, _, cnt := r.matchRoute(paths)
	assert.EqualValues(t, paths[1], matchedRoute.path.sections[0].path)
	assert.EqualValues(t, 2, cnt)

	paths = []string{"aa", "bb", "{id}"}
	newRoute(paths[0]).Group(filepath.Join(paths[1], paths[2]))
	matchedRoute, _, _ = r.matchRoute(paths)
	assert.EqualValues(t, paths[2], matchedRoute.path.sections[1].path)
}
