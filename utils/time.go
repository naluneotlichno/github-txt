package utils

import (
	"fmt"
	"time"
	"io"	// для работы с потоками данных
)

// Timer - структура для удобного измерения времени
type Timer struct {
	start time.Time
}

// StartTimer - создает новый таймер
func StartTimer() Timer {
	return Timer{start: time.Now()}
}

// Elapsed - возвращает прошедшее время в удобном формате
func (t Timer) Elapsed() string {
	return fmt.Sprintf("%v", time.Since(t.start))
}

// PrintElapsedTime - выводит в лог скольо времени прошло
func (t Timer) PrintElapsedTime(task string, log io.Writer) {
	fmt.Fprintf(log, "⏳ Время выполнения %s: %s\n", task, t.Elapsed())
}

// runTimedAction - запускает действие с таймером
func RunTimedAction(action func() error, stepName string, log io.Writer, retry bool) {
	timer := StartTimer() 
	HandleErrorRetry(action, stepName, log, retry)
	timer.PrintElapsedTime("Время выполнения: "+stepName, log) 
}