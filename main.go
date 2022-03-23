package main

import (
	"math"
	"os"
	"strconv"

	"github.com/JoelD7/auto-scheduler/source"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type UserSource struct {
	Mettings    [][]string
	DailyBounds []string
}

type UserSourceArray []UserSource

const (
	maxMeetingDuration = 120
	sourcePath         = "samples/source.json"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("duration argument missing")
		return
	}

	durationAsInt, err := strconv.ParseInt(os.Args[1], 10, 32)
	if err != nil {
		fmt.Println("Please introduce a valid duration number.")
		return
	}

	if durationAsInt > maxMeetingDuration {
		fmt.Println("Sorry. The maximum meeting duration must be two hours.")
		return
	}

	duration := time.Minute * time.Duration(durationAsInt)

	if duration == 0 {
		fmt.Println("duration must be greater than zero")
		return
	}

	data, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	usersSource := UserSourceArray{}
	err = json.Unmarshal(data, &usersSource)
	if err != nil {
		fmt.Println(fmt.Errorf("error reading users source file: %w", err))
		return
	}

	users := make([]source.User, 0, len(usersSource))
	for _, userSource := range usersSource {
		newUser, err := source.NewUser(userSource.Mettings, userSource.DailyBounds)
		if err != nil {
			fmt.Println(fmt.Errorf("error creating users from source file: %w", err))
			return
		}

		users = append(users, *newUser)
	}

	possibleMeetingSlots := getSuggestedMeetingSlots(duration, users...)

	for _, slot := range possibleMeetingSlots {
		fmt.Println(slot.String())
	}

	if len(possibleMeetingSlots) == 0 {
		fmt.Println("No possible meeting slots between users.")
	}
}

func getSuggestedMeetingSlots(duration time.Duration, users ...source.User) []source.Frame {
	if len(users) == 0 {
		return nil
	}

	if len(users) < 2 {
		return users[0].FreeSpace
	}

	candidateSlots := make([]source.Frame, 0)
	suggestedSlots := make([]source.Frame, 0)
	var overlaps []source.Frame

	for i := 0; i < len(users)-1; i++ {
		cur := users[i]
		next := users[i+1]

		overlaps = getOverlapBetween(duration, cur, next)
		if len(overlaps) == 0 {
			continue
		}

		splittedFrames := splitFramesInDuration(overlaps, duration)

		candidateSlots = append(candidateSlots, splittedFrames...)

		if len(users) == 2 {
			return candidateSlots
		}
	}

	type frameFrequency struct {
		frequency int
		frame     source.Frame
	}

	freqMap := make(map[string]frameFrequency)
	for _, frame := range candidateSlots {
		freqMap[frame.String()] = frameFrequency{
			frequency: freqMap[frame.String()].frequency + 1,
			frame:     frame,
		}
	}

	for _, frameFrequency := range freqMap {
		if frameFrequency.frequency == len(users)-1 {
			suggestedSlots = append(suggestedSlots, frameFrequency.frame)
		}
	}

	return suggestedSlots
}

// getOverlapBetween returns shared free times between two users
func getOverlapBetween(duration time.Duration, cur, next source.User) []source.Frame {
	longestFrameList := cur.FreeSpace
	shortestFrameList := next.FreeSpace

	if len(cur.FreeSpace) < len(next.FreeSpace) {
		longestFrameList = next.FreeSpace
		shortestFrameList = cur.FreeSpace
	}

	var startGreater, endGreater time.Time
	var overlaps []source.Frame

	for li := 0; li < len(longestFrameList); li++ {
		curLFrame := longestFrameList[li]

		for si := 0; si < len(shortestFrameList); si++ {
			curSFrame := shortestFrameList[si]

			if !isOverlapPossible(&curLFrame, &curSFrame) {
				continue
			}

			startGreater = getStartGreater(&curLFrame, &curSFrame)
			endGreater = getEndGreater(&curLFrame, &curSFrame)

			if !areFrameComponentsNil(startGreater, endGreater) && isFrameWithinDuration(startGreater, endGreater, duration) {
				overlaps = append(overlaps, source.Frame{Start: startGreater, End: endGreater})
				startGreater, endGreater = time.Time{}, time.Time{}
			}
		}
	}

	return overlaps
}

func getStartGreater(f *source.Frame, f2 *source.Frame) time.Time {
	if isTimeInFrameIncl(f.Start, f2) {
		return f.Start
	}

	if isTimeInFrameIncl(f2.Start, f) {
		return f2.Start
	}

	return time.Time{}
}

func getEndGreater(f *source.Frame, f2 *source.Frame) time.Time {
	if isTimeInFrameIncl(f.End, f2) {
		return f.End
	}

	if isTimeInFrameIncl(f2.End, f) {
		return f2.End
	}

	return time.Time{}
}

func isOverlapPossible(f *source.Frame, f2 *source.Frame) bool {
	return isTimeInFrameExcl(f.Start, f2) || isTimeInFrameExcl(f.End, f2) || isTimeInFrameExcl(f2.Start, f) || isTimeInFrameExcl(f2.End, f)
}

// isTimeInFrameIncl checks if the given time is between the start and end of the frame, inclusively
func isTimeInFrameIncl(t time.Time, frame *source.Frame) bool {
	return (t.After(frame.Start) || t.Equal(frame.Start)) && (t.Before(frame.End) || t.Equal(frame.End))
}

// isTimeInFrameExcl checks if the given time is between the start and end of the frame, exclusively
func isTimeInFrameExcl(t time.Time, frame *source.Frame) bool {
	return t.After(frame.Start) && t.Before(frame.End)
}

func areFrameComponentsNil(start, end time.Time) bool {
	return start.IsZero() || end.IsZero()
}

func isFrameWithinDuration(start, end time.Time, duration time.Duration) bool {
	return math.Abs(end.Sub(start).Seconds()) >= duration.Seconds()
}

func splitFramesInDuration(frames []source.Frame, duration time.Duration) []source.Frame {
	result := make([]source.Frame, 0)

	for _, frame := range frames {
		start := frame.Start
		end := frame.End

		for {
			if end.Sub(start) > duration {
				result = append(result, source.Frame{Start: start, End: start.Add(duration)})
				start = start.Add(duration)
			}

			if end.Sub(start) == duration {
				result = append(result, source.Frame{Start: start, End: end})
				break
			}

			if end.Sub(start) < duration {
				break
			}
		}
	}

	return result
}
