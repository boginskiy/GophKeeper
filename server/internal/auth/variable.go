package auth

var (
	// Keys for context for request.
	EmailCtx = keycontext{}
	PhoneCtx = keycontext{}
)

// keycontext - is type of key for values for context request.
type keycontext struct{}
