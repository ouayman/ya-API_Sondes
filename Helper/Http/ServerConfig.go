package helperhttp

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type Timeout struct {
	Read  int `json:"read"`
	Write int `json:"write"`
	Idle  int `json:"idle"`
}

type TLS struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

type ServerConfig struct {
	Host       string  `json:"host"`
	Port       int     `json:"port"`
	TLS        *TLS    `json:"tls,omitempty"`
	TimeoutSec Timeout `json:"timeoutSec"`
}

func CreateService(config ServerConfig, router http.Handler) *http.Server {
	var tlsConfig *tls.Config
	if config.TLS != nil {
		tlsConfig = &tls.Config{
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			/*
				// Check compatibility with HTTP/2
				PreferServerCipherSuites: true,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				},
			*/
		}

		if tlsCertificate, err := tls.LoadX509KeyPair(config.TLS.Cert, config.TLS.Key); err == nil {
			tlsConfig.Certificates = []tls.Certificate{tlsCertificate}
		}
	}

	return &http.Server{
		Addr:      fmt.Sprintf("%v:%v", config.Host, config.Port),
		TLSConfig: tlsConfig,
		// Good practice to set timeouts to avoid Slowloris attacks.
		ReadTimeout:  time.Second * time.Duration(config.TimeoutSec.Read),
		WriteTimeout: time.Second * time.Duration(config.TimeoutSec.Write),
		IdleTimeout:  time.Second * time.Duration(config.TimeoutSec.Idle),
		Handler:      router,
	}
}
