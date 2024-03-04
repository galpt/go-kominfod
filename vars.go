package main

import (
	"crypto/tls"
	"net/http"
	"runtime"
	"time"

	"github.com/spf13/afero"
)

const (
	Gigabyte    = 1 << 30
	Megabyte    = 1 << 20
	Kilobyte    = 1 << 10
	timeoutTr   = 30 * time.Second
	hostPortGin = "7777"
	usrAgent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"
)

var (
	mem         runtime.MemStats
	HeapAlloc   string
	SysMem      string
	Frees       string
	NumGCMem    string
	timeElapsed string
	latestLog   string

	CertFilePath = "/etc/letsencrypt/live/net.0ms.dev/fullchain.pem"
	KeyFilePath  = "/etc/letsencrypt/live/net.0ms.dev/privkey.pem"

	tlsConf = &tls.Config{
		InsecureSkipVerify: true,
	}

	h1Tr = &http.Transport{
		DisableKeepAlives:      false,
		DisableCompression:     false,
		ForceAttemptHTTP2:      false,
		TLSClientConfig:        tlsConf,
		TLSHandshakeTimeout:    timeoutTr,
		ResponseHeaderTimeout:  timeoutTr,
		IdleConnTimeout:        timeoutTr,
		ExpectContinueTimeout:  1 * time.Second,
		MaxIdleConns:           1000,     // Prevents resource exhaustion
		MaxIdleConnsPerHost:    100,      // Increases performance and prevents resource exhaustion
		MaxConnsPerHost:        0,        // 0 for no limit
		MaxResponseHeaderBytes: 64 << 10, // 64k
		WriteBufferSize:        64 << 10, // 64k
		ReadBufferSize:         64 << 10, // 64k
	}

	h1Client = &http.Client{
		Transport: h1Tr,
		Timeout:   timeoutTr,
	}

	// untuk kominfod
	memFS       = afero.NewMemMapFs()
	kominfodDir = "domains.txt"
)
