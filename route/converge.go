package route

import (
	"reflect"
)

/* ConvergeChannels takes in an array of handle.Request streams and converges them out
to a single handle.Request stream. The out stream must be buffered. It can also receive
new streams from the newStreams channel in the event that new handlers are allocated on
the fly. ConvergeChannels finishes when all channels have closed */
func ConvergeChannels(inStreams []<-chan RequestPair, newStreams <-chan <-chan RequestPair, outStream chan<- RequestPair) {
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
			newStream := value.Interface().(<-chan RequestPair)
			selectCases = append(selectCases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(newStream)})
		} else {
			request := value.Interface().(RequestPair)
			outStream <- request
		}
	}
}
