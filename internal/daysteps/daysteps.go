package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("parsePackage: expected 2 fields (steps,duration), got %d", len(parts))
	}

	stepsStr := strings.TrimSpace(parts[0])
	durStr := strings.TrimSpace(parts[1])

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("parsePackage: invalid steps %q: %w", stepsStr, err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("parsePackage: steps must be > 0, got %d", steps)
	}

	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return 0, 0, fmt.Errorf("parsePackage: invalid duration %q: %w", durStr, err)
	}

	return steps, dur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, dur, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		log.Println("steps must be > 0")
		return ""
	}

	// Дистанция в км: шаги * длина шага (м) / метров в км
	distKm := (float64(steps) * stepLength) / float64(mInKm)

	// Калории по модели из пакета spentcalories
	cals, err := spentcalories.WalkingSpentCalories(steps, weight, height, dur)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distKm, cals,
	)
}
