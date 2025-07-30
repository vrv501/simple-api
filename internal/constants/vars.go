package constants

// ENV Vars
const (
	LogLevel   = "LOG_LEVEL"
	ServerPort = "SERVER_PORT"
)

// Headers
const (
	HeaderRequestID = "X-Request-ID"
)

// Logging Vars
const (
	LogFieldRequestID  = "request_id"
	LogFieldClientIP   = "client_ip"
	LogFieldHTTPMethod = "http_method"
	LogFieldURLPath    = "url_path"
	LogFieldStackTrace = "stack_trace"
	LogFieldPanic      = "panic"
)

// Context Keys
type CtxKey int

const (
	_ CtxKey = iota
	RequestIDKey
	ClientIPAddr
	HTTPMethod
	URLPath
)
