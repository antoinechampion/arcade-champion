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

func Search(mamePath, query string) ([]Game, error) {
	out, err := exec.Command(mamePath, "-listfull").Output()
	if err != nil {
		return nil, fmt.Errorf("mame -listfull: %w", err)
	}
	return parseListFull(out, query), nil
}

func parseListFull(output []byte, query string) []Game {
	lowerQuery := strings.ToLower(query)
	var games []Game

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		romID, name, ok := parseLine(line)
		if !ok {
			continue
		}
		if !strings.Contains(strings.ToLower(name), lowerQuery) {
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
