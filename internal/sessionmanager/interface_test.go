package sessionmanager_test

import (
	"fmt"
	"testing"
	"time"

	"bmstu.codes/developers34/SBWeb/internal/model"
	sm "bmstu.codes/developers34/SBWeb/internal/sessionmanager"
	"github.com/alicebob/miniredis"
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
	}, true)

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
	}, true)

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
