package app

import (
	"Solution/config"
	"Solution/internal/controller"
	"Solution/internal/repo"
	"Solution/internal/usecase"
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

func Run(logger *zap.Logger, cfg *config.Config) {
	repository := repo.New(logger)
	sender := usecase.NewSender(logger, cfg)
	competition := usecase.NewCompetition(logger, cfg, sender, repository, repository, repository)
	ctrl := controller.New(logger, competition, competition, competition)

	runCompetition(logger, ctrl, cfg.CompetitionFileName, cfg.OutputLogFileName)
}

func runCompetition(logger *zap.Logger, ctrl controller.Controller, competitionFile, outputLog string) {
	file, err := os.Open(competitionFile)
	if err != nil {
		logger.Fatal("failed to open file", zap.Error(err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		eventID, err := extractEventID(line)
		if err != nil {
			logger.Error("Failed to extract event ID",
				zap.Int("line", lineNumber),
				zap.String("event", line),
				zap.Error(err))
			continue
		}

		logMsg, err := processEvent(eventID, line, ctrl)
		if err != nil {
			logger.Error("Failed to process event",
				zap.Int("line", lineNumber),
				zap.String("event", line),
				zap.Error(err))
			continue
		}

		if outputLog != "" {
			f, err := os.OpenFile(outputLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logger.Error("Failed to open log file", zap.Error(err))
			}

			if _, err := f.WriteString(logMsg + "\n"); err != nil {
				logger.Error("Failed to write log message", zap.Error(err))
			}
			err = f.Close()
			if err != nil {
				logger.Error("Failed to close log file", zap.Error(err))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Fatal("failed to scan file", zap.Error(err))
	}

	err = ctrl.WriteReport()
	if err != nil {
		logger.Fatal("failed to write report", zap.Error(err))
	}
}

func extractEventID(line string) (int, error) {
	if !strings.HasPrefix(line, "[") {
		return 0, fmt.Errorf("invalid event format: missing time prefix")
	}

	endTimeIdx := strings.Index(line, "]")
	if endTimeIdx == -1 {
		return 0, fmt.Errorf("invalid event format: missing closing bracket")
	}

	rest := strings.TrimSpace(line[endTimeIdx+1:])
	parts := strings.Fields(rest)
	if len(parts) < 1 {
		return 0, fmt.Errorf("invalid event format: missing EventID")
	}

	eventID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid EventID: %w", err)
	}

	return eventID, nil
}

func processEvent(eventID int, line string, ctrl controller.Controller) (string, error) {
	switch eventID {
	case 1:
		return ctrl.Registered(line)
	case 2:
		return ctrl.StartTimeWasSet(line)
	case 3:
		return ctrl.OnStartLine(line)
	case 4:
		return ctrl.Started(line)
	case 5:
		return ctrl.OnFiringRange(line)
	case 6:
		return ctrl.TargetWasHit(line)
	case 7:
		return ctrl.LeftFiringRange(line)
	case 8:
		return ctrl.OnPenaltyLap(line)
	case 9:
		return ctrl.LeftPenaltyLap(line)
	case 10:
		return ctrl.EndedMainLap(line)
	case 11:
		return ctrl.CantContinue(line)
	default:
		return "", fmt.Errorf("unknown event ID: %d", eventID)
	}
}
