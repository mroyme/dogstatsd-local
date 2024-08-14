package raw

import (
	"fmt"
	"github.com/mroyme/dogstatsd-local/internal/handler"
)

func NewHandler() handler.MessageHandler {
	return func(msg []byte) error {
		fmt.Println(string(msg))
		return nil
	}
}
