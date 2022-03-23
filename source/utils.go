package source

import "time"

const (
	TimeLayout = "15:04"
)

func CalendarToMeetingFrame(stringTimeFrame []string) (Frame, error) {
	start, err := time.Parse(TimeLayout, stringTimeFrame[0])
	if err != nil {
		return Frame{}, err
	}

	end, err := time.Parse(TimeLayout, stringTimeFrame[1])
	if err != nil {
		return Frame{}, err
	}

	return Frame{start, end}, nil
}
