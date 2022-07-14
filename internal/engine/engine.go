package engine

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alwindoss/ark"
	"github.com/alwindoss/ark/internal/vault"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httptransport "github.com/go-kit/kit/transport/http"
)

func Run(cfg *ark.Config) error {
	repo := vault.NewFSRepository(cfg.TempFolder)
	svc := vault.NewService(repo)

	saveHandler := httptransport.NewServer(
		makeSaveEndpoint(svc),
		decodeSaveRequest,
		encodeSaveResponse,
	)

	retrieveHandler := httptransport.NewServer(
		makeRetrieveEndpoint(svc),
		decodeRetrieveRequest,
		encodeRetrieveResponse,
	)

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post("/vault", saveHandler.ServeHTTP)
	r.Get("/vault/{key}", retrieveHandler.ServeHTTP)
	addr := fmt.Sprintf(":%d", cfg.Port)
	http.ListenAndServe(addr, r)
	return nil
}

// func RunGin(cfg *ark.Config) error {

// 	repo := vault.NewFSRepository(cfg.TempFolder)
// 	svc := vault.NewService(repo)

// 	saveHandler := httptransport.NewServer(
// 		makeSaveEndpoint(svc),
// 		decodeSaveRequest,
// 		encodeSaveResponse,
// 	)

// 	retrieveHandler := httptransport.NewServer(
// 		makeRetrieveEndpoint(svc),
// 		decodeRetrieveRequest,
// 		encodeRetrieveResponse,
// 	)

// 	// http.Handle("/uppercase", uppercaseHandler)
// 	// http.Handle("/count", countHandler)

// 	r := gin.Default()
// 	r.GET("/ping", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "pong",
// 		})
// 	})
// 	r.POST("/vault", gin.WrapH(saveHandler))
// 	r.GET("/vault/:key", gin.WrapH(retrieveHandler))
// 	addr := fmt.Sprintf(":%d", cfg.Port)
// 	err := r.Run(addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
