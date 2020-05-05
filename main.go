package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	handler "github.com/image-store/webservice/handlers"
	"github.com/image-store/webservice/kafka-services"
	"github.com/rs/cors"
	"golang.org/x/net/context"
)

var (
	nc     kafka.NewConsumer
	logger *log.Logger
)

func main() {

	logger = log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile)
	imageAlbum := handler.NewImageAlbum(logger)
	// using mux router instead of http.ServerMux to take care of my routing
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("Get")
	router.Handle("/api/v1/create-album", imageAlbum).Methods("POST")
	router.Handle("/api/v1/create-image", imageAlbum).Methods("POST")
	router.Handle("/api/v1/delete-album", imageAlbum).Methods("DELETE")
	router.Handle("/api/v1/delete-image", imageAlbum).Methods("DELETE")
	router.Handle("/api/v1/image", imageAlbum).Methods("GET")
	router.Handle("/api/v1/images", imageAlbum).Methods("GET")

	server := &http.Server{
		Addr:         ":9098",
		Handler:      cors.Default().Handler(gziphandler.GzipHandler(noCacheMW(router))),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 40 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Println("something went wrong while starting the server", err)
			return
		}
	}()

	signals := make(chan os.Signal, 0)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)

	<-signals
	logger.Println("recieved termination, gracefull shutdown", <-signals)

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(timeOutCtx)

	os.Exit(0)
}

func noCacheMW(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		h.ServeHTTP(w, r)
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome! Please hit the `/api/v1/create-album` API to create new album"))
}
