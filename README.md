# ChanMonitor
ChanMonitor is a lib for watch golang channel's load

to know more information see the example

we also provide the http server to get the result

use``/debug/chan/all/json`` to get all watched channels with json format
```json
[{"len":7,"cap":10,"percent":70,"name":"example-1","isOverflow":true},{"len":0,"cap":10,"percent":0,"name":"example-2","isOverflow":false}]
```
use ``/debug/chan/overflow/json`` to get overflow channels with json format
```json
[{"len":7,"cap":10,"percent":70,"name":"example-1","isOverflow":true}]
```
use ``/debug/chan/overflow/string`` to get overflow channels with string format
```
example-1 overflowd channel len 7 cap 10 len/cap is over 70
```
## Usage

1. import the lib in go.mod
```
require github.com/z2665/chanmonitor v0.1.0
```
2. import the pkg and add channels to monitor
```golang
import "github.com/z2665/chanmonitor/pkg/chanmonitor"
mon := chanmonitor.NewChanMonitor(50)
c := make(chan int, 10)
mon.AddChan("example-1", c)
```
3. you can use GetSnapshot or GetOverFlowSnapshot get the channels array
```golang
tmp:=mon.GetSnapshot()
tmp=mon.GetOverFlowSnapshot()
```
4. do your own business

5. if you want to set  the monitor's interval time,use the ``NewChanMonitorWithInterval`` to create the monitor

