package util

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"io/ioutil"
	"log"
)

func GetTLSConfig(certPemPath, keyPemPath string) *tls.Config {
	var certKeyPair *tls.Certificate
	cert, _ := ioutil.ReadFile(certPemPath)
	key, _ := ioutil.ReadFile(keyPemPath)

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Printf("TLS KeyPair err: %v\n", err)
	}

	certKeyPair = &pair

	return &tls.Config{
		Certificates: []tls.Certificate{*certKeyPair},
		NextProtos:   []string{http2.NextProtoTLS},
	}
}
