package timer

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Timer struct {
	Work       time.Duration
	ShortBreak time.Duration
	LongBreak  time.Duration
	Cycle      int
	current    int
	stop       chan bool
}

func New(work, short, long time.Duration, cycle int) *Timer {
	return &Timer{
		Work:       work,
		ShortBreak: short,
		LongBreak:  long,
		Cycle:      cycle,
		current:    0,
		stop:       make(chan bool, 1),
	}
}

func (t *Timer) Start() {
	// Manejar Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n⏹️  Timer interrumpido")
		t.stop <- true
	}()

	for {
		select {
		case <-t.stop:
			return
		default:
		}

		fmt.Printf("🍅 Pomodoro %d: trabajando %v...\n", t.current+1, t.Work)
		if !t.runPhase(t.Work) {
			return // Se interrumpió
		}

		t.current++

		if t.current >= t.Cycle {
			fmt.Printf("😴 Descanso largo %v...\n", t.LongBreak)
			if !t.runPhase(t.LongBreak) {
				return
			}
			t.current = 0
			fmt.Println("🎉 ¡Ciclo completo! Comenzando nuevo ciclo...")
		} else {
			fmt.Printf("☕ Descanso corto %v...\n", t.ShortBreak)
			if !t.runPhase(t.ShortBreak) {
				return
			}
		}
	}
}

func (t *Timer) runPhase(d time.Duration) bool {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	remaining := int(d.Seconds())

	for remaining > 0 {
		select {
		case <-ticker.C:
			remaining--
			minutes := remaining / 60
			seconds := remaining % 60
			fmt.Printf("\r⏱️  Tiempo restante: %02d:%02d", minutes, seconds)
		case <-t.stop:
			fmt.Println()
			return false
		}
	}

	fmt.Println("\n🔔 ¡Tiempo completado!")
	return true
}

func (t *Timer) Stop() {
	select {
	case t.stop <- true:
	default:
	}
}
