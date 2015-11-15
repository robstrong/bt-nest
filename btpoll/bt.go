package btpoll

import (
	"errors"
	"fmt"
	"os/exec"
	"time"
)

//f is the frequency to check for devices
func New(f time.Duration) *Poller {
	return &Poller{
		frequency: f,
		devices:   map[string]*Device{},
	}
}

type Poller struct {
	frequency time.Duration
	devices   map[string]*Device //devices to poll for
}

func (p *Poller) AddDevice(mac string, onFound DeviceEvent, onMissing DeviceEvent) {
	p.devices[mac] = &Device{
		mac:       mac,
		onFound:   onFound,
		onMissing: onMissing,
	}
}

type Device struct {
	mac           string //mac address of device
	currentStatus DeviceStatus
	onFound       func(string)
	onMissing     func(string)
}

func (d *Device) checkStatus() error {
	cmd := exec.Command("bt-device", "-s", d.mac)
	err := cmd.Run()
	switch err.(type) {
	case *exec.ExitError:
		d.deviceNotFound()
	case nil:
		d.deviceFound()
	default:
		return err
	}
	return nil
}

func (d *Device) setStatus(s DeviceStatus) {
	if d.currentStatus == s {
		return
	}
	d.currentStatus = s
}

func (d *Device) deviceFound() {
	fmt.Println("found")
	d.setStatus(DeviceStatusNearby)
	d.onFound(d.mac)
}

func (d *Device) deviceNotFound() {
	fmt.Println("not found")
	d.setStatus(DeviceStatusNotFound)
	d.onMissing(d.mac)
}

type DeviceEvent func(string)

type DeviceStatus string

const (
	DeviceStatusUnknown  DeviceStatus = ""
	DeviceStatusNearby                = "nearby"
	DeviceStatusNotFound              = "missing"
)

//this is a blocking function
func (p *Poller) Start() error {
	//make sure a device has been added
	if len(p.devices) == 0 {
		return errors.New("btpoll: no devices added")
	}

	//make sure this is running as root
	//have to do it this way instead of user package because
	//of cross compiling issue
	cmd := exec.Command("whoami")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	if string(out) != "root\n" {
		return fmt.Errorf("btpoll: must run as root not %s", out)
	}

	//start listener
	for _ = range time.Tick(p.frequency) {
		for _, d := range p.devices {
			err := d.checkStatus()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
