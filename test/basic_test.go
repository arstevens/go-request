package rtest

/*
func TestRequestLibrary(t *testing.T) {
	requestCount := 100
	datapoints := make([]byte, requestCount)
	for i := 0; i < requestCount; i++ {
		datapoints[i] = byte(rand.Int())
	}

	listener := &TestListener{datapoints}
	done := make(chan struct{})

	gen := &TestGenerator{cap: 3}
	alloc := allocate.NewCyclicJobAllocator(10, gen)
	handlers := map[int]handle.RequestHandler{
		0: &TestHandler{capacity: 10},
		1: &TestHandler{capacity: 15},
		2: &TestHandler{capacity: 10},
		3: alloc,
	}

	go route.UnpackAndRoute(listener, done, handlers, UnpackTestRequest, ReadTestRequestFromConn)
	time.Sleep(time.Second * 5)
}
*/
