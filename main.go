// Copyright 2021 Ilia Frenkel. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	version = `¯\_(ツ)_/¯`
)

// options describes command line flags and options
type options struct {
	SiteRoot string `short:"r" long:"site-root" description:"root folder of the static website to serve"`
	Port     int    `short:"p" long:"port" default:"8080" description:"web server port to listen on"`
	Version  bool   `short:"v" long:"version" default:"false" description:"show version info and exit"`
}

var opts options

// printVersion prints out version, license and contact information.
func printVersion() {
	fmt.Println("static-with-password", version)
	fmt.Println("Copyright (c) 2021 Ilia Frenkel")
	fmt.Println("MIT License <https://opensource.org/licenses/MIT>")
	fmt.Println("Source code <https://github.com/iliafrenkel/static-with-password/>")
	fmt.Println("\nWritten by Ilia Frenkel<frenkel.ilia@gmail.com>")
	fmt.Println()
}

// Declare and parse command line flags.
func init() {
	flag.StringVar(&opts.SiteRoot, "r", "", "root folder of the static website to serve")
	flag.IntVar(&opts.Port, "p", 8080, "web server port to listen on")
	flag.BoolVar(&opts.Version, "v", false, "show version info and exit")
	flag.Parse()

	if opts.Version {
		printVersion()
		os.Exit(0)
	}
}

// authMiddleware authenticates the request and calls the next handler if
// successful. Returns 403 if authentication fails.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth middleware called.")
		// TODO: Do the actual authentication here.
		next.ServeHTTP(w, r)
	})
}

func main() {
	// quit channel listens on keyboard interrupts - SIGINT and SIGTERM
	// so that we can gracefully shutdown our server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// errc channel will receive an error if our server fails to start
	errc := make(chan error, 1)

	hdlr := authMiddleware(http.FileServer(http.Dir(opts.SiteRoot)))

	srv := &http.Server{
		Addr:         "localhost:" + fmt.Sprintf("%d", opts.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      hdlr,
	}

	go func() {
		fmt.Println("Server is listening on ", srv.Addr)
		errc <- srv.ListenAndServe()
	}()

	// Wait indefinitely for either one of the OS signals (SIGTERM or SIGINT)
	// or for one of the servers to return an error.
	select {
	case <-quit:
		fmt.Println("Shutting down ...")
	case err := <-errc:
		fmt.Printf("ERROR Startup failed, exiting: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Web server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("Web server is down.")
	}

	fmt.Println("Sayōnara!")
}
