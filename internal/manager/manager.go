package manager

import (
	"context"
	"sync"
	"time"

	"github.com/BON4/timedQ/pkg/ttlstore"
	"github.com/sirupsen/logrus"
)

type TaskType int8

const (
	GetTask TaskType = iota
	SetTask
)

type Task struct {
	Key      string
	Val      string
	RespChan chan string
	Type     TaskType

	mapIndex int
}

type Worker struct {
	index        int
	valTTL       time.Duration
	store        *ttlstore.MapStore[string, string]
	logger       *logrus.Entry
	notFoundChan chan *Task
	next         *Worker
}

func newWorker(index int,
	valTTL time.Duration,
	store *ttlstore.MapStore[string, string],
	logger *logrus.Entry,
	notFoundChan chan *Task,
	next *Worker) *Worker {
	return &Worker{
		index:        index,
		valTTL:       valTTL,
		store:        store,
		logger:       logger,
		notFoundChan: notFoundChan,
		next:         next,
	}
}

func (w *Worker) Listen(ctx context.Context, taskChan chan *Task, wg *sync.WaitGroup) {
	defer wg.Done()

	respond := func(t *Task) {
		val, ok := w.store.Get(ctx, t.Key)
		if !ok {
			// w.logger.Infof("Not Found. Passing to Next Worker. Key: %s", t.Key)
			if t.mapIndex == -1 {
				t.mapIndex = w.index
			}
			w.next.notFoundChan <- t
			return
		}

		t.RespChan <- val

		//Refresh TTL
		if err := w.store.Set(ctx, t.Key, val, w.valTTL); err != nil {
			w.logger.Errorf("got error while refreshing value: %s", err.Error())
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-taskChan:
			switch t.Type {
			case GetTask:
				// w.logger.Infof("Getting. Key: %s", t.Key)
				respond(t)
			case SetTask:
				w.logger.Info("Setting.")

				if err := w.store.Set(ctx, t.Key, t.Val, w.valTTL); err != nil {
					w.logger.Errorf("got error while setting key-value: %s", err.Error())
				}
			}
		case t := <-w.notFoundChan:
			if t.mapIndex == w.index {
				t.RespChan <- ""
			} else {
				respond(t)
			}
		}
	}
}

type WorkerRing struct {
	len int
	cur *Worker
}

func (wr *WorkerRing) Push(w *Worker) {
	if wr.cur == nil {
		wr.cur = w
	} else if wr.len >= 1 {
		w.next = wr.cur.next
	}

	// 0 -> 1
	//   <-
	wr.cur.next = w
	wr.cur = w
	wr.len++
}

func (wr *WorkerRing) Range(f func(w *Worker)) {
	temp := wr.cur
	for {
		f(temp)
		temp = temp.next
		if temp.index == wr.cur.index {
			break
		}
	}
}

type WorkerManager struct {
	logger *logrus.Logger

	ctx    context.Context
	cancel context.CancelFunc
	waitG  *sync.WaitGroup

	ReqChan     chan *Task
	WorkerArena *WorkerRing
}

// NewWorkerManager - creates new worker manager, length of stroes MUST be == to cfg.Manager.WorkerNum
func NewWorkerManager(ctx context.Context, stores []*ttlstore.MapStore[string, string], logger *logrus.Logger, cfg ManagerConfig) *WorkerManager {
	//TODO CHAN SIZE???
	wm := &WorkerManager{
		ReqChan:     make(chan *Task, 100),
		WorkerArena: &WorkerRing{},
		waitG:       &sync.WaitGroup{},
		logger:      logger,
	}

	// Build a cercualr list of workers
	for widx := 0; widx < int(cfg.WorkerNum); widx++ {
		wm.WorkerArena.Push(newWorker(
			widx,
			cfg.ValTTL,
			stores[widx],
			logger.WithField("worker", widx),
			make(chan *Task, 10), nil),
		)
	}

	wm.ctx, wm.cancel = context.WithCancel(ctx)

	return wm
}

func (wm *WorkerManager) Get(key string) string {
	respChan := make(chan string, 1)

	t := &Task{
		mapIndex: -1,
		Key:      key,
		RespChan: respChan,
		Type:     GetTask,
	}
	wm.ReqChan <- t

	return <-respChan
}

func (wm *WorkerManager) Set(key string, val string) {
	t := &Task{
		mapIndex: -1,
		Key:      key,
		Val:      val,
		Type:     SetTask,
	}
	wm.ReqChan <- t
}

func (wm *WorkerManager) Run() {
	wm.WorkerArena.Range(func(w *Worker) {
		wm.logger.Infof("Worker%d. Listening.", w.index)
		go w.Listen(wm.ctx, wm.ReqChan, wm.waitG)
		wm.waitG.Add(1)
	})

	wm.logger.Info("Running...")
}

func (wm *WorkerManager) Stop() {
	wm.logger.Info("Stoping...")
	wm.cancel()
	wm.waitG.Wait()
}
