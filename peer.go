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
			globalJson        []Global
			domainsJson       []Domain
			query             string
			queryList         []string
			queryTotal        int = 0
			domainList        []string
			timestampMultiArr []float64
			timeAvgTotal      float64 = 0
		)
		query = c.DefaultQuery("domain", "")
		timestampTotal := time.Now()
		globalJson = nil
		domainsJson = nil
		queryList = nil
		domainList = nil
		timestampMultiArr = nil

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

			queryTotal = len(queryList)

			for qryIdx := range queryList {

				timestampMulti := time.Now()

				if contains(domainList, queryList[qryIdx]) {

					timeSinceMulti := time.Since(timestampMulti)
					timestampMultiArr = append(timestampMultiArr, float64(timeSinceMulti))

					if xstrings.Count(queryList[qryIdx], ".") > 1 {

						domainsJson = append(domainsJson, Domain{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(timeSinceMulti/time.Millisecond), float64(timeSinceMulti/time.Microsecond)), Domain: queryList[qryIdx], DomainIndex: chkIdx(domainList, queryList[qryIdx]), IsSubdomain: true, IsBlocked: true})

					} else {
						domainsJson = append(domainsJson, Domain{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(timeSinceMulti/time.Millisecond), float64(timeSinceMulti/time.Microsecond)), Domain: queryList[qryIdx], DomainIndex: chkIdx(domainList, queryList[qryIdx]), IsSubdomain: false, IsBlocked: true})

					}

				} else if !contains(domainList, queryList[qryIdx]) {

					timeSinceMulti := time.Since(timestampMulti)
					timestampMultiArr = append(timestampMultiArr, float64(timeSinceMulti))

					if xstrings.Count(queryList[qryIdx], ".") > 1 {
						domainsJson = append(domainsJson, Domain{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(timeSinceMulti/time.Millisecond), float64(timeSinceMulti/time.Microsecond)), Domain: queryList[qryIdx], DomainIndex: chkIdx(domainList, queryList[qryIdx]), IsSubdomain: true, IsBlocked: false})

					} else {
						domainsJson = append(domainsJson, Domain{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(timeSinceMulti/time.Millisecond), float64(timeSinceMulti/time.Microsecond)), Domain: queryList[qryIdx], DomainIndex: chkIdx(domainList, queryList[qryIdx]), IsSubdomain: false, IsBlocked: false})

					}

				}

			}
		} else {

			queryTotal = 1
			timestampSingle := time.Now()

			if contains(domainList, query) {

				timeSinceSingle := time.Since(timestampSingle)

				if xstrings.Count(query, ".") > 1 {
					domainsJson = append(domainsJson, Domain{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(timeSinceSingle/time.Millisecond), float64(timeSinceSingle/time.Microsecond)), Domain: query, DomainIndex: chkIdx(domainList, query), IsSubdomain: true, IsBlocked: true})

				} else {
					domainsJson = append(domainsJson, Domain{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(timeSinceSingle/time.Millisecond), float64(timeSinceSingle/time.Microsecond)), Domain: query, DomainIndex: chkIdx(domainList, query), IsSubdomain: false, IsBlocked: true})

				}
			} else if !contains(domainList, query) {

				timeSinceSingle := time.Since(timestampSingle)

				if xstrings.Count(query, ".") > 1 {
					domainsJson = append(domainsJson, Domain{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(timeSinceSingle/time.Millisecond), float64(timeSinceSingle/time.Microsecond)), Domain: query, DomainIndex: chkIdx(domainList, query), IsSubdomain: true, IsBlocked: false})

				} else {
					domainsJson = append(domainsJson, Domain{ExecTime: fmt.Sprintf("%.2f ms | %.2f μs", float64(timeSinceSingle/time.Millisecond), float64(timeSinceSingle/time.Microsecond)), Domain: query, DomainIndex: chkIdx(domainList, query), IsSubdomain: false, IsBlocked: false})

				}

			}

		}

		timeSinceTotal := time.Since(timestampTotal)

		for timestampIdx := range timestampMultiArr {
			timeAvgTotal = float64(timeAvgTotal + timestampMultiArr[timestampIdx])
		}
		timeAvgTotal = float64(timeAvgTotal) / float64(len(timestampMultiArr))
		timeAvgDuration := time.Duration(timeAvgTotal)

		globalJson = append(globalJson, Global{ExecTimeTotal: fmt.Sprintf("%.2f ms | %.2f μs", float64(timeSinceTotal/time.Millisecond), float64(timeSinceTotal/time.Microsecond)), ExecTimeAverage: fmt.Sprintf("%.2f ms | %.2f μs", float64(timeAvgDuration/time.Millisecond), float64(timeAvgDuration/time.Microsecond)), DomainQueryTotal: queryTotal, Domains: domainsJson})

		c.IndentedJSON(http.StatusOK, globalJson)

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
