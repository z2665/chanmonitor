package chanmonitor

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	chanNameExistsError = "chan name is exists"
)

//ChanMonitor is monitor struce
type ChanMonitor struct {
	chans        map[string]interface{}
	mux          *sync.Mutex
	threshold    int
	interval     int
	isrunning    int32
	chansnapshot []Chaninfo
}
type Chaninfo struct {
	Len        int    `json:"len"`
	Cap        int    `json:"cap"`
	Percent    int    `json:"percent"`
	Name       string `json:"name"`
	IsOverflow bool   `json:"isOverflow"`
}

//NewChanMonitor return a new chanMonitor with threshold value adn default interval 10s
func NewChanMonitor(threshold int) *ChanMonitor {
	//default interval is 10s
	return NewChanMonitorWithInterval(threshold, 10)
}

//NewChanMonitorWithInterval return a new chanMonitor with threshold value and interval
func NewChanMonitorWithInterval(threshold int, interval int) *ChanMonitor {
	cm := &ChanMonitor{
		chans:        make(map[string]interface{}),
		mux:          &sync.Mutex{},
		threshold:    threshold,
		interval:     interval,
		isrunning:    1,
		chansnapshot: make([]Chaninfo, 0),
	}
	go cm.run()
	return cm
}

//AddChan is used to add a channel to monitor
func (cm *ChanMonitor) AddChan(channame string, c interface{}) error {
	cm.mux.Lock()
	defer cm.mux.Unlock()
	if _, ok := cm.chans[channame]; ok {
		return errors.New(chanNameExistsError)
	}
	cm.chans[channame] = c
	return nil
}

//Stop the monitor
func (cm *ChanMonitor) Stop() {
	atomic.StoreInt32(&cm.isrunning, 0)
}
func (cm *ChanMonitor) run() {
	isOverflow := false
	for atomic.LoadInt32(&cm.isrunning) != 0 {
		cm.mux.Lock()
		cm.chansnapshot = cm.chansnapshot[:0]
		for k, v := range cm.chans {
			vc := reflect.ValueOf(v)
			//get floor of result
			per := int(math.Floor(float64(vc.Len()) / float64(vc.Cap()) * 100))
			if per >= cm.threshold {
				isOverflow = true
			} else {
				isOverflow = false
			}
			cm.chansnapshot = append(cm.chansnapshot, Chaninfo{
				Len:        vc.Len(),
				Cap:        vc.Cap(),
				Percent:    per,
				Name:       k,
				IsOverflow: isOverflow,
			})
		}
		cm.mux.Unlock()
		time.Sleep(time.Duration(cm.interval) * time.Second)
	}
}

//GetSnapshot get the snapshot of the chan monitor
func (cm *ChanMonitor) GetSnapshot() []Chaninfo {
	cm.mux.Lock()
	defer cm.mux.Unlock()
	result := make([]Chaninfo, 0)
	for _, s := range cm.chansnapshot {
		result = append(result, s)
	}
	return result
}

// GetOverFlowSnapshot get over flow snapshot
func (cm *ChanMonitor) GetOverFlowSnapshot() []Chaninfo {
	snap := cm.GetSnapshot()
	result := make([]Chaninfo, 0)
	for _, s := range snap {
		if s.IsOverflow {
			result = append(result, s)
		}
	}
	return result
}

//SnapshotToString make Chaninfo array to string array
func (cm *ChanMonitor) SnapshotToString(snap []Chaninfo) []string {
	result := make([]string, 0)
	for _, s := range snap {
		if s.IsOverflow {
			result = append(result, fmt.Sprintf("%s overflowd channel len %d cap %d len/cap is over %d", s.Name, s.Len, s.Cap, s.Percent))
		} else {
			result = append(result, fmt.Sprintf("%s channel len %d cap %d len/cap is over %d", s.Name, s.Len, s.Cap, s.Percent))
		}

	}
	return result
}

//SnapshotToJSON make ChainInfo array to json format
func (cm *ChanMonitor) SnapshotToJSON(snap []Chaninfo) ([]byte, error) {
	return jsoniter.Marshal(snap)
}
