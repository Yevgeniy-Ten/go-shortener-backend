package shutdown

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GraceFullShutDown is a function that handles graceful shutdown of the server
func GraceFullShutDown(server *http.Server, timeout time.Duration) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigs
	log.Println("graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error when shutdown %v", err)
	} else {
		log.Println("Server is shutdowned")
	}
}
