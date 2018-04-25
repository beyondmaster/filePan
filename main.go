package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"filePan/config"
	"filePan/model"
	"filePan/router"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	webServer := &http.Server{
		Addr:         ":" + fmt.Sprintf("%d", config.ServerConfig.Port),
		Handler:      webHandler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fileServer := &http.Server{
		Addr:         ":" + fmt.Sprintf("%d", config.ServerConfig.FilePort),
		Handler:      fileHandler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	g.Go(func() error {
		return webServer.ListenAndServe()
	})
	g.Go(func() error {
		return fileServer.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func webHandler() http.Handler {
	fmt.Println("gin.Version: ", gin.Version)
	if config.ServerConfig.Env != model.DevelopmentMode {
		// Disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()
		// Logging to a file.
		logFile, err := os.OpenFile(config.ServerConfig.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(-1)
		}
		gin.DefaultWriter = io.MultiWriter(logFile)
	}

	// Creates a router without any middleware by default
	app := gin.New()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	maxSize := int64(config.ServerConfig.MaxMultipartMemory)
	app.MaxMultipartMemory = maxSize << 20 // 3 MiB

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	app.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	app.Use(gin.Recovery())

	router.Route(app)

	return app
}

func fileHandler() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.Static("/", config.ServerConfig.FilePanDir)

	return e
}
