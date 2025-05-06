package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Laps                   int           `json:"laps"`
	LapLen                 float64       `json:"lapLen"`
	PenaltyLen             float64       `json:"penaltyLen"`
	FiringLines            int           `json:"firingLines"`
	Start                  string        `json:"start"`
	StartDelta             time.Duration `json:"startDelta"`
	Targets                int
	CompetitionFileName    string
	OutputLogFileName      string
	OutgoingEventsFileName string
	ReportFileName         string
}

func New(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return parse(data)
}

func parse(data []byte) (*Config, error) {
	type configJSON struct {
		Laps        int     `json:"laps"`
		LapLen      float64 `json:"lapLen"`
		PenaltyLen  float64 `json:"penaltyLen"`
		FiringLines int     `json:"firingLines"`
		Start       string  `json:"start"`
		StartDelta  string  `json:"startDelta"`
		Targets     int
	}

	var temp configJSON
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if _, err := time.Parse("15:04:05", temp.Start); err != nil {
		return nil, fmt.Errorf("invalid start time format, expected HH:MM:SS: %w", err)
	}

	var startDelta time.Duration
	if temp.StartDelta != "" {
		parts := strings.Split(temp.StartDelta, ":")
		if len(parts) < 2 || len(parts) > 3 {
			return nil, fmt.Errorf("invalid startDelta format, expected HH:MM:SS or MM:SS")
		}

		var hours, minutes, seconds int
		var err error
		if len(parts) == 3 {
			hours, err = strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid hours in startDelta: %w", err)
			}
			minutes, err = strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid minutes in startDelta: %w", err)
			}
			seconds, err = strconv.Atoi(parts[2])
			if err != nil {
				return nil, fmt.Errorf("invalid seconds in startDelta: %w", err)
			}
		} else {
			minutes, err = strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid minutes in startDelta: %w", err)
			}
			seconds, err = strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid seconds in startDelta: %w", err)
			}
		}
		startDelta = time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	}

	var ok bool
	inputFilename, ok := os.LookupEnv("COMPETITION_FILE")
	if !ok || inputFilename == "" {
		return nil, fmt.Errorf("COMPETITION_FILE environment variable not set")
	}

	outputFilename, ok := os.LookupEnv("OUTPUT_LOG_FILE")
	if !ok || outputFilename == "" {
		return nil, fmt.Errorf("OUTPUT_LOG_FILE environment variable not set")
	}

	outgoingEventsFilename, ok := os.LookupEnv("OUTGOING_EVENTS_FILE")
	if !ok || outgoingEventsFilename == "" {
		return nil, fmt.Errorf("OUTGOING_EVENTS_FILE environment variable not set")
	}

	reportFilename, ok := os.LookupEnv("REPORT_FILE")
	if !ok || reportFilename == "" {
		return nil, fmt.Errorf("REPORT_FILE environment variable not set")
	}

	cfg := &Config{
		Laps:                   temp.Laps,
		LapLen:                 temp.LapLen,
		PenaltyLen:             temp.PenaltyLen,
		FiringLines:            temp.FiringLines,
		Start:                  temp.Start,
		StartDelta:             startDelta,
		Targets:                5,
		CompetitionFileName:    inputFilename,
		OutputLogFileName:      outputFilename,
		OutgoingEventsFileName: outgoingEventsFilename,
		ReportFileName:         reportFilename,
	}

	if cfg.Laps <= 0 {
		return nil, fmt.Errorf("laps must be positive")
	}
	if cfg.LapLen <= 0 {
		return nil, fmt.Errorf("lapLen must be positive")
	}
	if cfg.PenaltyLen <= 0 {
		return nil, fmt.Errorf("penaltyLen must be positive")
	}
	if cfg.FiringLines <= 0 {
		return nil, fmt.Errorf("firingLines must be positive")
	}
	if cfg.StartDelta <= 0 {
		return nil, fmt.Errorf("startDelta must be positive")
	}

	return cfg, nil
}
