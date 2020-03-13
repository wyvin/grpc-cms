package main

import (
	"grpc-content/cmd"
	"grpc-content/models"
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
