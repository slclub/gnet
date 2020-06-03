package permission

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewChain(t *testing.T) {
	var invoke = NewInvoke()
	invoke.AutoSet("func1", 0)
	assert.Equal(t, SCOPE_USED, invoke.GetScopeById(0))

	invoke.GetChainScope()
	inv_box := invoke.GetInvoke()
	assert.Equal(t, 1, inv_box.Size())

	// Set not exist func

	err := invoke.Set(100, SCOPE_USED)
	assert.NotNil(t, err)

	id, err := invoke.AutoSet("func2", SCOPE_UNUSED)
	assert.Equal(t, 1, id)

	//test remove
	err = invoke.Remove(1)
	assert.Nil(t, err)
	assert.Equal(t, 2, invoke.GetInvoke().Size())
	assert.Equal(t, 1, len(invoke.GetChainScope()))
}

func TestAnFlow(t *testing.T) {
	var invoke_chain = NewInvoke()
	id, _ := invoke_chain.AutoSet("func1", SCOPE_UNUSED)

	var access_chain, ok = NewAccess("access_node1")
	assert.Nil(t, ok)

	ret := invoke_chain.Validate(id, access_chain)
	assert.Equal(t, SCOPE_UNUSED, invoke_chain.GetScopeById(id))
	assert.False(t, ret)
	assert.Equal(t, SCOPE_USED, invoke_chain.GetInvoke().GetScope(id))

	access_chain.Set(id, SCOPE_USED)
	ret = invoke_chain.Validate(id, access_chain)
	assert.True(t, ret)

	id2, _ := invoke_chain.AutoSet("func2", 0)
	ret = invoke_chain.Validate(id2, access_chain)
	assert.True(t, ret)
}
