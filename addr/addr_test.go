package addr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddr(t *testing.T) {
	ad := NewAddr("127.0.0.1", "2324", "temp")
	assert.Equal(t, "127.0.0.1", ad.Host())
	assert.Equal(t, "2324", ad.Port())
	assert.Equal(t, "temp", ad.Network())
	assert.Equal(t, "127.0.0.1:2324", ad.String())
	assert.Equal(t, 2324, ad.PortInt())

	ad = NewAddr("127.0.0.1", "2324")
	assert.Equal(t, "127.0.0.1", ad.Host())
	assert.Equal(t, "2324", ad.Port())
	assert.Equal(t, "", ad.Network())

	ad = NewAddr("127.0.0.1:80")
	assert.Equal(t, "127.0.0.1", ad.Host())
	assert.Equal(t, "80", ad.Port())

	assert.Empty(t, ad.MyIP())

	ad = NewAddr("127.0.0.1")
	assert.Equal(t, 0, ad.PortInt())
}
