package main

import "time"
import "github.com/sparrc/go-ping"

type Statistics struct {
	AvgRtt float32
}

func Ping(host string) (*Statistics, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return nil, err
	}

	pinger.SetPrivileged(true)
	pinger.Timeout = 4 * time.Second
	pinger.Count = 3

	pinger.Run()
	statistics := pinger.Statistics()

	return &Statistics{AvgRtt: float32(statistics.AvgRtt) / float32(time.Millisecond)}, nil
}