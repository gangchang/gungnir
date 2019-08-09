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
	*node
}

func newTree(path string) *tree {
	paths := getPaths(path)
	n := &node{paths: paths}
	t := &tree{node: n}
	return t
}

type node struct {
	paths []string

	children      map[string]*node
	wildcardChild *node

	// exec before handler, if return false, just return
	beforeExecFns []beforeExecFn
	// exec after handler, if return false, just return
	afterExecFns []afterExecFn

	pathParams map[int]string
	// key is method
	handlers map[string]handler
}

func (n *node) Register(api interface{}) {
	ph, ok := api.(typePath)
	if !ok {
		// TODO
		panic("must support type path method")
	}
	typeName := ph.Path()

	// create
	coh, ok := api.(createOneHandler)
	if ok {
		n.addHandlerToChildren(typeName, http.MethodPost, coh.CreateOne)
	}
	cmh, ok := api.(createManyHandler)
	if ok {
		n.addHandlerToChildren(getCollectionPath(typeName), http.MethodPost, cmh.CreateMany)
	}

	// read
	roh, ok := api.(readOneHandler)
	if ok {
		n.addHandlerToChildren(GetIDPath(typeName), http.MethodGet, roh.ReadOne)
	}
	rmh, ok := api.(readManyHandler)
	if ok {
		n.addHandlerToChildren(getCollectionPath(typeName), http.MethodGet, rmh.ReadMany)
	}

	// update
	uoh, ok := api.(updateOneHandler)
	if ok {
		n.addHandlerToChildren(GetIDPath(typeName), http.MethodPatch, uoh.UpdateOne)
	}
	umh, ok := api.(updateManyHandler)
	if ok {
		n.addHandlerToChildren(getCollectionPath(typeName), http.MethodPatch, umh.UpdateMany)
	}

	// delete
	doh, ok := api.(deleteOneHandler)
	if ok {
		n.addHandlerToChildren(GetIDPath(typeName), http.MethodDelete, doh.DeleteOne)
	}
	dmh, ok := api.(deleteManyHandler)
	if ok {
		n.addHandlerToChildren(getCollectionPath(typeName), http.MethodDelete, dmh.DeleteMany)
	}
}

func (n *node) AddBeforeExecFns(befns ...beforeExecFn) {
	for _, v := range befns {
		n.beforeExecFns = append(n.beforeExecFns, v)
	}
}

func (n *node) AddAfterExecFns(aefns ...afterExecFn) {
	for _, v := range aefns {
		n.afterExecFns = append(n.afterExecFns, v)
	}
}

func (n *node) addHandlerToChildren(path, method string, h handler) *node {
	subNode, exists := n.children[path]
	if !exists {
		subNode = n.addChildren(path)
	}
	if subNode.handlers == nil {
		subNode.handlers = make(map[string]handler)
	}
	if _, exists := subNode.handlers[method]; !exists {
		subNode.handlers[method] = h
	}

	return subNode
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
	cn.beforeExecFns = append(cn.beforeExecFns, n.beforeExecFns...)
	cn.afterExecFns = append(cn.afterExecFns, n.afterExecFns...)
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
		if n.wildcardChild != nil {
			// TODO
			panic("confilct")
		}
		n.wildcardChild = cn
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

func (t *tree) matchHandler(path, method string) (*node, handler, pathParams) {
	paths := getPaths(path)

	cnt := 0
	pp := newPathParams()

	nowNode := t.node
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
		if nowNode.wildcardChild != nil {
			cnt += len(nowNode.paths)
			nowNode = nowNode.wildcardChild
		} else {
			// NOT found
			return nil, nil, nil
		}
	}

handlerLabel:
	h := nowNode.matchHandler(method)
	return nowNode, h, pp
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
