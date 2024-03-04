package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huandu/xstrings"
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

		var (
			resultJson []JsonOutput
			query      string
			queryList  []string
			domainList []string
		)
		query = c.DefaultQuery("domain", "")
		queryList = nil
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

				timestampMulti := time.Now()

				if contains(domainList, queryList[qryIdx]) {

					if xstrings.Count(queryList[qryIdx], ".") > 1 {
						resultJson = append(resultJson, JsonOutput{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(time.Since(timestampMulti)/time.Millisecond), float64(time.Since(timestampMulti)/time.Microsecond)), Domain: queryList[qryIdx], DomainIndex: chkIdx(domainList, queryList[qryIdx]), IsSubdomain: true, IsBlocked: true})

					} else {
						resultJson = append(resultJson, JsonOutput{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(time.Since(timestampMulti)/time.Millisecond), float64(time.Since(timestampMulti)/time.Microsecond)), Domain: queryList[qryIdx], DomainIndex: chkIdx(domainList, queryList[qryIdx]), IsSubdomain: false, IsBlocked: true})

					}
				} else if !contains(domainList, queryList[qryIdx]) {

					if xstrings.Count(queryList[qryIdx], ".") > 1 {
						resultJson = append(resultJson, JsonOutput{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(time.Since(timestampMulti)/time.Millisecond), float64(time.Since(timestampMulti)/time.Microsecond)), Domain: queryList[qryIdx], DomainIndex: chkIdx(domainList, queryList[qryIdx]), IsSubdomain: true, IsBlocked: false})

					} else {
						resultJson = append(resultJson, JsonOutput{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(time.Since(timestampMulti)/time.Millisecond), float64(time.Since(timestampMulti)/time.Microsecond)), Domain: queryList[qryIdx], DomainIndex: chkIdx(domainList, queryList[qryIdx]), IsSubdomain: false, IsBlocked: false})

					}

				}

			}
		} else {

			timestampSingle := time.Now()

			if contains(domainList, query) {

				if xstrings.Count(query, ".") > 1 {
					resultJson = append(resultJson, JsonOutput{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(time.Since(timestampSingle)/time.Millisecond), float64(time.Since(timestampSingle)/time.Microsecond)), Domain: query, DomainIndex: chkIdx(domainList, query), IsSubdomain: true, IsBlocked: true})

				} else {
					resultJson = append(resultJson, JsonOutput{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(time.Since(timestampSingle)/time.Millisecond), float64(time.Since(timestampSingle)/time.Microsecond)), Domain: query, DomainIndex: chkIdx(domainList, query), IsSubdomain: false, IsBlocked: true})

				}
			} else if !contains(domainList, query) {

				if xstrings.Count(query, ".") > 1 {
					resultJson = append(resultJson, JsonOutput{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(time.Since(timestampSingle)/time.Millisecond), float64(time.Since(timestampSingle)/time.Microsecond)), Domain: query, DomainIndex: chkIdx(domainList, query), IsSubdomain: true, IsBlocked: false})

				} else {
					resultJson = append(resultJson, JsonOutput{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(time.Since(timestampSingle)/time.Millisecond), float64(time.Since(timestampSingle)/time.Microsecond)), Domain: query, DomainIndex: chkIdx(domainList, query), IsSubdomain: false, IsBlocked: false})

				}

			}

		}

		c.IndentedJSON(http.StatusOK, resultJson)

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
