package permission

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNewBox(t *testing.T) {
	var box = NewBox()
	id, err := box.Registe("func1")

	assert.Equal(t, 0, id)
	assert.Nil(t, err)

	id, err = box.Registe("func2")
	assert.Equal(t, 1, id)
	assert.Nil(t, err)

	// test size
	assert.Equal(t, 2, box.Size())

	// test remove
	err = box.Remove("funcn")
	assert.NotNil(t, err)
	err = box.Remove("func2")
	assert.Nil(t, err)
	assert.Equal(t, 1, box.Size())

	//test exist
	ret := box.Exist(-1, "funcn")
	assert.Equal(t, false, ret)

	ret = box.Exist(1, "")
	assert.Equal(t, false, ret)

	// register agian after remove
	id, err = box.Registe("func2")

	assert.Equal(t, 2, box.Size())
	assert.Equal(t, 1, id)

	// add exist func
	id, err = box.Registe("func1")
	assert.Equal(t, 0, id)
	assert.NotNil(t, err)

	name, ok := box.GetName(0)
	assert.Equal(t, true, ok)
	assert.Equal(t, "func1", name)

	id, ok = box.GetId("func2")
	assert.Equal(t, 1, id)
	assert.True(t, ok)

	//GetScope
	scope := box.GetScope(1)
	assert.Equal(t, SCOPE_USED, scope)

	//GetScopeByName
	scope = box.GetScopeByName("func1")
	assert.Equal(t, SCOPE_USED, scope)
	scope = box.GetScopeByName("funcn")
	assert.Equal(t, SCOPE_EMPTY, scope)

	//Update
	err = box.Update(1, SCOPE_UNUSED)
	assert.Nil(t, err)
	assert.Equal(t, SCOPE_UNUSED, box.GetScope(1))

	err = box.Update(2, SCOPE_USED)
	assert.NotNil(t, err)
	assert.Equal(t, SCOPE_UNUSED, box.GetScope(1))

	err = box.Update("funcn", SCOPE_USED)
	assert.NotNil(t, err)
	assert.Equal(t, 2, box.Size())

	err = box.Update("func2", SCOPE_USED)
	assert.Nil(t, err)
	assert.Equal(t, SCOPE_USED, box.GetScopeByName("func2"))

	err = box.Update(int32(1), SCOPE_UNUSED)
	assert.NotNil(t, err)
	assert.Equal(t, SCOPE_USED, box.GetScope(1))
}

func TestBigger32Func(t *testing.T) {
	var box = NewBox()
	for i := 0; i < 35; i++ {
		name := "func" + strconv.Itoa(i)
		box.Registe(name)
	}
	assert.Equal(t, 35, box.Size())
}

func TestCommon(t *testing.T) {
	test_println("test debug print function")
}
