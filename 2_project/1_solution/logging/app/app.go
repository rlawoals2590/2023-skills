package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

func determineLogLevel(statusCode int) LogLevel {
	if statusCode >= 200 && statusCode < 300 {
		return INFO
	} else if statusCode >= 300 && statusCode < 400 {
		return WARN
	} else {
		return ERROR
	}
}

func main() {
	// 로그 파일 열기 또는 생성
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// 로그 파일에 기록할 로거 생성
	logger := log.New(logFile, "", 0)

	// HTTP 라우터 설정
	r := gin.Default()

	// 미들웨어를 사용하여 로그를 작성하는 핸들러 함수 추가
	r.Use(func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		protocol := c.Request.Proto
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()

		logLevel := determineLogLevel(statusCode)

		logEntry := fmt.Sprintf("%s - [%s] [%s] \"%s %s %s %d %.1f \"%s\"\"\n",
			clientIP, endTime.Format("2006-01-02T15:04:05Z"), logLevel,
			method, path, protocol, statusCode, float64(latency.Microseconds()), userAgent)

		// 로그를 파일과 표준 출력에 기록
		logger.Print(logEntry)
		fmt.Print(logEntry)
	})

	r.GET("/v1/color/red", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"color": "red"})
	})

	r.GET("/v1/color/orange", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"color": "orange"})
	})

	r.GET("/v1/color/melon", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"color": "melon"})
	})

	r.GET("/v1/error/5xx", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	})

	r.GET("/v1/error/4xx", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	})

	r.GET("v1/error/3xx", func(c *gin.Context) {
		c.JSON(http.StatusNotModified, gin.H{"message": "Not Modified"})
	})

	// 서버 시작
	r.Run(":8080")
}
