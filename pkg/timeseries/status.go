package timeseries

import "time"

type Status interface {
	Timestamp() time.Time
	Fields() map[string]any
}
