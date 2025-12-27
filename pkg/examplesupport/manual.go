package examplesupport

import "time"

// ManualRunNow returns a TriggerEvent when a user explicitly requests a run-now trigger.
//
// This helper keeps cmd examples consistent and easy to unit test.
func ManualRunNow(runNow bool, now func() time.Time) (TriggerEvent, bool) {
	if !runNow {
		return TriggerEvent{}, false
	}
	if now == nil {
		now = time.Now
	}
	firedAt := now()
	return TriggerEvent{
		Kind:      "manual",
		Occurred:  firedAt,
		SourceRef: "manual://run-now",
	}, true
}
