package lib

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	logger *log.Logger
	file   *os.File
}

func StartLogger() *Logger {

	dir := fmt.Sprintf("./log/%s/", time.Now().Format("2006-01-02"))

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0644); err != nil {
			log.Fatalf("Failed to create logs directory: %v", err)
		}
	}
	file, err := os.OpenFile(dir+"info.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Fatalln("Failed to create log file:", err)
	}

	// 创建一个新的日志记录器
	return &Logger{
		logger: log.New(file, "", log.Ldate|log.Ltime),
		file:   file,
	}

}

func (l Logger) Error(v ...interface{}) {
	l.logger.Println("< error >", v)
}

func (l Logger) Info(v ...interface{}) {
	l.logger.Println("< info >", v)
}
func (l Logger) StopLogger() {
	l.file.Close()
}
