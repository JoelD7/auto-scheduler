package source

import (
	"fmt"
	"time"
)

// Frame represents a time range
type Frame struct {
	Start time.Time
	End   time.Time
}

func (f *Frame) String() string {
	return fmt.Sprintf("{%s, %s}", f.Start.Format(TimeLayout), f.End.Format(TimeLayout))
}
