package handler

import "log"

func HandleAmqpMsg(msg interface{}) error {
	log.Println(msg)
	return nil
}
