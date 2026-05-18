package raw

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/mroyme/dogstatsd-local/internal/messages"
)

type Handler struct {
	Logger *log.Logger
}

func (h *Handler) New() messages.OutputHandler {
	return func(msg messages.DogStatsDMessage) error {
		if data := msg.Data(); data != nil {
			fmt.Println(string(data))
		}
		return nil
	}
}
