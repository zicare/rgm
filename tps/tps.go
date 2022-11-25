// TPS package.
// TPS is an acronym for Transactions Per Second package.
// This package allows for TPS control.
// Users's actual TPS is recalculated on every request and a time penalty
// is accounted for in case the user's TPS rate is exceeded.
// Penalties are accumulative, so if a user who is currently voided insists
// making new requests exceeding her TPS quota, new penalties will add to
// existing ones.
// Penalties are stored in-memory in the TPSmap registry. Obsolete penalties
// are removed from the registry automatically.
package tps

import (
	"time"

	"github.com/zicare/rgm/config"
)

// TPS exported
type TPS struct {

	// A buffered channel, used to calculate the actual TPS.
	// A channel capacity value must be set in the configuration file.
	// The requests' timestamps are sent to the channel.
	// When the channel is full, difference in seconds from
	// oldest to newest timestamps divided by the channel
	// capacity will tell the actual TPS.
	// Once the channel is full, each new request pops out the oldest one
	// before recalculating the TPS and possibly a penalty.
	// The channel capacity must be 3 or more.
	ch chan time.Time

	// Channel more recent update timestamp
	// The newest timestamp in the channel matches this value
	// Used by the cleanUp go-routine.
	ts time.Time

	// The penalty. Requests won't be authorized before this time.
	// A nil value means no penalty.
	op *time.Time
}

// TPSmap exported
type TPSmap map[string]map[string]*TPS

// The map.
// Each acl.User required to authorize a transaction
// have an entry on this map.
// Lifetime of an entry with no new transactions is determined by mclnup.
var tpsmap TPSmap

// Number of requests needed to measure actual TPS.
// This imposses the channel capacity.
var chcap int

// TPS map clean up cycle length
var mclnup time.Duration

// Penalty factor
var pf float32

// Init function initializes TPS control.
// chcap is the number of request needed to calculate TPS.
// After a period of a user's inactivity, her TPS control is reset
// to save server resources. mapcln param sets the time window of
// inactivity required in order to reset. mapcln starts counting
// after any TPS related penalty is fulfilled.
func Init() error {

	chcap = config.Config().GetInt("tps.precision")
	mclnup = config.Config().GetDuration("tps.clean_up_cycle")
	pf = float32(config.Config().GetFloat64("tps.penalty_factor"))

	if (chcap < 3) || (chcap > 10) {
		e := new(PrecisionRange).SetArgs("3", "10")
		return &e
	} else if (mclnup < 1) || (mclnup > 10) {
		e := new(CleanUpCycleRange).SetArgs("1", "10")
		return &e
	} else if (pf < 0) || (pf > 10) {
		e := new(PenaltyFactorRange).SetArgs("0", "10")
		return &e
	} else {
		tpsmap = map[string]map[string]*TPS{}
		cleanUp()
		return nil
	}
}

// Transaction takes note of the request,
// recalculates the actual TPS, and returns a datetime
// if the TPS rate is exceeded. In such case, the user
// shouldn't be granted access before said datetime.
// Even if a user is blocked, new request should also be
// accounted for here, a more distant datetime could be returned.
func Transaction(t string, uid string, tpsMax float32) *time.Time {

	now := time.Now()

	if _, ok := tpsmap[t]; !ok {
		tpsmap[t] = map[string]*TPS{}
	}

	if _, ok := tpsmap[t][uid]; !ok {
		tps := new(TPS)
		tps.ch = make(chan time.Time, chcap)
		tpsmap[t][uid] = tps
	}

	if len(tpsmap[t][uid].ch) == cap(tpsmap[t][uid].ch) {
		tpsmap[t][uid].setOp(<-tpsmap[t][uid].ch, now, tpsMax)
	}

	tpsmap[t][uid].ch <- now
	tpsmap[t][uid].ts = now

	return tpsmap[t][uid].op
}

// Return actual current tps and
// Sets tps.op to block transactions if max tps limit is exceeded
func (tps *TPS) setOp(tch time.Time, now time.Time, tpsMax float32) float32 {

	t := float32(now.UnixNano()-tch.UnixNano()) / float32(time.Second)

	//Actual current tps
	tpsNow := float32(chcap) / t

	if tpsNow > tpsMax {
		if tps.op == nil {
			tps.op = &now
		}
		op := (*tps.op).Add(time.Minute * time.Duration(pf*(tpsNow/tpsMax)))
		tps.op = &op
	}

	return tpsNow
}

// IsEnabled exported.
// Just by calling Init() once, IsEnabled will return true.
func IsEnabled() bool {
	if tpsmap != nil {
		return true
	}
	return false
}

// Remove obsolete entries from the TPSmap registry to save resources.
// The registry is inspected with a frequency determined by mclnup.
// mclnup is the wait in minutes for a new clean up cycle to start
// after the previous cycle ended.
// mclnup must be set in the configuration file.
// An entry is considered obsolete if 2 conditions are met.
// 1. The more recent channel update (TPS.ts), that is the last request,
// was made at least 2 clean up cycles ago.
// 2. The penalty is nil or fulfilled
func cleanUp() {

	go func(m TPSmap) {
		d := mclnup * time.Minute
		for {
			for _, v1 := range m {
				for k2, v2 := range v1 {
					if v2.ts.Before(time.Now().Add(-1*d)) &&
						(v2.op == nil || v2.op.Before(time.Now())) {
						delete(v1, k2)
					}
				}
			}
			time.Sleep(d)
		}
	}(tpsmap)
}
