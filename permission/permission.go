package permission

// default permission

var (
	invoke_box = NewBox()
	access_box = NewBox()
)

func NewInvoke() *ChainInvoke {
	return NewChainInvoke(invoke_box)
}

type IAccess interface {
	GetAccess() Accesser
	SetAccess(Accesser)
}

// apply an accesser.
func NewAccess(name string) (Accesser, error) {
	id, err := access_box.Registe(name)
	if err != nil {
		return nil, err
	}
	ca := NewChainAccess(invoke_box, access_box)
	ca.access_id = id
	return ca, nil
}

// directly regist invoke handle by name
func RegisteInvoke(name string) (int, error) {
	return invoke_box.Registe(name)
}
