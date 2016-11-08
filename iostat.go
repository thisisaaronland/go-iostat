package iostat

import (
	_ "log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func PrepString(str string) []string {

	str = strings.Trim(str, " ")
	// log.Println(str)
	
	re_whitespace := regexp.MustCompile(`\s+`)
	return re_whitespace.Split(str, -1)
}

type IOStatResults struct {
	Averages *IOStatAverages
	Devices  map[string]*IOStatDevice
}

// %user   %nice %system %iowait  %steal   %idle

type IOStatAverages struct {
	User   float64
	Nice   float64
	System float64
	IOWait float64
	Steal  float64
	Idle   float64
}

// Device:         rrqm/s   wrqm/s     r/s     w/s    rkB/s    wkB/s avgrq-sz avgqu-sz   await r_await w_await  svctm  %util

type IOStatDevice struct {
	Name              string
	RRQMPerSecond     float64
	WRQMPerSecond     float64
	ReadsPerSecond    float64
	WritesPerSecond   float64
	ReadsKBPerSecond  float64
	WritesKBPerSecond float64
	AvgRQSize         float64
	AvgQUSize         float64
	Await             float64
	AwaitRead         float64
	AwaitWrite        float64
	Svctm             float64
	Util              float64
}

func NewIOStatAveragesFromString(str string) (*IOStatAverages, error) {

	parts := PrepString(str)

	user, _ := strconv.ParseFloat(parts[0], 64)
	nice, _ := strconv.ParseFloat(parts[1], 64)
	system, _ := strconv.ParseFloat(parts[2], 64)
	iowait, _ := strconv.ParseFloat(parts[3], 64)
	steal, _ := strconv.ParseFloat(parts[4], 64)
	idle, _ := strconv.ParseFloat(parts[5], 64)

	avg := IOStatAverages{
		User:   user,
		Nice:   nice,
		System: system,
		IOWait: iowait,
		Steal:  steal,
		Idle:   idle,
	}

	return &avg, nil
}

func NewIOStatDeviceFromString(str string) (*IOStatDevice, error) {

	parts := PrepString(str)

	name := parts[0]

	rrqm, _ := strconv.ParseFloat(parts[1], 64)
	wrqm, _ := strconv.ParseFloat(parts[2], 64)
	rs, _ := strconv.ParseFloat(parts[3], 64)
	ws, _ := strconv.ParseFloat(parts[4], 64)
	rkbs, _ := strconv.ParseFloat(parts[5], 64)
	wkbs, _ := strconv.ParseFloat(parts[6], 64)
	avgrqsz, _ := strconv.ParseFloat(parts[7], 64)
	avgqusz, _ := strconv.ParseFloat(parts[8], 64)
	await, _ := strconv.ParseFloat(parts[9], 64)
	await_r, _ := strconv.ParseFloat(parts[10], 64)
	await_w, _ := strconv.ParseFloat(parts[11], 64)
	svctm, _ := strconv.ParseFloat(parts[12], 64)
	util, _ := strconv.ParseFloat(parts[13], 64)

	device := IOStatDevice{
		Name:              name,
		RRQMPerSecond:     rrqm,
		WRQMPerSecond:     wrqm,
		ReadsPerSecond:    rs,
		WritesPerSecond:   ws,
		ReadsKBPerSecond:  rkbs,
		WritesKBPerSecond: wkbs,
		AvgRQSize:         avgrqsz,
		AvgQUSize:         avgqusz,
		Await:             await,
		AwaitRead:         await_r,
		AwaitWrite:        await_w,
		Svctm:             svctm,
		Util:              util,
	}

	return &device, nil
}

func NewIOStatResults() (*IOStatResults, error) {

	var out []byte
	var err error

	iters := 2

	for i := 0; i < iters; i++ {

		cmd := exec.Command("/usr/bin/iostat", "-x", "-z")
		out, err = cmd.Output()

		if err != nil {
			return nil, err
		}

		time.Sleep(time.Second * 2)
	}

	stats := strings.Trim(string(out), "\n")
	lines := strings.Split(stats, "\n")

	/*
		for _, ln := range lines {
			log.Println(ln)
		}
	*/

	avg, err := NewIOStatAveragesFromString(lines[3])

	if err != nil {
		return nil, err
	}

	devices := make(map[string]*IOStatDevice)

	for _, ln := range lines[6:] {

		device, err := NewIOStatDeviceFromString(ln)

		if err != nil {
			return nil, err
		}

		name := device.Name
		devices[name] = device
	}

	results := IOStatResults{
		Averages: avg,
		Devices:  devices,
	}

	return &results, nil
}
