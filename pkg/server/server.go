package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/tenmozes/sanchez/pkg/spirit"
)

type Server struct {
	maker  *spirit.Service
	router *httprouter.Router
	s      http.Server
}

func NewServer(maker *spirit.Service) *Server {
	s := &Server{
		maker: maker,
	}
	s.router = httprouter.New()
	s.router.POST("/cocktail/:name", s.MakeCocktail)
	return s
}

func (s *Server) MakeCocktail(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	err := s.maker.Make(params.ByName("name"))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (s *Server) Run(port int) error {
	s.s = http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.router,
	}
	log.Printf("http server starts on %d", port)
	return s.s.ListenAndServe()
}

func (s *Server) Close() error {
	s.maker.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.s.Shutdown(ctx)
}
