package sessionmanager_test

import (
	"testing"

	sm "bmstu.codes/developers34/SBWeb/internal/sessionmanager"
	"github.com/alicebob/miniredis"
)

func TestInitSM(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// test initSM with bad address
	_, err = sm.InitConnSM(sm.Config{
		DBAddress:      "fakeDBaddress",
		TockenLength:   100,
		ExpirationTime: 100,
	})

	if err == nil {
		t.Error("SM can't be initiated with such address")
	}

	// test initSM with good address
	SM, err := sm.InitConnSM(sm.Config{
		DBAddress:      `redis://user:@localhost:` + s.Port() + `/0`,
		TockenLength:   100,
		ExpirationTime: 100,
	})
	if err != nil {
		t.Error("SM must have been initiated")
	}

	SM.CreateSession(nil)
}
