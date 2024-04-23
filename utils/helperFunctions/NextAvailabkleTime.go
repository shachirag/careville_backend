package helper

import (
	"careville_backend/entity"
	"time"
)

func GetLastTimeAvailable(slots []entity.Slots) string {
	var lastEndTime string
	for _, slot := range slots {
		for _, breakingSlot := range slot.BreakingSlots {
			endTime := breakingSlot.EndTime
			if endTime > lastEndTime {
				lastEndTime = endTime
			}
		}
	}
	return lastEndTime
}

func HasBreakingSlots(slots []entity.Slots) bool {
	for _, slot := range slots {
		for _, breakingSlot := range slot.BreakingSlots {
			startTime, _ := time.Parse("15:04", breakingSlot.StartTime)
			endTime, _ := time.Parse("15:04", breakingSlot.EndTime)
			currentTime := time.Now().UTC()

			if currentTime.After(startTime) && currentTime.Before(endTime) {
				return true
			}
		}
	}
	return false
}

func ContainsDay(days []string, target string) bool {
	for _, day := range days {
		if day == target {
			return true
		}
	}
	return false
}

func DayAfterCurrentDay(day string, currentTime time.Time) bool {
	currentWeekday := currentTime.Weekday().String()
	if day == currentWeekday {
		return false
	}
	daysMap := map[string]int{
		"Sunday":    0,
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
	}
	currentDayNum := daysMap[currentWeekday]
	targetDayNum := daysMap[day]

	return currentDayNum < targetDayNum
}

func GetUpcomingStartAndLastTime(slots []entity.Slots) string {
	currentTime := time.Now().UTC()
	for _, slot := range slots {
		if ContainsDay(slot.Days, currentTime.Weekday().String()) && DayAfterCurrentDay(slot.Days[0], currentTime) {
			continue
		}
		for _, day := range slot.Days {
			if DayAfterCurrentDay(day, currentTime) {
				startTime := slot.StartTime
				// lastTime := getLastTimeAvailable(slots)
				return startTime
			}
		}
	}
	return ""
}

// func GetLastTimeAvailable(slots []entity.Slots) string {
// 	var lastEndTime string
// 	for _, slot := range slots {
// 		for _, breakingSlot := range slot.BreakingSlots {
// 			endTime := breakingSlot.EndTime
// 			if endTime > lastEndTime {
// 				lastEndTime = endTime
// 			}
// 		}
// 	}
// 	return lastEndTime
// }
