package timeseries

import (
	"encoding/binary"
	"os"
	"time"

	"github.com/bytedance/sonic"

	"github.com/akrylysov/pogreb"
)

type TimeSeries struct {
	store *pogreb.DB
}

func (ts *TimeSeries) Items() map[time.Time]map[string]any {
	items := make(map[time.Time]map[string]any)
	iter := ts.store.Items()
	var key, value []byte
	var err error
	for {
		key, value, err = iter.Next()
		if err != nil {
			break
		}
		var fields = make(map[string]any)
		err = sonic.Unmarshal(value, &fields)
		if err != nil {
			_, _ = os.Stderr.WriteString(err.Error())
			continue
		}
		items[bytesToTime(key)] = fields
	}
	return items
}

func New(path string) (*TimeSeries, error) {
	store, err := pogreb.Open(path, nil)
	if err != nil {
		return nil, err
	}
	return &TimeSeries{store: store}, nil
}

func (ts *TimeSeries) Close() error {
	return ts.store.Close()
}

func timeToBytes(t time.Time) []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(t.UnixNano()))
	return bs
}

func bytesToTime(bs []byte) time.Time {
	return time.Unix(0, int64(binary.BigEndian.Uint64(bs)))
}

func (ts *TimeSeries) IngestVictron(s Status) error {
	dat, err := sonic.Marshal(s.Fields())
	if err != nil {
		return err
	}
	return ts.store.Put(timeToBytes(s.Timestamp()), dat)

}
