package defined

const (
	// success code
	SUCCESS = 0
	// jwt login
	ERR_AUTH_GENERATE              = -200001
	ERR_AUTH                       = -200002
	ERROR_AUTH_CHECK_TOKEN_FAIL    = -200003
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = -200004

	// jump current process node.
	CODE_JUMP_CURRENT_NODE = -91001
	// jump current goroutine.
	CODE_JUMP_CURRENT_FLOW = -91002

	//gcore.Middle
	// empty handle used.
	CODE_NOT_ALLOWED_EMPTY_HANDLE = -91010
)
