package memstore

import (
	"time"
	"github.com/devxfactor/quicklog/utils"
	"errors"
	"fmt"
)

type Memstore interface {
	Trace(time time.Time, line string)
	Debug(time time.Time, line string)
	Info(time time.Time, line string)
	Warn(time time.Time, line string)
	Error(time time.Time, line string)
	Log(time time.Time, level Level, line string) error
	Logf(time time.Time, level Level, format string, args ...interface{}) error

	Each(f func (Entry))
}

type Entry interface {
	Time() time.Time
	Level() Level
	Line() string
}

type Level string

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
)

var (
	LEVELS = []string{ "TRACE", "DEBUG", "INFO", "WARN", "ERROR" }
)

type memstore struct {
	slices [][]entry
}

type entry struct {
	time time.Time
	level int
	line string
}

func NewMemstore() Memstore {
	m := &memstore{slices: [][]entry{}}
	m.slices = append(m.slices, []entry{})
	return m
}

func (m *memstore) Trace(time time.Time, line string) {
	m.log(time, TRACE, line)
}
func (m *memstore) Tracef(time time.Time, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args)
	m.Trace(time, line)
}

func (m *memstore) Debug(time time.Time, line string) {
	m.log(time, DEBUG, line)
}
func (m *memstore) Debugf(time time.Time, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args)
	m.Debug(time, line)
}

func (m *memstore) Info(time time.Time, line string) {
	m.log(time, INFO, line)
}
func (m *memstore) Infof(time time.Time, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args)
	m.Info(time, line)
}

func (m *memstore) Warn(time time.Time, line string) {
	m.log(time, WARN, line)
}
func (m *memstore) Warnf(time time.Time, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args)
	m.Warn(time, line)
}

func (m *memstore) Error(time time.Time, line string) {
	m.log(time, ERROR, line)
}
func (m *memstore) Errorf(time time.Time, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args)
	m.Error(time, line)
}

func (m *memstore) Logf(time time.Time, level Level, format string, args... interface{}) error {
	line := fmt.Sprintf(format, args)
	return m.Log(time, level, line)
}

func (m *memstore) Log(time time.Time, level Level, line string) error {
	l, err := levelInt(level)
	if err != nil {
		return err
	}
	m.log(time, l, line)
	return nil
}

func (m *memstore) log(time time.Time, level int, line string) {
	m.slices[0] = append(m.slices[0], entry{time: time, level: level, line: line})
}

func (m *memstore) Each(f func (Entry)) {
	for _, slice := range m.slices {
		for _, entry := range slice {
			f(entry)
		}
	}
}

func (e entry) Time() time.Time {
	return e.time
}

func (e entry) Level() Level {
	return Level(LEVELS[e.level])
}

func (e entry) Line() string {
	return e.line
}

func levelInt(level Level) (int, error) {
	index, err := utils.StringIndex(string(level), LEVELS)
	if err != nil {
		return -1, errors.New("invalid level")
	}
	return index, nil
}