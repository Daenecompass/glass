package timer

import (
	"encoding/json"
	"fmt"
	"time"
)

//A timeline holds a number of lines
//for a single file
type Timeline struct {
	Edges [][]time.Time `json:"edges"`
}

func NewTimeline() *Timeline {
	return &Timeline{}
}

func (tl *Timeline) Length() time.Duration {
	res := time.Millisecond * 0
	for _, e := range tl.Edges {
		if len(e) == 0 {
			continue
		}

		first := e[0]
		last := e[len(e)-1]
		res += last.Sub(first)
	}

	return res
}

func (tl *Timeline) OpenAt(t time.Time) {
	tl.Edges = append(tl.Edges, []time.Time{t})
}
func (tl *Timeline) CloseAt(t time.Time) {
	tl.ProgressTo(t)
}

func (tl *Timeline) ProgressTo(t time.Time) {
	if len(tl.Edges) == 0 {
		return
	}

	tl.Edges[len(tl.Edges)-1] = append(tl.Edges[len(tl.Edges)-1], t)
}

//A distributor managed various files
//for a single timer using timelines
type Distributor struct {
	data *distrData
}

type distrData struct {
	ActiveFile string               `json:"active_file"`
	Timelines  map[string]*Timeline `json:"timelines"`
}

var OverheadTimeline = "__overhead"

func NewDistributor() *Distributor {
	d := &Distributor{
		data: &distrData{
			Timelines: map[string]*Timeline{},
		},
	}

	d.init()
	return d
}

func (d *Distributor) init() {}

//close the open timeline
func (d *Distributor) Break(t time.Time) {
	if atl, ok := d.data.Timelines[d.data.ActiveFile]; ok {
		atl.CloseAt(t)
	}
}

//register a new point on the timeline
func (d *Distributor) Register(fpath string, t time.Time) {
	if fpath == "" {
		fpath = OverheadTimeline
	}

	var tl *Timeline
	var ok bool
	if tl, ok = d.data.Timelines[fpath]; !ok {
		tl = NewTimeline()
		d.data.Timelines[fpath] = tl
	}

	var atl *Timeline
	if atl, ok = d.data.Timelines[d.data.ActiveFile]; !ok {
		atl = nil
	}

	if atl != tl {
		if atl != nil {
			atl.CloseAt(t)
		}

		tl.OpenAt(t)
	} else {
		atl.ProgressTo(t)
	}

	d.data.ActiveFile = fpath
}

//extract the time spent on a file from the first point a given point in time
func (d *Distributor) Extract(fpath string, upto time.Time) (time.Duration, error) {
	if fpath == "" {
		fpath = OverheadTimeline
	}

	if tl, ok := d.data.Timelines[fpath]; !ok {
		return 0, fmt.Errorf("No known timeline for file '%s'", fpath)
	} else {
		return tl.Length(), nil
	}
}

func (d *Distributor) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &d)
	if err != nil {
		return err
	}

	d.init()
	return nil
}

func (d *Distributor) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.data)
}
