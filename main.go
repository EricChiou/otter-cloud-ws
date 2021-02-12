package main

import (
	"log"

	"otter-cloud-ws/acl"
	"otter-cloud-ws/config"
	"otter-cloud-ws/db/mysql"
	"otter-cloud-ws/jobqueue"
	"otter-cloud-ws/minio"
	"otter-cloud-ws/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// load config
	if err := config.Load(config.ConfigPath); err != nil {
		panic(err)
	}
	cfg := config.Get()

	// init db
	if err := mysql.Init(); err != nil {
		panic(err)
	}
	defer mysql.Close()

	// init minio
	if err := minio.Init(); err != nil {
		panic(err)
	}

	// init jobqueue
	jobqueue.Init()

	// load acl
	if err := acl.Load(); err != nil {
		panic(err)
	}

	// set headers
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS")
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// init api
	router.Init()

	// start http server
	// if err := router.ListenAndServe(cfg.ServerPort); err != nil {
	// 	panic(err)
	// }
	// start https server
	if err := router.ListenAndServeTLS(cfg.ServerPort, cfg.SSLCertFilePath, cfg.SSLKeyFilePath); err != nil {
		panic(err)
	}

	// waiting for jobqueue finished
	jobqueue.Wait()

	defer func() {
		if err := recover(); err != nil {
			log.Println("start server error:", err)
		}
	}()
}
