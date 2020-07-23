package chanhttp

import (
	"net/http"

	"github.com/z2665/chanmonitor/pkg/chanmonitor"
)

type chanhttp struct {
	h    *http.ServeMux
	moni *chanmonitor.ChanMonitor
}

func NewChanHTTP(h *http.ServeMux, m *chanmonitor.ChanMonitor) *chanhttp {
	result := &chanhttp{
		h:    h,
		moni: m,
	}
	tmph := h
	if h == nil {
		tmph = http.DefaultServeMux
	}
	tmph.HandleFunc("/debug/chan/all/json", result.HandleChanInfoJSON)
	tmph.HandleFunc("/debug/chan/overflow/json", result.HandleChanOverFlowInfoJSON)
	tmph.HandleFunc("/debug/chan/overflow/string", result.HandleChanOverFlowInfoString)
	return result
}
func (ch *chanhttp) HandleChanInfoJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	snip := ch.moni.GetSnapshot()
	result, _ := ch.moni.SnapshotToJSON((snip))
	w.Write(result)
}
func (ch *chanhttp) HandleChanOverFlowInfoJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	snip := ch.moni.GetOverFlowSnapshot()
	result, _ := ch.moni.SnapshotToJSON((snip))
	w.Write(result)
}
func (ch *chanhttp) HandleChanOverFlowInfoString(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	snip := ch.moni.GetOverFlowSnapshot()
	result := ch.moni.SnapshotToString(snip)
	for _, v := range result {
		w.Write([]byte(v + "\n"))
	}

}
