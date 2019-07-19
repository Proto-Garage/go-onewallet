package onewallet


type Request struct {
	Action string      `json:"type"`
	Data   interface{} `json:"data"`
}

type Arguments []Request

// Request is the holds the information of the onewallet request
type RequestFormat struct {
	Arguments Arguments `json:"arguments"`
}
