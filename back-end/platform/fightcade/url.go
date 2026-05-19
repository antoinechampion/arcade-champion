package fightcade

import (
	"fmt"
	"os/exec"
	"runtime"
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

func openURL(url string) error {
	fmt.Println("opening url: %s", url)
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}
