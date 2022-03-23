package source

import (
	"errors"
)

type User struct {
	DailyBounds     Frame
	MeetingSchedule []Frame
	FreeSpace       []Frame
}

var ErrInvalidParams = errors.New("invalid parameters for user creation")

func NewUser(meetings [][]string, dailyBounds []string) (*User, error) {
	if len(meetings) == 0 || len(dailyBounds) == 0 {
		return nil, ErrInvalidParams
	}

	meetingFrames := make([]Frame, 0, len(meetings))
	for _, meeting := range meetings {
		frame, err := CalendarToMeetingFrame(meeting)
		if err != nil {
			return nil, ErrInvalidParams
		}

		meetingFrames = append(meetingFrames, frame)
	}

	dailyBoundsFrame, err := CalendarToMeetingFrame(dailyBounds)
	if err != nil {
		return nil, ErrInvalidParams
	}

	return &User{
		DailyBounds:     dailyBoundsFrame,
		MeetingSchedule: meetingFrames,
		FreeSpace:       getFreeSpaceOf(meetingFrames, dailyBoundsFrame),
	}, nil
}

func getFreeSpaceOf(meetingFrames []Frame, dailyBounds Frame) []Frame {
	freeSpace := make([]Frame, 0, len(meetingFrames)-1)

	firstMeeting := meetingFrames[0]
	if dailyBounds.Start.Before(firstMeeting.Start) {
		freeSpace = append(freeSpace, Frame{Start: dailyBounds.Start, End: firstMeeting.Start})
	}

	for i := 0; i < len(meetingFrames)-1; i++ {
		startFree := meetingFrames[i].End
		endFree := meetingFrames[i+1].Start

		if !startFree.Equal(endFree) {
			freeSpace = append(freeSpace, Frame{Start: startFree, End: endFree})
		}
	}

	lastMeetTime := meetingFrames[len(meetingFrames)-1]
	if dailyBounds.End.After(lastMeetTime.End) {
		freeSpace = append(freeSpace, Frame{Start: lastMeetTime.End, End: dailyBounds.End})
	}

	return freeSpace
}
