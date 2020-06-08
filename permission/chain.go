package permission

import (
	"github.com/slclub/gerror"
)

// ==========================================basic chain==================================
// ---------------------------------interface -----------------------------------
type ChainSetter interface {
	Set(id int, scope uint8) (err error)
}

type ChainGetter interface {
	GetChainScope() map[int]uint8
	GetScopeById(id int) uint8
	GetInvoke() BoxInvoker
	GetId(name string) (int, bool)
}

// ---------------------------------struct -----------------------------------
type Chain struct {
	box         BoxInvoker
	keys_handle map[int]uint8
}

func (c *Chain) GetChainScope() map[int]uint8 {
	return c.keys_handle
}

func (c *Chain) GetScopeById(idx int) uint8 {
	return c.keys_handle[idx]
}

func (c *Chain) GetInvoke() BoxInvoker {
	return c.box
}

func (c *Chain) Set(id int, scope uint8) (err error) {
	if c.box.Exist(id, "") == false {
		return gerror.New("-90004:[PERMISSION][ACCESS][SET]ID[", id, "]")
	}
	c.keys_handle[id] = scope
	return
}

func (c *Chain) GetId(name string) (int, bool) {
	return c.box.GetId(name)
}

// ==========================================basic chain==================================
// ==========================================invoke chain==================================

// ---------------------------------interface -----------------------------------
type InvokeSetter interface {
	ChainSetter
	AutoSet(name string, scope uint8) (int, error)
	Remove(id int) (err error)
}

type Invoker interface {
	InvokeSetter
	ChainGetter

	// valid interface
	Validate(id int, access AccessGetter) bool
	ValidateByName(name string, access AccessGetter) bool
}

var _ Invoker = &ChainInvoke{}

// ---------------------------------struct -----------------------------------
type ChainInvoke struct {
	Chain
}

func NewChainInvoke(box BoxInvoker) *ChainInvoke {
	ci := &ChainInvoke{}

	// this line can't write in defined. will create complie error.
	// Can't use promoted field Chain.box in struct literal of type ChainInvoke
	ci.box = box
	ci.keys_handle = make(map[int]uint8)
	return ci
}

func (ci *ChainInvoke) AutoSet(name string, scope uint8) (int, error) {
	if ci.box.Exist(-1, name) {
		id, _ := ci.box.GetId(name)
		if scope == SCOPE_EMPTY {
			ci.Set(id, ci.box.GetScope(id))
		} else {
			ci.Set(id, scope)
		}
		return id, nil
	}
	// default scope == SCOPE_USED
	id, err := ci.box.Registe(name)
	if err != nil {
		return 0, err
	}
	ci.Set(id, ci.box.GetScope(id))
	if scope != SCOPE_EMPTY {
		ci.Set(id, scope)
		//ci.box.Update(id, scope)
	}
	return id, nil
}

func (ci *ChainInvoke) Remove(id int) (err error) {
	if !ci.box.Exist(id, "") {
		return gerror.New("-90011:[PERMISSION][CHAIN][INVOKE][REMOVE]ID[", id, "]")
	}
	delete(ci.keys_handle, id)
	return
}

// handle execution validate.
func (ci *ChainInvoke) Validate(id int, access AccessGetter) bool {
	scope := access.GetScopeById(id)
	switch scope {
	case SCOPE_USED:
		return true
	case SCOPE_UNUSED:
		return false
	}

	// There is no id in access object.
	// Used the chain default scope.
	scope = ci.GetScopeById(id)
	if scope == SCOPE_USED {
		return true
	}
	return false
}

func (ci *ChainInvoke) ValidateByName(name string, access AccessGetter) bool {
	id, ok := ci.GetId(name)
	if !ok {
		return false
	}
	return ci.Validate(id, access)
}

// ==========================================invoke chain==================================
// ==========================================Access chain==================================

// ---------------------------------interface -----------------------------------
type AccessGetter interface {
	ChainGetter
	GetScope() uint8
	GetAID() int
}

type AccessSetter interface {
	ChainSetter
}

type Accesser interface {
	AccessGetter
	AccessSetter
}

var _ Accesser = &ChainAccess{}

// ---------------------------------struct -----------------------------------
type ChainAccess struct {
	Chain
	access    BoxAccesser
	access_id int
}

func NewChainAccess(inv BoxInvoker, access BoxAccesser) *ChainAccess {
	ci := &ChainAccess{}
	ci.box = inv
	ci.access = access
	ci.keys_handle = make(map[int]uint8)
	return ci
}

func (ca *ChainAccess) GetAID() int {
	return ca.access_id
}

func (ca *ChainAccess) GetScope() uint8 {
	return ca.access.GetScope(ca.access_id)
}

// ==========================================Access chain==================================
