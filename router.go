package gungnir

import (
	"net/http"
	"strings"
)

type pathParams map[string]string

func newPathParams() pathParams {
	return make(map[string]string)
}

func (pp pathParams) add(subPathParams map[string]string) {
	for k, v := range subPathParams {
		pp[k] = v
	}
}

type handler func(Ctx)

type tree struct {
	rootNode *node
}

func newTree(path string) *tree {
	paths := getPaths(path)
	n := &node{paths: paths}
	t := &tree{rootNode: n}
	return t
}

type node struct {
	paths                []string
	children             map[string]*node
	wildcardChildrenNode *node
	pathParams           map[int]string
	// key is method
	handlers map[string]handler
}

func newNode(path string) *node {
	paths := getPaths(path)
	pp := getPathParams(paths)
	return &node{
		paths:      paths,
		pathParams: pp,
	}
}

func (n *node) addChildren(path string) *node {
	cn := newNode(path)
	// not wildcard
	if len(cn.pathParams) == 0 || !strings.HasPrefix(cn.paths[0], ":") {
		if n.children == nil {
			n.children = make(map[string]*node)
		}
		_, exists := n.children[cn.paths[0]]
		if exists {
			// TODO
			panic("confilt")
		}
		n.children[cn.paths[0]] = cn
	} else {
		if n.wildcardChildrenNode != nil {
			// TODO
			panic("confilct")
		}
		n.wildcardChildrenNode = cn
	}

	return cn
}

// path length should > 0
func (n *node) Group(path string) *node {
	return n.addChildren(path)
}

func (n *node) GET(h handler) {
	n.addHandler(http.MethodGet, h)
}

func (n *node) POST(h handler) {
	n.addHandler(http.MethodPost, h)
}

func (n *node) addHandler(method string, h handler) {
	if n.handlers == nil {
		n.handlers = make(map[string]handler)
	}
	if _, exists := n.handlers[method]; exists {
		panic("this has exists")
	}

	n.handlers[method] = h
}

func (t *tree) matchHandler(path, method string) (handler, pathParams) {
	paths := getPaths(path)

	cnt := 0
	pp := newPathParams()

	nowNode := t.rootNode
	for {
	routeLabel:
		// first, itself check
		subPp, matched := nowNode.matched(paths[cnt:])
		if matched {
			pp.add(subPp)
			goto handlerLabel
		}
		// second, children chcek
		for k, v := range nowNode.children {
			if k == paths[cnt] {
				cnt += len(nowNode.paths) - 1
				nowNode = v
				goto routeLabel
			}
		}
		// last, wildcard children check
		if nowNode.wildcardChildrenNode != nil {
			cnt += len(nowNode.paths)
			nowNode = nowNode.wildcardChildrenNode
		} else {
			// NOT found
			return nil, nil
		}
	}

handlerLabel:
	h := nowNode.matchHandler(method)
	return h, pp
}

func (n *node) matchHandler(method string) handler {
	for m, h := range n.handlers {
		if m == method {
			return h
		}
	}

	return nil
}

func (n *node) matched(paths []string) (map[string]string, bool) {
	if len(n.paths) != len(paths) {
		return nil, false
	}

	pathParams := make(map[string]string)
	for k, v := range paths {
		if pathParam, exists := n.pathParams[k]; exists {
			pathParams[pathParam] = paths[k]
		} else {
			if v != n.paths[k] {
				return nil, false
			}
		}
	}

	return pathParams, true
}
