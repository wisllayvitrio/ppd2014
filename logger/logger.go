package logger

import (
	"fmt"
	"time"
	"runtime"
)

/*type MemStats struct {
	// General statistics.
	Alloc      uint64 // bytes allocated and still in use
	TotalAlloc uint64 // bytes allocated (even if freed)
	Sys        uint64 // bytes obtained from system (sum of XxxSys below)
	Lookups    uint64 // number of pointer lookups
	Mallocs    uint64 // number of mallocs
	Frees      uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    uint64 // bytes allocated and still in use
	HeapSys      uint64 // bytes obtained from system
	HeapIdle     uint64 // bytes in idle spans
	HeapInuse    uint64 // bytes in non-idle span
	HeapReleased uint64 // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  uint64 // bytes used by stack allocator
	StackSys    uint64
	MSpanInuse  uint64 // mspan structures
	MSpanSys    uint64
	MCacheInuse uint64 // mcache structures
	MCacheSys   uint64
	BuckHashSys uint64 // profiling bucket hash table
	GCSys       uint64 // GC metadata
	OtherSys    uint64 // other system allocations

	// Garbage collector statistics.
	NextGC       uint64 // next collection will happen when HeapAlloc ≥ this amount
	LastGC       uint64 // end time of last collection (nanoseconds since 1970)
	PauseTotalNs uint64
	PauseNs      [256]uint64 // circular buffer of recent GC pause durations, most recent at [(NumGC+255)%256]
	PauseEnd     [256]uint64 // circular buffer of recent GC pause end times
	NumGC        uint32
	EnableGC     bool
	DebugGC      bool

	// Per-size allocation statistics.
	// 61 is NumSizeClasses in the C code.
	BySize [61]struct {
		Size    uint32
		Mallocs uint64
		Frees   uint64
	}
}*/

// Stores the total time of count events for something named name
type timeData struct {
	total time.Time
	count int
}

type Logger struct {
	timeMap map[string]timeData
	// TODO: MemmoryMaps, CounterMaps, CPUMaps and whatnot
}

func NewLogger() *Logger {
	res := new(Logger)
	res.timeMap = make(map[string]timeData)
	
	return res
}

func (l *Logger) AddTime(name string, t time.Duration) {
	_, exists := l.timeMap[name]
	if !exists {
		// Initialize (Time{} is Epoch: 0001-01-01 00:00:00 +0000 UTC)
		l.timeMap[name] = timeData{time.Time{}, 0}
	}
	// Update (maps cant be set directly)
	aux := l.timeMap[name]
	aux.total = aux.total.Add(t)
	aux.count++
	l.timeMap[name] = aux
}

func (l *Logger) GetMean(name string) time.Duration {
	data, exists := l.timeMap[name]
	if !exists {
		// Ugly code to return zero duration (return 0.0 seems to work)
		return time.Time{}.Sub(time.Time{})
	}
	// REALLY UGLY code to get the total duration and divide by count
	totalDuration := data.total.Sub(time.Time{})
	truncatedMean := totalDuration.Nanoseconds() / int64(data.count)
	timeMean := time.Unix(0, truncatedMean)
	durationMean := timeMean.Sub(time.Unix(0, 0))
	return durationMean
}

// TODO: This was just a test function, probably should be deleted
func PrintEntireMemStats() {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	fmt.Println("LOG - MemStats:")
	fmt.Println("LOG - Alloc:", stats.Alloc)
	fmt.Println("LOG - TotalAlloc:", stats.TotalAlloc)
	fmt.Println("LOG - Sys:", stats.Sys)
	fmt.Println("LOG - Lookups:", stats.Lookups)
	fmt.Println("LOG - Mallocs:", stats.Mallocs)
	fmt.Println("LOG - Frees:", stats.Frees)
	fmt.Println("LOG - HeapAloc:", stats.HeapAlloc)
	fmt.Println("LOG - HeapSys:", stats.HeapSys)
	fmt.Println("LOG - HeapIdle:", stats.HeapIdle)
	fmt.Println("LOG - HeapInuse:", stats.HeapInuse)
	fmt.Println("LOG - HeapReleased:", stats.HeapReleased)
	fmt.Println("LOG - HeapObjects:", stats.HeapObjects)
	fmt.Println("LOG - StackInuse:", stats.StackInuse)
	fmt.Println("LOG - StackSys:", stats.StackSys)
	fmt.Println("LOG - MSpanInuse:", stats.MSpanInuse)
	fmt.Println("LOG - MSpanSys:", stats.MSpanSys)
	fmt.Println("LOG - MCacheInuse:", stats.MCacheInuse)
	fmt.Println("LOG - MCacheSys:", stats.MCacheSys)
	fmt.Println("LOG - BuckHashSys:", stats.BuckHashSys)
	fmt.Println("LOG - GCSys:", stats.GCSys)
	fmt.Println("LOG - OtherSys:", stats.OtherSys)
	fmt.Println("LOG - NextGC:", stats.NextGC)
	fmt.Println("LOG - LastGC:", stats.LastGC)
	fmt.Println("LOG - PauseTotalNs:", stats.PauseTotalNs)
	fmt.Println("LOG - PauseNs:", stats.PauseNs)
	//fmt.Println("LOG - PauseEnd:", stats.PauseEnd)
	fmt.Println("LOG - NumGC:", stats.NumGC)
	fmt.Println("LOG - EnableGC:", stats.EnableGC)
	fmt.Println("LOG - DebugGC:", stats.DebugGC)
	fmt.Println("LOG - BySize:", stats.BySize)
}