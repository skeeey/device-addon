package device

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var temperatures = []int{15, 20, 22, 25, 30}

type TemperatureHandler func(int)

type Thermometer struct {
	interrupt chan struct{}
	handlers  []TemperatureHandler
}

func NewThermometer() *Thermometer {
	return &Thermometer{
		interrupt: make(chan struct{}),
		handlers:  []TemperatureHandler{},
	}
}

func (t *Thermometer) AddHandlers(handlers ...TemperatureHandler) *Thermometer {
	t.handlers = append(t.handlers, handlers...)
	return t
}

func (t *Thermometer) Stop() {
	t.interrupt <- struct{}{}
}

func (t *Thermometer) Run() {
	fmt.Fprintln(os.Stdout, "Starting ...")
	for {
		select {
		case <-t.interrupt:
			fmt.Fprintln(os.Stdout, "Shutdown")
			return
		default:
			for _, handler := range t.handlers {
				handler(temperatures[rand.Intn(len(temperatures)-1)])
			}
			time.Sleep(5 * time.Second)
		}
	}
}
