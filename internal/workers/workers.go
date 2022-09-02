package workers

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/BON4/timedQ/pkg/ttlstore"
)

type Task struct {
	Key      string
	RespChan chan string
}

type Worker struct {
	valTTL time.Duration
	store  *ttlstore.MapStore[string, string]
}

func (w *Worker) Listen(ctx context.Context, taskChan chan *Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-taskChan:
			val, ok := w.store.Get(ctx, t.Key)
			if !ok {
				//DO SOME WORK TO FIGURE OUT WHAT VAL IS
				//...
				//Mutate val with found value
				val = strings.Repeat(t.Key, 3)
			}

			t.RespChan <- val

			//Refresh TTL
			if err := w.store.Set(ctx, t.Key, val, w.valTTL); err != nil {
				fmt.Printf("Got error: %s", err.Error())
			}
		}
	}
}

type WorkerManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	waitG  *sync.WaitGroup

	ReqChan     chan *Task
	WorkerArena []*Worker
}

// NewWorkerManager - creates new worker manager, length of stroes MUST be == to cfg.Manager.WorkerNum
func NewWorkerManager(ctx context.Context, stores []*ttlstore.MapStore[string, string], cfg ManagerConfig) *WorkerManager {
	//TODO CHAN SIZE???
	wm := &WorkerManager{
		ReqChan:     make(chan *Task, 100),
		WorkerArena: make([]*Worker, cfg.Manager.WorkerNum),
		waitG:       &sync.WaitGroup{},
	}

	for widx := 0; widx < len(wm.WorkerArena); widx++ {
		wm.WorkerArena[widx] = &Worker{
			valTTL: cfg.Manager.ValTTL,
			store:  stores[widx],
		}
	}

	wm.ctx, wm.cancel = context.WithCancel(ctx)

	return wm
}

func (wm *WorkerManager) Get(key string) string {
	respChan := make(chan string, 1)

	t := &Task{
		Key:      key,
		RespChan: respChan,
	}
	wm.ReqChan <- t

	return <-respChan
}

func (wm *WorkerManager) Run() {
	for _, worker := range wm.WorkerArena {
		go worker.Listen(wm.ctx, wm.ReqChan, wm.waitG)
		wm.waitG.Add(1)
	}
}

func (wm *WorkerManager) Stop() {
	wm.cancel()
	wm.waitG.Wait()
}
