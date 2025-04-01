package control

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type SingleFight struct {
	mu sync.Mutex
	m  map[string]*call
}

func (s *SingleFight) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	s.mu.Lock()
	if s.m == nil {
		s.m = make(map[string]*call)
	}
	if c, ok := s.m[key]; ok {
		s.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	s.m[key] = c
	s.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	s.mu.Lock()
	delete(s.m, key)
	s.mu.Unlock()

	return c.val, c.err
}
