package main

import (
	"log"
	"net/http"

	"encoding/base64"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	"github.com/volatiletech/authboss/v3"
	_ "github.com/volatiletech/authboss/v3/auth"
	_ "github.com/volatiletech/authboss/v3/logout"
)

var (
	database     = NewMemStorage()
	ab           = authboss.New()
	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer
)

// ExampleServer is example server for AuthBoss PoC.
// It provides
type ExampleServer struct {
	host    string
	port    string
	storage *MemStorage
	router  *mux.Router
}

// NewExampleServer creates new instance of ExampleServer.
func NewExampleServer() *ExampleServer {
	es := &ExampleServer{}
	es.port = "8080"
	es.host = "http://localhost"
	es.storage = database
	es.router = mux.NewRouter()
	es.routes()

	return es
}

func (es *ExampleServer) routes() {
	es.router.Use(ab.LoadClientStateMiddleware)
	es.router.HandleFunc("", es.status())
	es.router.HandleFunc("/foo", es.foo())
	es.router.HandleFunc("/bar", es.bar())
	es.router.HandleFunc("/sigma", es.sigma())
}

func (es *ExampleServer) status() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func (es *ExampleServer) foo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("foo"))
	}
}

func (es *ExampleServer) bar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bar"))
	}
}

func (es *ExampleServer) sigma() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := ab.CurrentUserID(r)
		if err != nil {
			log.Println(err)
		}

		switch hasAdminPermissions(user) {
		case false:
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("You must be admin user"))
		default:
			w.Write([]byte("sigma - Congratulations, admin user!"))
		}
	}
}

func main() {
	es := NewExampleServer()

	ab.Config.Storage.Server = database
	ab.Config.Paths.RootURL = es.host + ":" + es.port

	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	cookieStore := abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = false
	sessionStore := abclientstate.NewSessionStorer("test", sessionStoreKey, nil)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false

	log.Fatal(http.ListenAndServe(":"+es.port, es.router))
}
