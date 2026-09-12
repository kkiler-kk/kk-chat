package clock

import (
	"time"

	"server-go/internal/usecase/port"
)

// Real 真实时钟，实现 port.Clock。
type Real struct{}

func (Real) Now() time.Time { return time.Now() }

var _ port.Clock = Real{}
