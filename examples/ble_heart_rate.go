// +build example
//
// Do not build by default.

/*
 How to run
 Pass the Bluetooth address or name as the first param:

	go run examples/ble_heart_rate.go BB-1234

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	info := ble.NewDeviceInformationDriver(bleAdaptor)
	battery := ble.NewBatteryDriver(bleAdaptor)
	heartRate := ble.NewHeartRateDriver(bleAdaptor)

	work := func() {
		// info
		fmt.Println("=== Device Information ===")
		fmt.Println("Model number:", info.GetModelNumber())
		fmt.Println("Firmware rev:", info.GetFirmwareRevision())
		fmt.Println("Hardware rev:", info.GetHardwareRevision())
		fmt.Println("Manufacturer name:", info.GetManufacturerName())
		// battery
		fmt.Println("=== Battery Level ===")
		fmt.Println("Battery level:", battery.GetBatteryLevel())
		// heartRate
		fmt.Println("=== Body Sensor Location ===")
		loc, _ := heartRate.GetBodySensorLocation()
		fmt.Println("Body sensor location:", loc)
		fmt.Println("=== Heart Rate ===")
		heartRate.SubscribeHeartRate()
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{battery, heartRate},
		work,
	)

	robot.Start()
}
