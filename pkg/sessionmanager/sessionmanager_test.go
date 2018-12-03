// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package sessionmanager_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/orangejohny/SBWeb/pkg/model"
	sm "github.com/orangejohny/SBWeb/pkg/sessionmanager"
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

	SM.CreateSession(nil, true)
}

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
		t.Error(err)
	}

	sID, _ := SM.CreateSession(&model.Session{
		ID:        15,
		Login:     "aaa@eee.ru",
		UserAgent: "ieieie",
	}, true)

	_, err = SM.CheckSession(sID)
	if err != nil {
		t.Error("Key must exist")
	}

	s.FastForward(5 * time.Second)

	res, _ := SM.CheckSession(sID)
	if res != nil {
		t.Error("Key mustn't exist")
	}

	sID, _ = SM.CreateSession(&model.Session{
		ID:        15,
		Login:     "aaa@eee.ru",
		UserAgent: "ieieie",
	}, true)

	_, err = SM.CheckSession(sID)
	if err != nil {
		t.Error("Key must exist")
	}

	SM.DeleteSession(sID)

	res, _ = SM.CheckSession(sID)
	if res != nil {
		t.Error("Key mustn't exist")
	}

	sID, _ = SM.CreateSession(&model.Session{
		ID:        15,
		Login:     "aaa@eee.ru",
		UserAgent: "ieieie",
	}, false)

	_, err = SM.CheckSession(sID)
	if err != nil {
		t.Error("Key must exist")
	}

	s.FastForward(5 * time.Second)

	_, err = SM.CheckSession(sID)
	if err != nil {
		t.Error("Key must exist")
	}
}
