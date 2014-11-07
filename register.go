package main

type register struct {
	workers map[string]chan string
}

var reg = register{
	workers: make(map[string]chan string),
}

func (r *register) getWorkerChannel(id string) chan string {
	value, _ := r.workers[id]
	return value
}

func (r *register) exists(id string) bool {
	_, ok := r.workers[id]
	return ok
}

func (r *register) register(id string, wChan chan string) {
	r.workers[id] = wChan
}

func (r *register) unregister(id string) {
	delete(r.workers, id)
}

func (r *register) size() int {
	return len(r.workers)
}
