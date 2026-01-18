package logg

type Checker interface {
	CheckWithFatal(error, string)
	CheckWithError(error, string)
}

type Closer interface {
	Close()
}

type Raiser interface {
	RaiseInfo(string, Fields)
	RaiseWarn(string, Fields)
	RaiseError(error, string, Fields)
	RaiseFatal(error, string, Fields)
	RaisePanic(error, string, Fields)
}

type Logger interface {
	Checker
	Raiser
	Closer
}
