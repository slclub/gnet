package encoder

import (
	"encoding/json"
	//"fmt"
	"github.com/slclub/link"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoding(t *testing.T) {
	coding := Coding("")
	assert.NotNil(t, coding)

	data_map := map[string]interface{}{
		"had": "ndec",
		"key": 2321,
	}

	s1 := coding.Encode(data_map)

	link.DEBUG_PRINT("[json encoding]", s1)

	delete(data_map, "had")
	delete(data_map, "key")
	coding.Decode(s1, &data_map)
	assert.Equal(t, 2, len(data_map))

	ctype := coding.ContentType()
	assert.Equal(t, "application/json", ctype)
}

func TestCodingError(t *testing.T) {
	coding := Coding("json")

	s1 := coding.Encode(nil)
	link.DEBUG_PRINT("[CODING][ASSERT_ERROR]", s1)

	data_map := ""
	s2 := `{"had":"ndec","key":2321}`
	coding.Decode(s2, &data_map)

	// not get coding
	coding = Coding("xml")
	assert.Nil(t, coding)
}

func TestEncodingJson(t *testing.T) {
	coding := Coding("")
	data_map := map[string]interface{}{
		"had": "ndec",
		"key": 2321,
	}

	s1 := coding.Encode(data_map)

	delete(data_map, "had")
	delete(data_map, "key")

	_ = json.Unmarshal([]byte(s1), &data_map)

	//link.DEBUG_PRINT("[EncodingJson][json encoding]", err)
	assert.Equal(t, 2, len(data_map))
}
