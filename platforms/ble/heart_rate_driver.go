package ble

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
)

type HeartRateDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

func NewHeartRateDriver(a BLEConnector) *HeartRateDriver {
	n := &HeartRateDriver{
		name:       gobot.DefaultName("Heart Rate"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}
	return n
}

func (b *HeartRateDriver) Name() string { return b.name }

func (b *HeartRateDriver) SetName(n string) { b.name = n }

func (b *HeartRateDriver) Connection() gobot.Connection { return b.connection }

func (b *HeartRateDriver) Start() (err error) { return }

func (b *HeartRateDriver) Halt() (err error) { return }

func (b *HeartRateDriver) adaptor() BLEConnector {
	return b.Connection().(BLEConnector)
}

// HRS(Heart Rate Service) characteristics
const (
	cUUIDHeartRateMeasurement  = "2a37"
	cUUIDBodySensorLocation    = "2a38"
	cUUIDHeartRateControlPoint = "2a39"
)

// BodySensorLocation value
var mBodySensorLocation = map[uint8]string{
	0: "Other",
	1: "Chest",
	2: "Wrist",
	3: "Finger",
	4: "Hand",
	5: "Ear Lobe",
	6: "Foot",
}

func (b *HeartRateDriver) GetBodySensorLocation() (string, error) {
	c, err := b.adaptor().ReadCharacteristic(cUUIDBodySensorLocation)
	if err != nil {
		return "", err
	}
	val := uint8(c[0])
	ret := mBodySensorLocation[val]
	if ret == "" {
		return "", fmt.Errorf("undefined location %v", val)
	}
	return ret, nil
}

// HeartRateMeasurement flags
var mHeartRateFormat = map[uint8]string{
	0b0: "UINT8",
	0b1: "UINT16",
}
var mSensorContactStatus = map[uint8]string{
	0b00: "not supported",
	0b01: "not supported",
	0b10: "supported but contact is not detected",
	0b11: "supported and contact is detected",
}
var mEnergyExpandedStatus = map[uint8]string{
	0b0: "not present",
	0b1: "present",
}
var mRRInterval = map[uint8]string{
	0b0: "not present",
	0b1: "present (one or more)",
}

type cHRMFlags struct {
	heartRateFormat      uint8
	sensorContactStatus  uint8
	energyExpendedStatus uint8
	rrInterval           uint8
}

func (hrf cHRMFlags) String() string {
	return fmt.Sprintf("HeartRateFormat: %v\n", mHeartRateFormat[hrf.heartRateFormat]) +
		fmt.Sprintf("SensorContactStatus: %v\n", mSensorContactStatus[hrf.sensorContactStatus]) +
		fmt.Sprintf("EnergyExpandedStatus: %v\n", mEnergyExpandedStatus[hrf.energyExpendedStatus]) +
		fmt.Sprintf("RR-Interval: %v", mRRInterval[hrf.rrInterval])
}

func parseHeartRateFlags(flags byte) cHRMFlags {
	var hrf cHRMFlags
	hrf.heartRateFormat = flags & 0b1
	hrf.sensorContactStatus = flags >> 1 & 0b11
	hrf.energyExpendedStatus = flags >> 3 & 0b1
	hrf.rrInterval = flags >> 4 & 0b1
	return hrf
}

func parseHeartRate(data []byte) (heartRate int, hrf cHRMFlags, err error) {
	hrf = parseHeartRateFlags(data[0])
	if hrf.heartRateFormat == 0b0 {
		return int(data[1]), hrf, nil
	} else {
		hr := binary.LittleEndian.Uint16(data[1:3])
		return int(hr), hrf, nil
	}
}

//func (b *HeartRateDriver) SubscribeHeartRate() error {
//	err := b.adaptor().Subscribe(cUUIDHeartRateMeasurement,
//		func(c []byte, _ error) {
//			fmt.Println(time.Now().Format("15:04:05"), c)
//		})
//	return err
//}

func (b *HeartRateDriver) SubscribeHeartRate() error {
	err := b.adaptor().Subscribe(cUUIDHeartRateMeasurement,
		func(data []byte, e error) {
			if e != nil {
				fmt.Fprintf(os.Stderr, "err: %v", e)
				return
			}
			hr, hrf, e := parseHeartRate(data)
			if e != nil {
				fmt.Fprintf(os.Stderr, "err: %v", e)
				return
			}
			fmt.Println(time.Now().Format("15:04:05"))
			fmt.Println(hrf)
			fmt.Printf("HeartRate: %v\n", hr)
			fmt.Println("--------------------")
		})
	return err
}
