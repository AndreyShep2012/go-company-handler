package integration

import (
	"testing"
	"time"

	"github.com/AndreyShep2012/go-company-handler/internal/app"
	"github.com/AndreyShep2012/go-company-handler/internal/config"
)

var testConf config.Config

func TestMain(m *testing.M) {
	var err error
	testConf, err = config.Load("config.yml")
	if err != nil {
		panic(err)
	}

	go func() {
		app.Serve(testConf)
	}()

	time.Sleep(3 * time.Second) // let the server start

	m.Run()
}
