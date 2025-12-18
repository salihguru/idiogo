package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/salihguru/idiogo/internal/app/serve"
	"github.com/salihguru/idiogo/internal/rest"
	"github.com/salihguru/idiogo/pkg/server"
)

func init() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()
	if err := serve.Init(ctx); err != nil {
		panic(err)
	}
}

func main() {
	a := serve.Get()
	restServer := rest.New(rest.Config{
		Rest:      a.Config.Rest,
		I18n:      *a.Deps.I18n,
		Validator: *a.Deps.ValidationSrv,
		Locales:   a.Config.I18n.Locales,
		Routers:   a.Modules.Routers(),
	})
	wg := sync.WaitGroup{}
	wg.Add(2)
	server.Start("rest", restServer, wg.Done)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)
	go func() {
		defer wg.Done()
		<-shutdownCh
		log.Println("application is shutting down...")
		if err := a.Shutdown(context.Background(), restServer.Shutdown); err != nil {
			log.Fatalf("failed to disconnect: %v", err)
		}
	}()

	wg.Wait()
	fmt.Println("All servers are stopped.")
}
