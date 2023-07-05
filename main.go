package main

import (
	"context"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/elbombardi/power_monitor/models"
	storage "github.com/elbombardi/power_monitor/storage/sqlc"
	"github.com/elbombardi/power_monitor/util"
)

func main() {
	err := util.LoadConfig()
	if err != nil {
		log.Fatal("Error while loading configuration: ", err)
	}
	dataStore, err := storage.GetStoreInstance()
	if err != nil {
		log.Fatal("Error while initializing data store: ", err)
	}
	defer storage.Finalize()
	logChannel := make(chan *models.PowerStateChange)
	defer close(logChannel)

	go powerStateLogger(logChannel, dataStore)
	watchPowerState(logChannel)
}

func powerStateLogger(logChannel chan *models.PowerStateChange, dataStore storage.DataStore) {
	for {
		powerStateChange := <-logChannel
		err := dataStore.Store(context.Background(), powerStateChange)
		for err != nil {
			log.Println("Error while storing power state change: ", err, " retrying in 1 second")
			time.Sleep(time.Second)
			err = dataStore.Store(context.Background(), powerStateChange)
		}
		log.Printf("New log saved : #%v %v Power %v \n", powerStateChange.Id, powerStateChange.Time.Format("2006-01-02 15:04:05"), powerStateChange.PowerState)
	}
}

func logPowerStateChange(logChannel chan *models.PowerStateChange, powerStateChange *models.PowerStateChange) {
	logChannel <- powerStateChange
	log.Printf("New log added : #%v %v Power %v \n",
		powerStateChange.Id,
		powerStateChange.Time.Format("2006-01-02 15:04:05"),
		powerStateChange.PowerState)
}

func watchPowerState(logChannel chan *models.PowerStateChange) {
	powerState := detectPowerState()
	counter := 1
	logPowerStateChange(logChannel, &models.PowerStateChange{
		Time:       time.Now(),
		PowerState: powerState,
	})
	for {
		time.Sleep(time.Millisecond * 100)
		newPowerState := detectPowerState()
		if powerState != newPowerState {
			counter++
			logPowerStateChange(logChannel, &models.PowerStateChange{
				Id:         counter,
				Time:       time.Now(),
				PowerState: newPowerState,
			})
		}
		powerState = newPowerState
	}
}

func detectPowerState() models.PowerState {
	upowerOuput, err := runCommand("upower -i /org/freedesktop/UPower/devices/battery_BAT0 | grep 'state:'")
	if err != nil {
		log.Fatal("Error while reading battery state: ", err)
	}
	state := strings.TrimSpace(string(upowerOuput))
	state = strings.TrimPrefix(state, "state:")
	state = strings.TrimSpace(state)
	if state == "discharging" {
		return "OFF"
	}
	return "ON"
}

func runCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
