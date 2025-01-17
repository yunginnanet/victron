package victron

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/rosenstand/go-vedirect/vedirect"
)

var ErrBadChecksumModulus = errors.New("bad checksum modulus")

type Blocks struct {
	Fields    map[string]string
	blocks    []vedirect.Block
	lastBlock vedirect.Block // for debugging
	mu        sync.RWMutex
}

func newBlocks() *Blocks {
	return &Blocks{Fields: make(map[string]string)}
}

func (b *Blocks) Validate() bool {
	b.mu.RLock()
	for _, block := range b.blocks {
		if !block.Validate() {
			return false
		}
	}
	_, ok := b.Fields[PrefixSerial]
	b.mu.RUnlock()
	return ok
}

func (b *Blocks) DropInvalid() (n int) {
	b.mu.Lock()
	var valid []vedirect.Block
	for _, block := range b.blocks {
		if block.Validate() {
			valid = append(valid, block)
			continue
		}
		n++
	}
	b.blocks = valid
	b.mu.Unlock()
	return
}

func (b *Blocks) Len() int {
	b.mu.RLock()
	blen := len(b.blocks)
	b.mu.RUnlock()
	return blen
}

func GetFields(nb vedirect.Block, target map[string]string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			println("[WARN] GetFields panic:", r)
			if _, ok := r.(error); !ok {
				err = errors.New(fmt.Sprintf("%v", r))
			} else {
				err = r.(error)
			}
		}
	}()

	reflected := reflect.ValueOf(&nb).Elem().FieldByName("fields")
	if !reflected.IsValid() {
		return errors.New("reflection failure: no fields in block")
	}
	if reflected.Kind() != reflect.Map {
		return errors.New("reflection failure: fields is not a map")
	}
	for _, key := range reflected.MapKeys() {
		target[key.String()] = reflected.MapIndex(key).String()
	}
	return nil
}

func (b *Blocks) readBlock(stream *vedirect.Stream) error {
	var nb vedirect.Block
	var n int
	if nb, n = stream.ReadBlock(); n != 0 || !nb.Validate() {
		b.mu.Lock()
		b.lastBlock = nb // for debugging, so set it even if it's invalid
		b.mu.Unlock()
		return ErrBadChecksumModulus
	}
	if err := GetFields(nb, b.Fields); err != nil {
		return err
	}
	b.mu.Lock()
	b.lastBlock = nb
	b.blocks = append(b.blocks, nb)
	b.mu.Unlock()
	return nil
}
