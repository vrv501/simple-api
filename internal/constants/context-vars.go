package constants

type CtxKey int

const (
	_ CtxKey = iota

	RequestIDKey

	ClientIPAddr

	HTTPMethod

	URLPath
)
