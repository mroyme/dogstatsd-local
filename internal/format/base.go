package format

import "github.com/mroyme/dogstatsd-local/v2/internal/messages"

type Handler interface {
	New() messages.OutputHandler
}
