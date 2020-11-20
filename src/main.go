package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

//Sets up a logger, a webHookServer, prints the address and port, starts the server
func main() {
	logger := log.New(os.Stdout, "", 0)

	webHookServ := setupServer(logger)

	logger.Printf("KubeLinterBot is listening on http://localhost%s\n", webHookServ.Addr) //TODO: Address

	webHookServ.ListenAndServe()

	/*
		react to webhook
			check if .yaml was changed
				Call KubeLinter
				if exit-code 0
					no review comment? Some other kind of feedback?
				else
					interpret Linter-output
					review comment via github
			else
				do nothing? Feedback?
	*/
}

//Setup method, needs an already set up logger and returns a http.Server-Pointer
func setupServer(logger *log.Logger) *http.Server {
	return &http.Server{
		Addr:    ":4567", //TODO: hardcoded is bad
		Handler: newServer(logWith(logger)),

		//TODO: standard, check if that makes sense
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

type Server struct {
	mux    *http.ServeMux
	logger *log.Logger
}

func newServer(options ...Option) *Server {
	s := &Server{logger: log.New(ioutil.Discard, "", 0)}

	for _, o := range options {
		o(s)
	}

	s.mux = http.NewServeMux()

	//Multiplexer, maybe need several paths?
	s.mux.HandleFunc("/", s.index)
	s.mux.HandleFunc("/push/", s.push)
	s.mux.HandleFunc("/pull/", s.pull)

	return s
}

type Option func(*Server)

func logWith(logger *log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.log("Only POST allowed.")
		//return
	}

	reqBody, err := ioutil.ReadAll(r.Body) //TODO ReadAll may be bad for large messages
	if err != nil {
		log.Fatal(err)
	}
	s.log("Webhook received.")
	makeJSON(s, reqBody)
}

//TODO: Parse JSON. Marshal? Decode?
func makeJSON(s *Server, body []byte) {
	s.log("\n%s", body)
}

func (s *Server) log(format string, v ...interface{}) {
	s.logger.Printf(format+"\n", v...)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("KubeLinterBot is running here."))
}

func (s *Server) push(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("KubeLinterBot has received a push"))
	fmt.Println("push")
}

func (s *Server) pull(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("KubeLinterBot has received a pull"))
	fmt.Println("pull")
}
