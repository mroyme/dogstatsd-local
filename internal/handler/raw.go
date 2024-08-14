package handler

import "fmt"

func NewRawMessageHandler() MessageHandler {
	return func(msg []byte) error {
		fmt.Println(string(msg))
		return nil
	}
}
