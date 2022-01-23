package timer

import (
	"container/heap"
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

const Forever = -1

// ID represents ID of timer task
type ID int64

func (id ID) Valid() bool { return id > 0 }

// Task represents timer task
type Task interface {
	ExecTimer(ID)
}

// TaskFunc wraps function as a task
type TaskFunc func(ID)

// ExecTimer implements Task ExecTimer method
func (fn TaskFunc) ExecTimer(id ID) { fn(id) }

// Scheduler schedules timers
type Scheduler interface {
	// Start starts the Scheduler
	Start() error
	// Shutdown shutdowns the Scheduler
	Shutdown()
	// Add adds a new timer task
	Add(next, duration time.Duration, task Task, times int) ID
	// Remove removes a timer task by ID
	Remove(id ID)
}

type group struct {
	next   int64
	timers []timer
}

func (g *group) remove(id ID) bool {
	n := len(g.timers)
	for i := 0; i < n; i++ {
		if g.timers[i].id == id {
			copy(g.timers[i:n-1], g.timers[i+1:])
			g.timers = g.timers[:n-1]
			return true
		}
	}
	return false
}

type timer struct {
	id       ID
	task     Task
	times    int
	next     int64
	duration int64
}

// memoryScheduler implements Scheduler in memory
type memoryScheduler struct {
	groups  []group
	indices map[int64]int // next => indexof(groups)
	timers  map[ID]int64  // id => next

	addChan   chan timer
	closeChan chan struct{}

	nextId  int64
	running int32

	queue  *list.List
	locker sync.Mutex
	cond   *sync.Cond
}

// NewMemoryScheduler creates in-memory Scheduler
func NewMemoryScheduler() Scheduler {
	s := &memoryScheduler{
		indices:   make(map[int64]int),
		timers:    make(map[ID]int64),
		addChan:   make(chan timer, 128),
		closeChan: make(chan struct{}),
		queue:     list.New(),
	}
	s.cond = sync.NewCond(&s.locker)
	return s
}

// Start implements Scheduler Start method
func (s *memoryScheduler) Start() error {
	if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
		return nil
	}
	go s.receive()
	go s.schedule()
	return nil
}

// Shutdown implements Scheduler Shutdown method
func (s *memoryScheduler) Shutdown() {
	if atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		s.cond.Signal()
		s.closeChan <- struct{}{}
	}
}

func (s *memoryScheduler) receive() {
	for atomic.LoadInt32(&s.running) == 1 {
		s.cond.L.Lock()
		for s.queue.Len() == 0 {
			s.cond.Wait()
		}
		front := s.queue.Front()
		task := front.Value.(timer)
		s.queue.Remove(front)
		s.cond.L.Unlock()
		s.addChan <- task
	}
}

func (s *memoryScheduler) schedule() {
	var timer *time.Timer
	var timerExpired bool
	for {
		if len(s.groups) == 0 {
			select {
			case x := <-s.addChan:
				if x.id < 0 {
					s.removeTimer(-x.id)
				} else {
					s.addTimer(x)
				}
			case <-s.closeChan:
				return
			}
		}
		now := time.Duration(time.Now().UnixNano()/1000000) * time.Millisecond

		first := heap.Pop(s).(group)
		next := first.next
		dt := time.Duration(next)*time.Millisecond - now

		if dt <= 0 {
			s.execGroup(first)
			continue
		}

		if timer == nil {
			timer = time.NewTimer(dt)
		} else {
			if !timerExpired {
				if !timer.Stop() {
					<-timer.C
				}
			}
			timer.Reset(dt)
		}
		timerExpired = false

	WAIT_DO_FIRST:
		for {
			select {
			case <-timer.C:
				timerExpired = true
				s.execGroup(first)
				break WAIT_DO_FIRST
			case x := <-s.addChan:
				if x.id < 0 {
					next, ok := s.timers[-x.id]
					if ok {
						if next == first.next {
							delete(s.timers, -x.id)
							first.remove(-x.id)
						} else {
							s.removeTimerByNext(-x.id, next)
						}
					}
				} else {
					if x.next == first.next {
						first.timers = append(first.timers, x)
					} else {
						s.addTimer(x)
						if x.next < first.next {
							heap.Push(s, first)
							break WAIT_DO_FIRST
						}
					}
				}
			case <-s.closeChan:
				return
			}
		}
	}
}

// Add implements Scheduler Add method
func (s *memoryScheduler) Add(next, duration time.Duration, task Task, times int) ID {
	id := ID(atomic.AddInt64(&s.nextId, 1))

	s.locker.Lock()
	t := timer{
		id:       id,
		times:    times,
		task:     task,
		duration: int64(duration / time.Millisecond),
		next:     int64(next / time.Millisecond),
	}
	s.queue.PushBack(t)
	l := s.queue.Len()
	s.locker.Unlock()

	if l == 1 {
		s.cond.Signal()
	}

	return id
}

// Remove implements Scheduler Remove method
func (s *memoryScheduler) Remove(id ID) {
	s.locker.Lock()
	t := timer{
		id: -id,
	}
	s.queue.PushBack(t)
	l := s.queue.Len()
	s.locker.Unlock()

	if l == 1 {
		s.cond.Signal()
	}
}

func (s *memoryScheduler) addTimer(x timer) {
	s.timers[x.id] = x.next
	if i, ok := s.indices[x.next]; ok {
		s.groups[i].timers = append(s.groups[i].timers, x)
	} else {
		g := group{
			next:   x.next,
			timers: make([]timer, 0, 8),
		}
		g.timers = append(g.timers, x)
		heap.Push(s, g)
	}
}

func (s *memoryScheduler) removeTimer(id ID) {
	next, ok := s.timers[id]
	if !ok {
		return
	}
	s.removeTimerByNext(id, next)
}

func (s *memoryScheduler) removeTimerByNext(id ID, next int64) {
	delete(s.timers, id)

	i, ok := s.indices[next]
	if !ok {
		return
	}

	if s.groups[i].remove(id) {
		if len(s.groups[i].timers) == 0 {
			heap.Remove(s, i)
		}
	}
}

func (s *memoryScheduler) execGroup(g group) {
	n := 0
	for i := range g.timers {
		g.timers[i].task.ExecTimer(g.timers[i].id)
		if g.timers[i].times > 0 {
			g.timers[i].times--
		}
		if g.timers[i].times != 0 {
			g.timers[i].next += g.timers[i].duration
			if i != n {
				g.timers[n] = g.timers[i]
				n++
			}
		} else {
			delete(s.timers, g.timers[i].id)
		}
	}
	g.timers = g.timers[:n]
	for i := range g.timers {
		s.addTimer(g.timers[i])
	}
}

// Len implements heap.Interface Len method
func (s *memoryScheduler) Len() int { return len(s.groups) }

// Less implements heap.Interface Less method
func (s *memoryScheduler) Less(i, j int) bool { return s.groups[i].next < s.groups[j].next }

// Swap implements heap.Interface Swap method
func (s *memoryScheduler) Swap(i, j int) {
	s.groups[i], s.groups[j] = s.groups[j], s.groups[i]
	s.indices[s.groups[i].next] = i
	s.indices[s.groups[j].next] = j
}

// Push implements heap.Interface Push method
func (s *memoryScheduler) Push(x any) {
	g := x.(group)
	l := len(s.groups)
	s.groups = append(s.groups, g)
	s.indices[g.next] = l
}

// Pop implements heap.Interface Pop method
func (s *memoryScheduler) Pop() any {
	l := len(s.groups)
	x := s.groups[l-1]
	s.groups = s.groups[:l-1]
	delete(s.indices, x.next)
	return x
}

var globalScheduler = NewMemoryScheduler()

func init() {
	globalScheduler.Start()
}

// SetTimeout add a timeout timer
func SetTimeout(d time.Duration, task Task) ID {
	return globalScheduler.Add(time.Duration(time.Now().Add(d).UnixNano()), d, task, 1)
}

// SetTimeoutFunc add a timeout timer func
func SetTimeoutFunc(d time.Duration, fn TaskFunc) ID {
	return SetTimeout(d, fn)
}

// ClearTimeout removes the timeout timer by ID
func ClearTimeout(id ID) {
	globalScheduler.Remove(id)
}

// SetInterval add interval timer
func SetInterval(d time.Duration, task Task) ID {
	return globalScheduler.Add(time.Duration(time.Now().Add(d).UnixNano()), d, task, Forever)
}

// SetIntervalFunc add a interval timer func
func SetIntervalFunc(d time.Duration, fn TaskFunc) ID {
	return SetInterval(d, fn)
}

// ClearTimeout removes the interval timer by ID
func ClearInterval(id ID) {
	globalScheduler.Remove(id)
}
