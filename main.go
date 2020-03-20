package main

import (
	"grpc-cms/cmd"
	"grpc-cms/models"
	"grpc-cms/pkg/setting"
)

func init() {
	setting.Setup()
	models.Setup()
}

func main() {
	cmd.Execute()
}
