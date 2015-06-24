package jsonrpc

// JSON-RPC 2.0 message types. Version field is ommited for brevity.

// Request is a JSON-RPC request object.
type Request struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

// Notification is a JSON-RPC notification object.
type Notification struct {
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

// Response is a JSON-RPC response object.
type Response struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  *ResError   `json:"error,omitempty"`
}

// ResError is a JSON-RPC response error object.
type ResError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// message is a JSON-RPC request, response, or notification message.
type message struct {
	ID     string      `json:"id,omitempty"`
	Method string      `json:"method,omitempty"`
	Params interface{} `json:"params,omitempty"`
	Result interface{} `json:"result,omitempty"`
	Error  *ResError   `json:"error,omitempty"`
}
