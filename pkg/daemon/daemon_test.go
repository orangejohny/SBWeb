// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package daemon_test

import (
	"testing"

	"bmstu.codes/developers34/SBWeb/pkg/api"
	"bmstu.codes/developers34/SBWeb/pkg/daemon"
	"bmstu.codes/developers34/SBWeb/pkg/db"
	"bmstu.codes/developers34/SBWeb/pkg/sessionmanager"
)

func TestRunService(t *testing.T) {
	cfg := &daemon.Config{
		DB: db.Config{
			DBAddress:    "postgresql://runner:@postgres/data?sslmode=disable",
			MaxOpenConns: 10,
		},
		SM: sessionmanager.Config{
			DBAddress:      "redis://redis:6379/0",
			TockenLength:   32,
			ExpirationTime: 86400,
		},
		API: api.Config{
			Address:      ":54000",
			ReadTimeout:  "10s",
			WriteTimeout: "10s",
			IdleTimeout:  "10s",
		},
	}

	err := daemon.RunService(cfg)
	if err != nil {
		t.Error("Unexpected error: ", err.Error())
	}
}
