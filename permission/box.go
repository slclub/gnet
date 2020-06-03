package permission

import (
	"fmt"
	"github.com/slclub/gerror"
	"strconv"
)

const (
	SCOPE_EMPTY  = uint8(0)
	SCOPE_UNUSED = uint8(1)
	SCOPE_USED   = uint8(2)
)

type Boxer interface {
	Registe(name string) (id int, err error)
	GetId(name string) (id int, ok bool)
	GetScopeByName(name string) (sp uint8)
	GetScope(id int) (sp uint8)
	Exist(id int, name string) bool
	Update(key interface{}, value uint8) error
	Size() int
}

type BoxerRemove interface {
	Remove(name string) (err error)
}

// All every part of function module optional or defined by user.
// Will  registe here and return an id.
// It is baseic class for invoke and access.
type Box struct {
	//
	scope    []uint8
	name_key map[string]int
	key_name map[int]string
}

func NewBox() *Box {
	return &Box{
		scope:    make([]uint8, 32),
		name_key: make(map[string]int),
		key_name: make(map[int]string),
	}
}

func (bx *Box) Registe(name string) (id int, err error) {
	id = bx.GetNextID()
	err = bx.Add(id, name)
	if err != nil {
		id = 0
		return
	}
	return
}

func (bx *Box) Size() int {
	return len(bx.name_key)
}

func (bx *Box) Remove(name string) (err error) {
	id, ok := bx.GetId(name)
	if !ok {
		return gerror.New("[PERMISSION][BOX][NAME NOT EXIST]", name)
	}
	delete(bx.key_name, id)
	delete(bx.name_key, name)

	bx.scope[id] = SCOPE_EMPTY
	return
}

// Get next id.
func (bx *Box) GetNextID() (id int) {
	// First get id from empty scope item.
	for i, v := range bx.scope {
		//test_println("next value", v)
		if v == SCOPE_EMPTY {
			id = i
			return id
		}
	}
	id = len(bx.scope)
	return
}

func (bx *Box) Check(id int, name string) (err error) {
	if _, ok := bx.key_name[id]; ok && id >= 0 {
		return gerror.New("-90001:[PERMISSION][ID EXIST]", strconv.Itoa(id))
	}
	if _, ok := bx.name_key[name]; ok && name != "" {
		return gerror.New("-90002:[PERMISSION][MODULE NAME EXIST]", name)
	}
	return
}

func (bx *Box) Exist(id int, name string) bool {
	// here I am lazy and write less code
	err := bx.Check(id, name)
	if err != nil {
		return true
	}
	return false
}

func (bx *Box) Add(id int, name string) (err error) {
	if ok := bx.Check(id, name); ok != nil {
		return ok
	}
	if id == len(bx.scope) {
		bx.scope = append(bx.scope, SCOPE_USED)
	} else {
		bx.scope[id] = SCOPE_USED
	}

	bx.name_key[name] = id
	bx.key_name[id] = name
	return
}

func (bx *Box) GetName(id int) (name string, ok bool) {
	name, ok = bx.key_name[id]
	return
}

func (bx *Box) GetId(name string) (id int, ok bool) {
	id, ok = bx.name_key[name]
	return
}

func (bx *Box) GetScope(id int) (sp uint8) {
	sp = bx.scope[id]
	return
}

// Will return zero value if can't find the name
func (bx *Box) GetScopeByName(name string) (sp uint8) {
	id, ok := bx.GetId(name)
	if ok == false {
		return SCOPE_EMPTY
	}
	sp = bx.scope[id]
	return
}

func (bx *Box) Update(key interface{}, value uint8) error {

	if id, ok := key.(int); ok {
		if !bx.Exist(id, "") {
			return gerror.New("-90003:[PERMISSION][BOX][UPDATE][NOT FOUND]ID[", id, "]")
		}
		bx.scope[id] = value
		return nil
	}

	if name, ok := key.(string); ok {
		if !bx.Exist(-1, name) {
			return gerror.New("-90003:[PERMISSION][BOX][UPDATE][NOT FOUND]NAME[", name, "]")
		}
		id, ret := bx.GetId(name)
		if !ret {
			return gerror.New("-90004:[PERMISSION][BOX][UPDATE][FOUND NAME][NOT FOUND ID]NAME[", name, "]")
		}
		bx.scope[id] = value
		return nil
	}

	return gerror.New("-90005:[PERMISSION][BOX][UPDATE][NOT SUPPORT OTHER TYPE]KEY ONLY[int, string]")
}

var _ Boxer = &Box{}

// ===========================================================================
// Here we split invoker and accesser in different type box.
// ===========================================================================

// where invoke function will regist.
type BoxInvoker interface {
	Boxer
}

type BoxAccesser interface {
	Boxer
}

func test_println(args ...interface{}) {
	fmt.Println("[BOX]", args)
}
