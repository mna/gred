package main

import "sync"

type args struct {
	c, n        int
	wg          *sync.WaitGroup
	start, stop chan struct{}
	files       []jsonFile
}

func launchWorkers(a args) []*worker {
	ws := make([]*worker, a.c)
	ix := 0
	for i := 0; i < a.c; i++ {
		w := &worker{
			id: i,
			n:  a.n,
			wg: a.wg,
		}
		w.files = append(w.files, a.files[ix])
		go w.work(a.start, a.stop)
		ws[i] = w
		ix++
		if ix >= len(a.files) {
			ix = 0
		}
	}

	// If there are more files than workers, assign remaining files in order
	if len(a.files) > a.c {
		ix = 0
		for i := 0; i < len(a.files)-a.c; i++ {
			ws[ix].files = append(ws[ix].files, a.files[i])
			ix++
			if ix >= len(ws) {
				ix = 0
			}
		}
	}
	return ws
}

type worker struct {
	id      int
	n       int
	wg      *sync.WaitGroup
	files   []jsonFile
	replies []interface{}
}

func (w *worker) work(start, stop chan struct{}) {
	// Wait for start signal
	<-start
	// TODO : Work...
	w.wg.Done()
}
