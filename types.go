package onewallet

// Request is the structure of the onewallet messaging
type Request struct {
	Action string      `json:"type"`
	Data   interface{} `json:"data"`
}
