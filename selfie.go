package selfie

import (
	"time"

	"github.com/pressly/selfie/config"
	"github.com/pressly/selfie/data"
	"github.com/pressly/selfie/logme"
	"github.com/pressly/selfie/web"
	"github.com/pressly/selfie/web/security"
	"github.com/tylerb/graceful"
)

//selfie main structuire for selfie's app
type selfie struct {
	conf *config.Config
}

//Start starts selfie App, listeing to specified port
func (r *selfie) Start() {
	graceful.Run(r.conf.Server.Bind, 10*time.Second, web.New())
}

//Exit stops the app
func (r *selfie) Exit() {
	data.DB.Close()
}

//New makes a new and setup releasifer app's settings
func New(conf *config.Config) (*selfie, error) {
	logme.Info("selfie started at " + conf.Server.Bind)

	app := &selfie{conf: conf}

	//setup security
	security.Setup(conf)

	//Start a new DB session
	_, err := data.NewDBWithConfig(conf)
	if err != nil {
		logme.Fatal(err)
	}

	return app, nil
}
