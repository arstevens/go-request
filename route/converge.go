package route

import (
	"reflect"

	"github.com/arstevens/go-request/handle"
)

/* ConvergeChannels takes in an array of handle.Request streams and converges them out
to a single handle.Request stream. It can also receive new streams from the newStreams
channel in the event that new handlers are allocated on the fly. ConvergeChannels finishes
when all channels have closed */
func ConvergeChannels(inStreams []<-chan handle.Request, newStreams <-chan <-chan handle.Request, outStream chan<- handle.Request) {
	defer close(outStream)
	selectCases := make([]reflect.SelectCase, len(inStreams)+1)
	selectCases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(newStreams)}
	for i, stream := range inStreams {
		selectCases[i+1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(stream)}
	}

	for len(selectCases) > 0 {
		idx, value, ok := reflect.Select(selectCases)
		if !ok {
			selectCases = append(selectCases[:idx], selectCases[idx+1:]...)
			continue
		} else if idx == 0 {
			newStream := value.Interface().(<-chan handle.Request)
			selectCases = append(selectCases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(newStream)})
		} else {
			request := value.Interface().(handle.Request)
			outStream <- request
		}
	}
}
