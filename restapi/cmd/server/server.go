// Starting the API server
package main // import "github.com/la4ezar/restapi"

import (
	"context"
	"sync"

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

	cfg, err := config.New()
	fatalOnError(err)

	err = cfg.Validate()
	fatalOnError(err)

	ctx, err = log.Configure(ctx, cfg.Logger)
	fatalOnError(err)

	db, err := storage.New(cfg.Storage)
	fatalOnError(err)
	defer func() {
		if err := db.Close(); err != nil {
			log.D().WithError(err).Error()
		}
	}()

	//r := routes-old.NewRouter(db.DB)

	//srv := server.New(cfg.Server, r)
	repository := storage.NewRepository(*db)
	ctr := controller.NewController(repository)
	srv := server.New(cfg.Server, *ctr)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go srv.Start(ctx, wg)

	wg.Wait()

	log.D().Println("Server shutdown.")
}

func fatalOnError(err error) {
	if err != nil {
		log.D().Fatal(err.Error())
	}
}
