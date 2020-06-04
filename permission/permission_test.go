package permission

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	var inv = NewInvoke()
	id, err := inv.AutoSet("func1", SCOPE_USED)
	assert.Equal(t, SCOPE_USED, inv.GetScopeById(id))
	assert.Nil(t, err)

	id2, err := RegisteInvoke("func3")
	assert.Equal(t, 2, id2)
	assert.Nil(t, err)
	fmt.Println(err)
}
