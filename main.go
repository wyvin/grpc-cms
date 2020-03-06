package main

import (
	"grpc-content/cmd"
	"grpc-content/models"
)

func init() {
	models.Setup()
}

func main() {
	cmd.Execute()
}
