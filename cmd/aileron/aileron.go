package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/util/register"
)

// version is the version of the binary.
// This is intended to be set by -ldflags on build like
// go build -ldflags "-X main.version=v1.2.3"
var version = "UNSET"

func init() {
	app.Version = version // Overwrite app version.
}

var svr = api.NewDefaultServeMux()

func main() {
	// Modify standard log output flag.
	// Debug outputs may use the standard logger.
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	a := app.New()
	a.ParseArgs(os.Args[1:])

	f := api.NewFactoryAPI()
	register.RegisterAll(f) // Register all APIs.

	_ = svr.Handle("core/", f) // Handle "core/*" APIs.
	_ = svr.Handle("app/", f)  // Handle "app/*" APIs.

	if err := a.Run(svr); err != nil {
		e := app.ErrAppMain.WithStack(err, nil)
		fmt.Printf("%s\n\n%s\n", e.Error(), e.StackTrace())
		app.Exit(1)
	}
}
