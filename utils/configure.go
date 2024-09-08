package utils

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
)

// Configure returns an elasticsearch.Config
func Configure() elasticsearch.Config {
	Debug("Attempting to configure elasticsearch connection")
	var conf elasticsearch.Config

	// Get elasticsearch connection settings
	address := viper.GetStringSlice("elasticsearch.address")
	if len(address) != 0 {
		conf.Addresses = address
	}
	timeout := viper.GetInt("core.timeout")
	t := time.Duration(timeout) * time.Second

	certPath := viper.GetString("elasticsearch.ca_cert_path")
	cloudID := viper.GetString("elasticsearch.cloud_id")
	cloudAPIKey := viper.GetString("elasticsearch.cloud_api_key")
	fingerprint := viper.GetString("elasticsearch.certificate_fingerprint")
	username := viper.GetString("elasticsearch.username")
	password := viper.GetString("elasticsearch.password")

	// Try Elastic cloud
	if cloudID != "" && cloudAPIKey != "" {
		Debug("Using Elastic Cloud configuration")
		conf.CloudID = cloudID
		conf.APIKey = cloudAPIKey
	} else if certPath != "" && username != "" && password != "" {
		// Try with HTTPS certificate
		Debug("Using HTTPS certificate configuration")
		cert, _ := os.ReadFile(certPath)
		conf.CACert = cert
		conf.Username = username
		conf.Password = password

	} else if fingerprint != "" && username != "" && password != "" {
		// Try with HTTPS certificate fingerprint
		Debug("Using HTTPS certificate fingerprint configuration")
		conf.CertificateFingerprint = fingerprint
		conf.Username = username
		conf.Password = password

	} else if username != "" && password != "" {
		// Try Basic authentication
		Debug("Using Basic authentication configuration")
		conf.Username = username
		conf.Password = password
	} else {
		Fatal("No authentication method correctly configured", errors.New("Configuration error"))
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
