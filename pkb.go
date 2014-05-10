package main

import (
	"github.com/hagna/dumbkb/rawkb"
	"log"
	"time"
)

type Pad struct {
	id          string
	boardstate  uint8
	keyspressed uint8
}

func (p *Pad) sendphone() {
	log.Println(p.keyspressed)
}

func (p *Pad) keyup(keyno uint8) {
	b := p.boardstate
	b = b & ^(1 << keyno)
	p.boardstate = b
	if b == 0 {
		p.sendphone()
		p.keyspressed = 0
	}
}

func (p *Pad) keydown(keyno uint8) {
	log.Println(p.id, keyno)
	b := p.boardstate
	b = b | (1 << keyno)
	p.boardstate = b
	p.keyspressed = p.keyspressed | (1 << keyno)
}

type Mcs struct {
	keys     map[uint16]uint8
	watched  map[uint16]func()
	states   map[string]func(k uint16) string
	kbevents chan uint16
	pad      [2]Pad
}

func (m *Mcs) state_log(k uint16) string {
	if callable, ok := m.watched[k]; ok {
		callable()
	}
	return "loggy"
}

func iskeyup(k uint16) bool {
	if k > 0x80 {
		return true
	}
	return false
}

func (mcs *Mcs) state_learn(k uint16) string {
	b := k
	log.Println("learning the keys starting with", k)
	for iskeyup(b) {
		log.Println("learn: discarding keyup", b)
		b = <-mcs.kbevents
	}
	for j := 0; j < 2; j++ {
		var p Pad = mcs.pad[j]
		for i := 0; i < 6; i++ {
			var j int = i
			mcs.watched[b] = func() {
				p.keydown(uint8(j))
			}
			b = <-mcs.kbevents
			mcs.watched[b] = func() {
				p.keyup(uint8(j))
			}
			log.Println("learned key", i)
			b = <-mcs.kbevents
		}
	}
	return "loggy"
}

func NewMcs() *Mcs {
	m := new(Mcs)
	m.states = make(map[string]func(k uint16) string)
	m.keys = make(map[uint16]uint8, 12)
	m.kbevents = make(chan uint16, 2)
	m.watched = make(map[uint16]func())
	m.states["loggy"] = m.state_log
	m.states["init"] = m.state_learn
	m.pad = [2]Pad{{"left", 0, 0}, {"right", 0, 0}}
	log.Println(m.pad)
	return m
}

func interact() {
	/*	logfile, err := os.Create("log")
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logfile)

	*/
	if rawkb.SetupKeyboard() == 0 {
		log.Println("must be on the console for raw keyboard access")
		return
	} else {
		log.Println("setting up rawkb")
		defer rawkb.RestoreKeyboard()
	}
	mcs := NewMcs()
	log.Println(mcs)
	log.Println(mcs.kbevents)
	go func() {
		mcs.kbevents <- 0x81
		for {
			b, ok := rawkb.ReadOnce()
			if ok != 255 {
				mcs.kbevents <- b

			}
			time.Sleep(1 * time.Microsecond)
		}
	}()

	m := mcs.states["init"]
inf:
	for {
		select {
		case b := <-mcs.kbevents:
			next_state := m(b)
			next_method, ok := mcs.states[next_state]
			if !ok {
				log.Println("no such state", next_state)
			} else {
				m = next_method
			}
			if b == 1 {
				break inf
			}
		}
	}

}

func main() {
	interact()
}
