package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("parseTraining: expected 3 fields (steps, activity, duration), got %d", len(parts))
	}

	stepsStr := strings.TrimSpace(parts[0])
	activity := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("parseTraining: invalid steps %q: %w", stepsStr, err)
	}

	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("parseTraining: invalid duration %q: %w", durationStr, err)
	}

	return steps, activity, dur, nil
}

func distance(steps int, height float64) float64 {
	// длина шага в метрах
	stepLen := height * stepLengthCoefficient
	// пройденная дистанция в километрах
	return float64(steps) * stepLen / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distKm := distance(steps, height) // км
	hours := duration.Hours()         // часы (float64)
	return distKm / hours             // км/ч
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, dur, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distKm := distance(steps, height)
	speed := meanSpeed(steps, height, dur)
	hours := dur.Hours()

	switch activity {
	case "Бег":
		cals, err := RunningSpentCalories(steps, weight, height, dur)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
			activity, hours, distKm, speed, cals,
		), nil

	case "Ходьба":
		cals, err := WalkingSpentCalories(steps, weight, height, dur)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
			activity, hours, distKm, speed, cals,
		), nil

	default:
		err := fmt.Errorf("неизвестный тип тренировки: %s", activity)
		log.Println(err)
		return "", err
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid steps: %d (must be > 0)", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight: %.2f (must be > 0)", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("invalid height: %.2f (must be > 0)", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration: %v (must be > 0)", duration)
	}

	speed := meanSpeed(steps, height, duration) // км/ч
	minutes := duration.Minutes()               // мин

	calories := (weight * speed * minutes) / float64(minInH)
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid steps: %d (must be > 0)", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight: %.2f (must be > 0)", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("invalid height: %.2f (must be > 0)", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration: %v (must be > 0)", duration)
	}

	speed := meanSpeed(steps, height, duration) // км/ч
	minutes := duration.Minutes()               // мин

	calories := ((weight * speed * minutes) / float64(minInH)) * walkingCaloriesCoefficient
	return calories, nil
}
