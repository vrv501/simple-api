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
	LogFieldRequestID    = "request_id"
	LogFieldClientIP     = "client_ip"
	LogFieldMethodAndURL = "url"
	LogFieldStatus       = "status"
	LogFieldStackTrace   = "stack_trace"
	LogFieldPanic        = "panic"
	LogFieldLatency      = "latency"
)
