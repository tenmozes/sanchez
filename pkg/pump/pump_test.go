package pump

import (
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"time"
	"strconv"
	"sync/atomic"
	"fmt"
)

type BusyError struct {
	n int
}

func (be BusyError) Error() string{
	return fmt.Sprintf("pump %d is working hard at this moment", be.n)
}

type Pump struct {
	Name string
	in  int
	pin gpio.PinIO
	busy *int32
}

func NewPump(name string, in int) (*Pump, error) {
	p := gpioreg.ByName(strconv.Itoa(in))
	var busy int32
	return &Pump{Name: name, in: in, pin: p, busy:&busy}, nil
}

func (p *Pump) Start(t time.Duration) error {
	if !atomic.CompareAndSwapInt32(p.busy, 0, 1) {
		return BusyError{n:p.in}
	}
	if err := p.pin.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(t)
	atomic.StoreInt32(p.busy, 0)
	return nil
}

func (p *Pump) Stop() error {
	return p.pin.Out(gpio.Low)
}

func (p *Pump) Work(t time.Duration) error {
	if err := p.Start(t); err != nil {
		return err
	}
	return p.Stop()
}