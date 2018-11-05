package daemon

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/orangejohny/SBWeb/internal/api"
	"github.com/orangejohny/SBWeb/internal/db"
	"github.com/orangejohny/SBWeb/internal/model"
	sm "github.com/orangejohny/SBWeb/internal/sessionmanager"
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

	sm, err := sm.InitConnSM(cfg.SM)
	if err != nil {
		log.Fatalln("Can't start session manager", err.Error())
		return err
	}

	m := model.New(db, sm)

	log.Println("Starting API server...")
	go api.StartServer(cfg.API, m)

	waitForSignal()

	return nil
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
