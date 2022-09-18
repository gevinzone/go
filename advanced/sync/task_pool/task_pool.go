package task_pool

type TaskPool interface {
	Do(f func())
}

type TaskPoolWithToken struct {
	ch chan struct{}
}

func (taskPool *TaskPoolWithToken) Do(f func()) {
	token := <-taskPool.ch
	go func() {
		f()
		taskPool.ch <- token
	}()
	//f()
	//taskPool.ch <- token
}

func NewTaskPoolWithToken(limit int) *TaskPoolWithToken {
	taskPool := &TaskPoolWithToken{
		ch: make(chan struct{}, limit),
	}
	for i := 0; i < limit; i++ {
		taskPool.ch <- struct{}{}
	}
	return taskPool
}

type TaskPoolWithCache struct {
	limit int
	ch    chan func()
}

func (t *TaskPoolWithCache) Do(f func()) {
	t.ch <- f
}

func (t *TaskPoolWithCache) Start() {
	for i := 0; i < t.limit; i++ {
		go func() {
			for {
				select {
				case task, ok := <-t.ch:
					if !ok {
						return
					}
					task()
				}
			}
		}()
	}
}

func (t *TaskPoolWithCache) Stop() {
	close(t.ch)
}

func (t *TaskPoolWithCache) Restart() {
	t.ch = make(chan func(), t.limit)
	t.Start()
}

func NewTaskPoolWithCache(limit int) *TaskPoolWithCache {
	t := &TaskPoolWithCache{
		limit: limit,
		ch:    make(chan func(), limit),
	}
	return t
}
