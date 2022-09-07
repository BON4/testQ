package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/BON4/timedQ/internal/manager"
	"github.com/BON4/timedQ/pkg/ttlstore"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// TODO: this constans shoud be loaded from cfg
var LOGGER_FILE_PATH = ""

func setUpLogger(fileName string) (*logrus.Logger, error) {
	// instantiation
	logger := logrus.New()

	if len(fileName) == 0 {
		logger.Out = os.Stdout
	} else {
		//Write to file
		src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return nil, err
		}
		//Set output
		logger.Out = src
	}

	//Set log level
	logger.SetLevel(logrus.DebugLevel)

	//Format log
	logger.SetFormatter(&logrus.TextFormatter{})
	return logger, nil
}

type Server struct {
	g      *gin.Engine
	logger *logrus.Logger
	wM     *manager.WorkerManager
	cfg    ServerConfig
	stores []*ttlstore.MapStore[string, string]
}

func NewServer(configPath string) (*Server, error) {
	ctx := context.Background()
	g := gin.Default()

	cfg, err := LoadServerConfig(configPath)
	if err != nil {
		return nil, err
	}

	log, err := setUpLogger(cfg.AppConfig.LogFile)
	if err != nil {
		return nil, err
	}

	log.Infof("Loaded config: %+v", cfg)

	//Construct maps
	stores := make([]*ttlstore.MapStore[string, string], cfg.ManagerCfg.WorkerNum)
	for i := uint(0); i < cfg.ManagerCfg.WorkerNum; i++ {
		ttlCfg := cfg.StoreCfg
		ttlCfg.SavePath = strings.TrimRight(ttlCfg.SavePath, "/") + fmt.Sprintf("/#store%d.db", i)

		log.Infof("Creating db file in: %s", ttlCfg.SavePath)

		stores[i] = ttlstore.NewMapStore[string, string](ctx, ttlCfg)
	}

	wM := manager.NewWorkerManager(ctx, stores, log, cfg.ManagerCfg)

	return &Server{
		g:      g,
		logger: log,
		stores: stores,
		wM:     wM,
		cfg:    cfg,
	}, nil
}

func (s *Server) Run() error {
	srv := &http.Server{
		Handler: s.g,
		Addr:    s.cfg.AppConfig.Port,
	}

	s.logger.Infof("Running on: %s", s.cfg.AppConfig.Port)

	// start every store
	for _, st := range s.stores {
		if err := st.Load(); err != nil {
			s.logger.Errorf("Error while load store: %s", err.Error())
		}

		if err := st.Run(); err != nil {
			s.logger.Errorf("Error while start store: %s", err.Error())
		}
	}

	//start manager
	s.wM.Run()

	if err := s.MapHandlers(); err != nil {
		return err
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Errorf("listen: %s\n", err)
			return
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Stop manager
	s.wM.Stop()

	// Stop every store
	for _, st := range s.stores {
		if err := st.Close(); err != nil {
			s.logger.Errorf("Error while closing store: %s", err.Error())
		}
	}

	s.logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error("Server Shutdown Err:", err)
		return err
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		s.logger.Info("timeout of 5 seconds.")
	}
	s.logger.Info("Server exiting")
	return nil
}
