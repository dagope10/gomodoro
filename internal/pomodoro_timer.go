package internal

import (
	"fmt"
	"gomodoro/internal/config"
	"gomodoro/internal/timer"
	"os"
	"strconv"
	"time"
)

func Start() {
	cfg := config.Load("pomodoro_config.json")

	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case "start":
		startTimer(cfg)
	case "work", "short", "long", "cycle":
		updateConfig(cfg, os.Args[1])
	case "status":
		showStatus(cfg)
	default:
		fmt.Printf("Comando no reconocido: %s\n", os.Args[1])
		usage()
	}
}

func startTimer(cfg *config.Config) {
	timer := timer.New(cfg.Work, cfg.ShortBreak, cfg.LongBreak, cfg.Cycle)

	fmt.Println("🚀 Iniciando Pomodoro Timer...")
	fmt.Printf("⏱️  Configuración: Trabajo %v | Descanso corto %v | Descanso largo %v | Ciclos %d\n",
		cfg.Work, cfg.ShortBreak, cfg.LongBreak, cfg.Cycle)
	fmt.Println("💡 Presiona Ctrl+C para detener")

	// El timer ya maneja las señales internamente
	timer.Start()
}

func updateConfig(cfg *config.Config, field string) {
	if len(os.Args) < 3 {
		fmt.Printf("Uso: gomo %s <valor>\n", field)
		return
	}

	value := os.Args[2]

	switch field {
	case "work":
		duration, err := parseDuration(value)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		oldValue := cfg.Work
		cfg.Work = duration
		fmt.Printf("⏱️  Tiempo de trabajo actualizado: %v → %v\n", oldValue, cfg.Work)

	case "short":
		duration, err := parseDuration(value)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		oldValue := cfg.ShortBreak
		cfg.ShortBreak = duration
		fmt.Printf("☕ Descanso corto actualizado: %v → %v\n", oldValue, cfg.ShortBreak)

	case "long":
		duration, err := parseDuration(value)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		oldValue := cfg.LongBreak
		cfg.LongBreak = duration
		fmt.Printf("😴 Descanso largo actualizado: %v → %v\n", oldValue, cfg.LongBreak)

	case "cycle":
		cycles, err := strconv.Atoi(value)
		if err != nil || cycles < 1 {
			fmt.Println("Error: el número de ciclos debe ser un número entero mayor a 0")
			return
		}
		oldValue := cfg.Cycle
		cfg.Cycle = cycles
		fmt.Printf("🔄 Ciclos actualizados: %d → %d\n", oldValue, cfg.Cycle)
	}

	if err := cfg.Save("pomodoro_config.json"); err != nil {
		fmt.Printf("❌ Error al guardar configuración: %v\n", err)
		return
	}
	fmt.Println("✅ Configuración guardada")
}

func showStatus(cfg *config.Config) {
	fmt.Println("📊 Configuración actual:")
	fmt.Printf("   🍅 Trabajo: %v\n", cfg.Work)
	fmt.Printf("   ☕ Descanso corto: %v\n", cfg.ShortBreak)
	fmt.Printf("   😴 Descanso largo: %v\n", cfg.LongBreak)
	fmt.Printf("   🔄 Ciclos: %d\n", cfg.Cycle)
}

func parseDuration(value string) (time.Duration, error) {
	// Intentar parsear como duración (25m, 5m, etc.)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration, nil
	}

	// Si falla, intentar como minutos simples
	if minutes, err := strconv.Atoi(value); err == nil && minutes > 0 {
		return time.Duration(minutes) * time.Minute, nil
	}

	return 0, fmt.Errorf("formato inválido. Usa '25m' o '25' (minutos)")
}

func usage() {
	fmt.Println("🍅 Gomodoro - Pomodoro Timer")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  gomo start                    - Iniciar el timer")
	fmt.Println("  gomo status                   - Ver configuración actual")
	fmt.Println("  gomo work <tiempo>            - Configurar tiempo de trabajo")
	fmt.Println("  gomo short <tiempo>           - Configurar descanso corto")
	fmt.Println("  gomo long <tiempo>            - Configurar descanso largo")
	fmt.Println("  gomo cycle <número>           - Configurar número de ciclos")
	fmt.Println()
	fmt.Println("Ejemplos:")
	fmt.Println("  gomo work 25m                 - 25 minutos de trabajo")
	fmt.Println("  gomo work 25                  - 25 minutos (formato corto)")
	fmt.Println("  gomo short 5m                 - 5 minutos de descanso")
	fmt.Println("  gomo cycle 4                  - 4 ciclos antes del descanso largo")
}
