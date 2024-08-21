package utils

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
)

// Configure returns an elasticsearch.Config
func Configure() elasticsearch.Config {
	var conf elasticsearch.Config

	// Get elasticsearch connection settings
	address := viper.GetStringSlice("elasticsearch.address")
	if len(address) != 0 {
		conf.Addresses = address
	}
	viper.SetDefault("core.timeout", 10)
	timeout := viper.GetInt("core.timeout")
	t := time.Duration(timeout) * time.Second
	username := viper.GetString("elasticsearch.username")
	if username != "" {
		conf.Username = username
	}
	password := viper.GetString("elasticsearch.password")
	if password != "" {
		conf.Password = password
	}
	caCert := viper.GetString("elasticsearch.ca_cert_path")
	if caCert != "" {
		cert, err := os.ReadFile(caCert)
		if err != nil {
			panic(err)
		}
		conf.CACert = cert
	}

	conf.Transport = &http.Transport{
		ResponseHeaderTimeout: t,
		MaxIdleConnsPerHost:   10,
		DialContext:           (&net.Dialer{Timeout: t}).DialContext,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	return conf
}
