package nest

import (
	"fmt"
	"time"

	"github.com/robstrong/nest-bt/btpoll"
)

var (
	watcherDuration     = 5 * time.Second
	awayEnabledDuration = 20 * time.Second
)

type NestHandler struct {
	deviceStatus        map[string]btpoll.DeviceStatus
	devicesMissingSince *time.Time
	away                bool
}

func NewNestHandler() *NestHandler {
	h := &NestHandler{
		deviceStatus: map[string]btpoll.DeviceStatus{},
	}
	go h.startWatcher()
	return h
}

func (n *NestHandler) Found(mac string) {
	if n.deviceStatus[mac] != btpoll.DeviceStatusNearby {
		n.deviceStatus[mac] = btpoll.DeviceStatusNearby
		n.statusChange(mac)
	}
}

func (n *NestHandler) NotFound(mac string) {
	if n.deviceStatus[mac] != btpoll.DeviceStatusNotFound {
		n.deviceStatus[mac] = btpoll.DeviceStatusNotFound
		n.statusChange(mac)
	}
}

func (n *NestHandler) statusChange(mac string) {
	fmt.Printf("device %s status changed to %s\n", mac, n.deviceStatus[mac])
	for _, s := range n.deviceStatus {
		if s == btpoll.DeviceStatusNearby {
			n.devicesMissingSince = nil
			return
		}
	}
	now := time.Now()
	n.devicesMissingSince = &now
}

func (n *NestHandler) startWatcher() {
	for _ = range time.Tick(watcherDuration) {
		if n.devicesMissingSince != nil && time.Since(*n.devicesMissingSince) > awayEnabledDuration {
			n.startAwayMode()
		} else {
			n.endAwayMode()
		}
	}
}

func (n *NestHandler) endAwayMode() {
	if !n.awayModeEnabled() {
		return
	}
	n.away = false
	fmt.Println("end away mode")
}

func (n *NestHandler) startAwayMode() {
	if n.awayModeEnabled() {
		return
	}
	n.away = true
	fmt.Println("enable away mode")
}

func (n *NestHandler) awayModeEnabled() bool {
	return n.away
}
