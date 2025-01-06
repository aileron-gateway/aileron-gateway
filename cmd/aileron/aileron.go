package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aileron-gateway/aileron-gateway/app"
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

	core := api.NewFactoryAPI()
	register.RegisterAll(core) // Register all core APIs.

	_ = svr.Handle("core/", core) // Handle "core/*" APIs.

	if err := a.Run(svr); err != nil {
		e := app.ErrAppMain.WithStack(err, nil)
		fmt.Printf("%s\n\n%s\n", e.Error(), e.StackTrace())
		app.Exit(1)
	}
}
