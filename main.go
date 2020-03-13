package main

import (
	"grpc-cms/cmd"
	"grpc-cms/models"
)

func init() {
	models.Setup()
}

func main() {
	//server.Port = "50052"
	//server.CertPemPath = "./certs/cert.pem"
	//server.KeyPemPath = "./certs/key.pem"
	//server.CertName = "localhost"
	cmd.Execute()
	//_ = server.Run()
}
