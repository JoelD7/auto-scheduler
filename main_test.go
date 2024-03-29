package main

import (
	"testing"
	"time"

	"github.com/JoelD7/auto-scheduler/source"

	"github.com/stretchr/testify/require"
)

func TestGetSuggestedMeetingSlots(t *testing.T) {
	c := require.New(t)

	users := getTestUsers()
	john := users[0]
	annie := users[1]
	mark := users[2]

	suggestions := getSuggestedMeetingSlots(time.Minute*30, john, annie)

	c.Equal("{11:30, 12:00}", suggestions[0].String())
	c.Equal("{15:00, 15:30}", suggestions[1].String())
	c.Equal("{15:30, 16:00}", suggestions[2].String())

	suggestions = getSuggestedMeetingSlots(time.Minute*30, john, annie, mark)
	c.Len(suggestions, 1)
	c.Equal("{11:30, 12:00}", suggestions[0].String())

	suggestions = getSuggestedMeetingSlots(time.Minute*12, john, annie)
	c.Equal("{11:30, 11:42}", suggestions[0].String())
	c.Equal("{11:42, 11:54}", suggestions[1].String())
	c.Equal("{15:00, 15:12}", suggestions[2].String())
	c.Equal("{15:12, 15:24}", suggestions[3].String())
	c.Equal("{15:24, 15:36}", suggestions[4].String())
	c.Equal("{15:36, 15:48}", suggestions[5].String())
	c.Equal("{15:48, 16:00}", suggestions[6].String())
}

func TestCalendarToMettingFrame(t *testing.T) {
	c := require.New(t)

	meetings := []string{"9:00", "10:30"}

	frame, err := source.CalendarToMeetingFrame(meetings)
	c.Nil(err)

	c.Equal("{09:00, 10:30}", frame.String())
}

func TestGetOverlapBetween(t *testing.T) {
	c := require.New(t)

	users := getTestUsers()
	john := users[0]
	annie := users[1]
	duration := time.Minute * 30

	overlaps := getOverlapBetween(duration, john, annie)
	c.Equal("{11:30, 12:00}", overlaps[0].String())
	c.Equal("{15:00, 16:00}", overlaps[1].String())
}

func TestIsTimeInFrame(t *testing.T) {
	c := require.New(t)

	ti, _ := time.Parse("15:04", "08:00")
	start, _ := time.Parse(source.TimeLayout, "11:30")
	end, _ := time.Parse(source.TimeLayout, "12:30")
	frame := source.Frame{Start: start, End: end}
	c.False(isTimeInFrameIncl(ti, &frame))

	ti, _ = time.Parse("15:04", "09:00")
	start, _ = time.Parse(source.TimeLayout, "11:30")
	end, _ = time.Parse(source.TimeLayout, "12:30")
	frame = source.Frame{Start: start, End: end}
	c.False(isTimeInFrameIncl(ti, &frame))

	ti, _ = time.Parse("15:04", "10:30")
	start, _ = time.Parse(source.TimeLayout, "11:30")
	end, _ = time.Parse(source.TimeLayout, "12:30")
	frame = source.Frame{Start: start, End: end}
	c.False(isTimeInFrameIncl(ti, &frame))

	ti, _ = time.Parse("15:04", "12:00")
	start, _ = time.Parse(source.TimeLayout, "11:30")
	end, _ = time.Parse(source.TimeLayout, "12:30")
	frame = source.Frame{Start: start, End: end}
	c.True(isTimeInFrameIncl(ti, &frame))

	ti, _ = time.Parse("15:04", "13:00")
	start, _ = time.Parse(source.TimeLayout, "11:30")
	end, _ = time.Parse(source.TimeLayout, "12:30")
	frame = source.Frame{Start: start, End: end}
	c.False(isTimeInFrameIncl(ti, &frame))

	ti, _ = time.Parse("15:04", "16:00")
	start, _ = time.Parse(source.TimeLayout, "11:30")
	end, _ = time.Parse(source.TimeLayout, "12:30")
	frame = source.Frame{Start: start, End: end}
	c.False(isTimeInFrameIncl(ti, &frame))
}

func TestSplitFramesInDuration(t *testing.T) {
	c := require.New(t)

	users := getTestUsers()
	john := users[0]
	annie := users[1]
	duration := time.Minute * 30

	overlaps := getOverlapBetween(duration, john, annie)
	frames := splitFramesInDuration(overlaps, duration)

	c.Equal("{11:30, 12:00}", frames[0].String())
	c.Equal("{15:00, 15:30}", frames[1].String())
	c.Equal("{15:30, 16:00}", frames[2].String())

	duration = time.Minute * 60
	overlaps = getOverlapBetween(duration, john, annie)
	frames = splitFramesInDuration(overlaps, duration)

	c.Equal("{15:00, 16:00}", frames[0].String())

	duration = time.Minute * 45
	overlaps = getOverlapBetween(duration, john, annie)
	frames = splitFramesInDuration(overlaps, duration)

	c.Equal("{15:00, 15:45}", frames[0].String())
}

func getTestUsers() []source.User {
	jMeetings := [][]string{{"9:00", "10:30"}, {"12:00", "13:00"}, {"16:00", "18:00"}}
	jDailyBounds := []string{"8:00", "18:00"}

	aMeetings := [][]string{{"10:00", "11:30"}, {"12:30", "14:30"}, {"14:30", "15:00"}, {"16:00", "17:00"}}
	aDailyBounds := []string{"10:00", "18:30"}

	mMeetings := [][]string{{"10:00", "11:00"}, {"13:00", "14:30"}, {"15:00", "17:00"}, {"17:30", "18:00"}}
	mDailyBounds := []string{"08:00", "18:30"}

	johnUser, _ := source.NewUser(jMeetings, jDailyBounds)
	annieUser, _ := source.NewUser(aMeetings, aDailyBounds)
	markUser, _ := source.NewUser(mMeetings, mDailyBounds)

	return []source.User{*johnUser, *annieUser, *markUser}
}
