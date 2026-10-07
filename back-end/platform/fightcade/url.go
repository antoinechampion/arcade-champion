package fightcade

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var nativeGGPO = map[string]bool{
	"flycast":     true,
	"duckstation": true,
	"custom":      true,
}

func buildMatchURL(emulator, gameid, quarkid string, playerid, port, delay int, ranked bool, token string) string {
	r := 0
	if ranked {
		r = 1
	}
	if nativeGGPO[emulator] {
		return fmt.Sprintf("fcade://launch/%s/%s/%s.%d,%d,%d,%d/%s",
			emulator, gameid, quarkid, playerid, port, delay, r, token)
	}
	return fmt.Sprintf("fcade://served/%s/%s/%s.%d,%d,%d,%d",
		emulator, gameid, quarkid, playerid, port, delay, r)
}

func buildPlayURL(emulator, gameid string) string {
	return fmt.Sprintf("fcade://play/%s/%s", emulator, gameid)
}

func buildTrainingURL(emulator, gameid string) string {
	return fmt.Sprintf("fcade://training/%s/%s", emulator, gameid)
}

func parseCommand(cmd string) (string, []string) {
	fields := strings.Fields(cmd)
	if len(fields) == 0 {
		return cmd, nil
	}
	return fields[0], fields[1:]
}

func openURL(fightcadeCmd, url string) error {
	log.Printf("opening url: %s with command: %s\n", url, fightcadeCmd)
	name, args := parseCommand(fightcadeCmd)
	args = append(args, url)
	return exec.Command(name, args...).Start()
}
