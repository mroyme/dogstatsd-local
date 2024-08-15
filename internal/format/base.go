package format

import "github.com/mroyme/dogstatsd-local/internal/messages"

type Handler interface {
	New() messages.OutputHandler
}
