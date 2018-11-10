package sessionmanager_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/orangejohny/SBWeb/internal/model"
	sm "github.com/orangejohny/SBWeb/internal/sessionmanager"
)

func TestInterfaceSession(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	SM, err := sm.InitConnSM(sm.Config{
		DBAddress:      `redis://user:@localhost:` + s.Port() + `/0`,
		TockenLength:   32,
		ExpirationTime: 1,
	})
	if err != nil {
		panic(err)
	}

	sID, err := SM.CreateSession(&model.Session{
		ID:        15,
		Login:     "aaa@eee.ru",
		UserAgent: "ieieie",
	})

	_, err = SM.CheckSession(sID)
	if err != nil {
		t.Error("Key must exist")
	}

	s.FastForward(5 * time.Second)

	res, err := SM.CheckSession(sID)
	fmt.Println(res)
	if res != nil {
		t.Error("Key mustn't exist")
	}

	sID, err = SM.CreateSession(&model.Session{
		ID:        15,
		Login:     "aaa@eee.ru",
		UserAgent: "ieieie",
	})

	_, err = SM.CheckSession(sID)
	if err != nil {
		t.Error("Key must exist")
	}

	SM.DeleteSession(sID)

	res, err = SM.CheckSession(sID)
	fmt.Println(res)
	if res != nil {
		t.Error("Key mustn't exist")
	}
}
