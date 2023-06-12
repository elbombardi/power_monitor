package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

func main() {
	batteryState := readBatteryState()
	fmt.Println("Battery State:", batteryState)
	for {
		time.Sleep(time.Millisecond * 100)
		newBatteryState := readBatteryState()
		if batteryState != newBatteryState {
			fmt.Println("Battery State:", newBatteryState)
		}

		batteryState = newBatteryState
	}
}

func readBatteryState() string {
	upowerOuput, err := runCommand("upower -i /org/freedesktop/UPower/devices/battery_BAT0 | grep 'state:'")
	if err != nil {
		log.Fatal("Error while reading battery state: ", err)
	}
	state := strings.TrimSpace(string(upowerOuput))
	state = strings.TrimPrefix(state, "state:")
	state = strings.TrimSpace(state)
	return state
}

func runCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
