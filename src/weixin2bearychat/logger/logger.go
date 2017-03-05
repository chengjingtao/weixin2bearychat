package logger

import "fmt"
import "time"

type Logger struct {
	id string
}

func New() *Logger {
	return &Logger{
		id: fmt.Sprint(time.Now().Nanosecond()),
	}
}

func (log *Logger) Debug(msg ...interface{}) {
	fmt.Printf("[%s] [DEBUG] [%s] %s \n", log.id, time.Now().Format("01-02 15:04:05"), fmt.Sprint(msg...))
}
func (log *Logger) Info(msg ...interface{}) {

	fmt.Printf("[%s] [INFO]  [%s] %s \n", log.id, time.Now().Format("01-02 15:04:05"), fmt.Sprint(msg...))
}
func (log *Logger) Warn(msg ...interface{}) {
	fmt.Printf("[%s] [WARN] [%s] %s \n", log.id, time.Now().Format("01-02 15:04:05"), fmt.Sprint(msg...))
}
func (log *Logger) Error(msg ...interface{}) {
	fmt.Printf("[%s] [ERROR] [%s] %s \n", log.id, time.Now().Format("01-02 15:04:05"), fmt.Sprint(msg...))
}
