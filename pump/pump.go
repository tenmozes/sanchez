package pump

import (
	"fmt"
	"io"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"strconv"
	"sync/atomic"
	"time"
	"log"
)

type Worker interface {
	io.Closer
	Work(time.Duration) error
	Name() string
}

type BusyError struct {
	n int
}

func (be BusyError) Error() string {
	return fmt.Sprintf("pump %d is working hard at this moment", be.n)
}

type Pump struct {
	name string
	in   int
	pin  gpio.PinIO
	busy *int32
}

func NewPump(name string, in int) *Pump {
	p := gpioreg.ByName(strconv.Itoa(in))
	var busy int32
	return &Pump{name: name, in: in, pin: p, busy: &busy}
}

func (p *Pump) Start(t time.Duration) error {
	if !atomic.CompareAndSwapInt32(p.busy, 0, 1) {
		return BusyError{n: p.in}
	}
	if err := p.pin.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(t)
	atomic.StoreInt32(p.busy, 0)
	return nil
}

func (p *Pump) Stop() error {
	log.Printf("stop working on %d", p.in)
	if p.pin == nil {
		return fmt.Errorf("can't stop nil pin(%d) did you run on program on RPi", p.in)
	}
	return p.pin.Out(gpio.Low)
}

func (p *Pump) Work(t time.Duration) error {
	log.Printf("start working on %d", p.in)
	if err := p.Start(t); err != nil {
		return err
	}
	return p.Stop()
}

func (p *Pump) Name() string {
	return p.name
}

func (p *Pump) Close() error {
	return p.Stop()
}
