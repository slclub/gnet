package encoder

type Encoder interface {
	ContentType() string

	Encode(data interface{}) string
	EncodeBytes(data interface{}) []byte

	Decode(data string, obj interface{})
}

func Coding(selecting string) Encoder {
	switch selecting {
	case "json", "", "application/json":
		return &Json{}
	}
	return nil
}
