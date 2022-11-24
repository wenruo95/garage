package tcp

import "fmt"

type ServerStatus int

func (status ServerStatus) String() string {
	switch status {
	case ServerStatusInit:
		return fmt.Sprintf("Init(%v)", status)
	case ServerStatusServing:
		return fmt.Sprintf("Serving(%v)", status)
	case ServerStatusClosed:
		return fmt.Sprintf("Closed(%v)", status)
	default:
		return fmt.Sprintf("Unknown(%v)", status)
	}
}

var (
	ServerStatusInit    ServerStatus = 1
	ServerStatusServing ServerStatus = 2
	ServerStatusClosed  ServerStatus = 3
)

type ServerConnStatus int

func (status ServerConnStatus) String() string {
	switch status {
	case ServerConnStatusInit:
		return fmt.Sprintf("Init(%v)", status)
	case ServerConnStatusServing:
		return fmt.Sprintf("Serving(%v)", status)
	case ServerConnStatusClosed:
		return fmt.Sprintf("Closed(%v)", status)
	default:
		return fmt.Sprintf("Unknown(%v)", status)
	}
}

var (
	ServerConnStatusInit    ServerConnStatus = 1
	ServerConnStatusServing ServerConnStatus = 2
	ServerConnStatusClosed  ServerConnStatus = 3
)

const (
	ServerMinTimeoutMs     = 100
	ServerDefaultTimeoutMs = 30 * 1000

	MinClientHBTimeout = 3 * 60 * 1000 // ms
)
