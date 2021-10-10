package listener

import (
	"fmt"
	"github.com/Rohan12152001/Syook_TimeSeries/listener/manager/listener/data"
)

var (
	DiscardedError = fmt.Errorf("object discarded")
)

type ListenerManager interface {
	DecryptAndEmit(enString, myListenerId string) (data.LiveData, error)
}
