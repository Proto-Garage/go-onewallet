package onewallet

import "github.com/djansyle/go-rabbit"

func requestFormatter(request interface{}) (interface{}, error) {
	request = RequestFormat{ Arguments: Arguments{request.(Request)}}

	return request, nil
}

// CreateClient creates an client that is used for onewallet communication
func CreateClient(opt *rabbit.CreateClientOption) (rabbit.Sender, error) {
	return rabbit.CreateClient(opt, requestFormatter)
}