package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
)

func peer() {

	duration := time.Now()

	// Use Gin as the HTTP router
	gin.SetMode(gin.ReleaseMode)
	recover := gin.New()
	recover.Use(gin.Recovery())
	ginroute := recover

	// Custom NotFound handler
	ginroute.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, fmt.Sprintln("[404] NOT FOUND"))
	})

	// Print homepage
	ginroute.GET("/", func(c *gin.Context) {
		runtime.ReadMemStats(&mem)
		NumGCMem = fmt.Sprintf("%v", mem.NumGC)
		timeElapsed = fmt.Sprintf("%v", time.Since(duration))

		latestLog = fmt.Sprintf("\n •===========================• \n • [SERVER STATUS] \n • Last Modified: %v \n • Completed GC Cycles: %v \n • Time Elapsed: %v \n •===========================• \n\n", time.Now().UTC().Format(time.RFC850), NumGCMem, timeElapsed)

		c.String(http.StatusOK, fmt.Sprintf("%v", latestLog))
	})

	// untuk kominfod API
	ginroute.GET("/kominfod", func(c *gin.Context) {

		timestamp := time.Now()
		var (
			resultPrint      string
			query            string
			queryList        []string
			queryListResult  []string
			domainList       []string
			domainBlocked    = false
			domainBlockedIdx = 0
		)
		query = c.DefaultQuery("domain", "google.com")
		queryList = nil
		queryListResult = nil
		domainList = nil

		// read the cached file
		kominfodReadFile, err := afero.ReadFile(memFS, kominfodDir)
		if err != nil {
			fmt.Println(" [kominfodReadFile] ", err)
			return
		}

		domainList, err = strSplit(string(kominfodReadFile[:]), "\n")
		if err != nil {
			fmt.Println(" [domainList] ", err)
			return
		}

		if strings.Contains(query, ",") {

			queryList, err = strSplit(query, ",")
			if err != nil {
				fmt.Println(" [queryList] ", err)
				return
			}

			for qryIdx := range queryList {

				for domainIdx := range domainList {

					if domainList[domainIdx] == queryList[qryIdx] {
						domainBlocked = true
						domainBlockedIdx = (domainIdx + 1)
						break
					} else {
						domainBlocked = false
						domainBlockedIdx = 0
					}
				}

				if domainBlocked {
					queryListResult = append(queryListResult, fmt.Sprintf("%v >> BLOCKED [Index: %v]", queryList[qryIdx], domainBlockedIdx))
				} else {
					queryListResult = append(queryListResult, fmt.Sprintf("%v >> NOT BLOCKED", queryList[qryIdx]))
				}

			}
		} else {
			for domainIdx := range domainList {

				if domainList[domainIdx] == query {
					domainBlocked = true
					domainBlockedIdx = (domainIdx + 1)
					break
				} else {
					domainBlocked = false
					domainBlockedIdx = 0
				}

			}

			if domainBlocked {
				queryListResult = append(queryListResult, fmt.Sprintf("%v >> BLOCKED [Index: %v]", query, domainBlockedIdx))
			} else {
				queryListResult = append(queryListResult, fmt.Sprintf("%v >> NOT BLOCKED", query))
			}

		}

		resultPrint = fmt.Sprintf("Selesai dalam %.2f ms | %.2f μs \n\n%v", float64(time.Since(timestamp)/time.Millisecond), float64(time.Since(timestamp)/time.Microsecond), strings.Join(queryListResult, "\n"))

		c.String(http.StatusOK, resultPrint)

	})

	tlsConf = &tls.Config{
		InsecureSkipVerify: true,
		// Certificates:       []tls.Certificate{serverTLSCert},
	}

	// HTTP proxy server Gin
	httpserverGin := &http.Server{
		Addr:              fmt.Sprintf(":%v", hostPortGin),
		Handler:           ginroute,
		TLSConfig:         tlsConf,
		MaxHeaderBytes:    64 << 10, // 64k
		ReadTimeout:       timeoutTr,
		ReadHeaderTimeout: timeoutTr,
		WriteTimeout:      timeoutTr,
		IdleTimeout:       timeoutTr,
	}
	httpserverGin.SetKeepAlivesEnabled(true)

	notifyGin := fmt.Sprintf("[go-kominfod] Server is running on %v", fmt.Sprintf(":%v", hostPortGin))

	fmt.Println()
	fmt.Println(notifyGin)
	fmt.Println()
	// httpserverGin.ListenAndServe()
	httpserverGin.ListenAndServeTLS(CertFilePath, KeyFilePath)

}
