package logger

import (
	"os"
	"log"
	"fmt"
	"time"
	"runtime"
	"syscall"
	"go/build"
)

// Stores the total time of count events for something
type timeData struct {
	total time.Time
	count int
}

type Logger struct {
	file *os.File
	delta time.Duration
	searchTime timeData
	writeTime timeData
	cpuTime int64
	memAlloc uint64
}

func NewLogger(fileName string, delta time.Duration) *Logger {
	// Open the file, or create a new one if it does not exist
	filePath := build.Default.GOPATH + "/src/github.com/wisllayvitrio/ppd2014/logs/" + fileName
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filePath)
			if err != nil {
				log.Fatal(fmt.Sprintln("ERROR creating file:", err))
			}
		} else {
			log.Fatal(fmt.Sprintln("ERROR opening file:", err))
		}
	}
	
	res := new(Logger)
	res.file = file
	res.delta = delta
	
	res.searchTime = timeData{time.Time{}, 0}
	res.writeTime = timeData{time.Time{}, 0}
	res.cpuTime = getCPUTime()
	res.memAlloc = getMemAlloc()
	
	return res
}

func (l *Logger) LogStart() {
	// Print an initial line (to separate executions)
	c, err := l.file.WriteString("----------------------------------------------------\n")
	if err != nil {
		log.Fatal(fmt.Sprintln("ERROR writing", c, "characters on log file", l.file.Name(), ":", err))
	}
	// Append on file at each delta time
	for ; ; <-time.After(l.delta) {
		sTime := l.GetTimeMean(true)
		wTime := l.GetTimeMean(false)
		cpu := l.GetCPU(l.delta)
		mem := l.GetMem()
	
		str := fmt.Sprintln(sTime, wTime, cpu, mem)
		c, err := l.file.WriteString(str)
		if err != nil {
			log.Fatal(fmt.Sprintln("ERROR writing", c, "characters on log file", l.file.Name(), ":", err))
		}
	}
}

func (l *Logger) AddTime(searchTime bool, t time.Duration) {
	if searchTime {
		l.searchTime.total = l.searchTime.total.Add(t)
		l.searchTime.count++
	} else {
		l.writeTime.total = l.writeTime.total.Add(t)
		l.writeTime.count++
	}
}

func (l *Logger) GetTimeMean(searchTime bool) int64 {
	var data timeData
	if searchTime {
		data = l.searchTime
	} else {
		data = l.writeTime
	}
	if data.count == 0 {
		return 0
	}
	// Code to get the total duration and divide by count
	totalDuration := data.total.Sub(time.Time{})
	truncatedMean := totalDuration.Nanoseconds() / int64(data.count)
	return truncatedMean
}

func (l *Logger) GetCPU(t time.Duration) float64 {
	newTime := getCPUTime()
	mean := float64(newTime - l.cpuTime) / float64(t.Nanoseconds())
	// TODO: divide by num of processors
	return (mean * 100)
}

func (l *Logger) GetMem() uint64 {
	l.memAlloc = getMemAlloc()
	return l.memAlloc
}

func getCPUTime() int64 {
	usage := syscall.Rusage{}
	syscall.Getrusage(syscall.RUSAGE_SELF, &usage)
	return usage.Utime.Usec + usage.Stime.Usec
}

func getMemAlloc() uint64 {
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	return mem.Alloc
}