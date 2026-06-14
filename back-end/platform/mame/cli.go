package mame

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Game struct {
	Name  string
	RomID string
}

func Search(mameCmd, query string) ([]Game, error) {
	name, args := parseCommand(mameCmd)
	args = append(args, "-listfull")
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return nil, fmt.Errorf("mame -listfull: %w", err)
	}
	return parseListFull(out, query), nil
}

func Launch(mameCmd, appId string) error {
	name, args := parseCommand(mameCmd)
	args = append(args, appId)
	cmd := exec.Command(name, args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch mame game %s: %w", appId, err)
	}
	return nil
}

func parseCommand(cmd string) (string, []string) {
	fields := strings.Fields(cmd)
	if len(fields) == 0 {
		return cmd, nil
	}
	return fields[0], fields[1:]
}

func parseListFull(output []byte, query string) []Game {
	normalizedQuery := strings.ReplaceAll(strings.ToLower(query), " ", "")
	var games []Game

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		romID, name, ok := parseLine(line)
		if !ok {
			continue
		}
		normalizedName := strings.ReplaceAll(strings.ToLower(name), " ", "")
		if !strings.Contains(normalizedName, normalizedQuery) {
			continue
		}
		games = append(games, Game{Name: name, RomID: romID})
	}
	return games
}

func parseLine(line string) (romID, name string, ok bool) {
	qStart := strings.IndexByte(line, '"')
	if qStart < 0 {
		return "", "", false
	}
	qEnd := strings.LastIndexByte(line, '"')
	if qEnd <= qStart {
		return "", "", false
	}
	romID = strings.TrimSpace(line[:qStart])
	if romID == "" {
		return "", "", false
	}
	name = line[qStart+1 : qEnd]
	return romID, name, true
}
