package raw

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mroyme/dogstatsd-local/internal/messages"
)

type Handler struct {
	Logger *log.Logger
}

func (H *Handler) New() messages.OutputHandler {
	return func(msg []byte) error {
		fmt.Println(string(msg))
		return nil
	}
}
