// Starting the API server
package main // import "github.com/la4ezar/restapi"

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/la4ezar/restapi/internal/config"
	"github.com/la4ezar/restapi/pkg/controller"
	"github.com/la4ezar/restapi/pkg/log"
	"github.com/la4ezar/restapi/pkg/server"
	"github.com/la4ezar/restapi/pkg/storage"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handleInterrupts(ctx, cancel)

	cfg, err := config.NewDefaultServerConfig()
	fatalOnError(err)

	err = cfg.Validate()
	fatalOnError(err)

	ctx, err = log.Configure(ctx, cfg.Logger)
	fatalOnError(err)

	db, err := storage.New(cfg.Storage)
	fatalOnError(err)
	defer func() {
		if err := db.Close(); err != nil {
			log.C(ctx).WithError(err).Error()
		}
	}()

	repository := storage.NewRepository(*db)
	ctr := controller.NewController(repository)
	srv := server.New(cfg.Server, *ctr)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go srv.Start(ctx, wg)

	wg.Wait()

	//log.C(ctx).Println("Server shutdown.")
}

func fatalOnError(err error) {
	if err != nil {
		log.D().Fatal(err.Error())
	}
}

func handleInterrupts(ctx context.Context, cancel context.CancelFunc) {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		select {
		case <-term:
			log.C(ctx).Println("Received OS Interrupt/Kill/Terminate signal, exiting gracefully...")
			cancel()
		case <-ctx.Done():
			return
		}
	}()
}
