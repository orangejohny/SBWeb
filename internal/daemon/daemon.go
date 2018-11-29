package daemon

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"bmstu.codes/developers34/SBWeb/internal/model"

	"bmstu.codes/developers34/SBWeb/internal/api"
	"bmstu.codes/developers34/SBWeb/internal/db"
	sm "bmstu.codes/developers34/SBWeb/internal/sessionmanager"
)

// Config is config structure for whole service
type Config struct {
	DB  db.Config
	SM  sm.Config
	API api.Config
}

// RunService is function that starts the whole service
func RunService(cfg *Config) error {
	db, err := db.InitConnDB(cfg.DB)
	if err != nil {
		log.Fatalln("Can't connect to database", err.Error())
		return err
	}
	log.Println("Connected to DB")

	sm, err := sm.InitConnSM(cfg.SM)
	if err != nil {
		log.Fatalln("Can't start session manager", err.Error())
		return err
	}
	log.Println("Connected to SM")

	m := model.New(db, sm)

	log.Println("Starting API server...")
	srv, ch := api.StartServer(cfg.API, m)

	waitForSignal(srv, ch)

	return nil
}

func waitForSignal(srv *http.Server, chSrv chan error) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

LOOP:
	for {
		select {
		case s := <-ch:
			srv.Shutdown(nil)
			log.Printf("Got signal: %v, exiting.", s)
			<-chSrv
			break LOOP
		case err := <-chSrv:
			log.Fatalln(err)
		}
	}
}
