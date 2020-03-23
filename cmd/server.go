package cmd

import (
	"github.com/spf13/cobra"
	"grpc-cms/server"
	"log"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC content server",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recover error : %v", err)
			}
		}()

		_ = server.Run()
	},
}

func init() {
	serverCmd.Flags().StringVarP(&server.CertPemPath, "cert-pem", "", "./certs/cert.pem", "cert pem path")
	serverCmd.Flags().StringVarP(&server.KeyPemPath, "key-pem", "", "./certs/key.pem", "key pem path")
	rootCmd.AddCommand(serverCmd)
}
