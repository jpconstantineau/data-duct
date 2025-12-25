package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jpconstantineau/data-duct/pkg/pipeline"
)

type SensorReading struct {
	ID          int
	Timestamp   time.Time
	Temperature float64
	Moisture    float64
}

type SensorBatchSummary struct {
	StartID     int
	EndID       int
	StartTime   time.Time
	EndTime     time.Time
	Count       int
	AvgTemp     float64
	AvgMoisture float64
	MinTemp     float64
	MaxTemp     float64
	MinMoisture float64
	MaxMoisture float64
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	src := func(ctx context.Context) (<-chan SensorReading, error) {
		ch := make(chan SensorReading)
		go func() {
			defer close(ch)

			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()

			timer := time.NewTimer(10 * time.Second)
			defer timer.Stop()

			id := 0
			for {
				select {
				case <-ctx.Done():
					return
				case <-timer.C:
					return
				case tickTime := <-ticker.C:
					id++

					// Generate stable, deterministic synthetic sensor data.
					t := float64(id)
					temp := 20.0 + 5.0*math.Sin(t*0.15)       // ~15..25
					moist := 50.0 + 15.0*math.Sin(t*0.07+1.0) // ~35..65

					reading := SensorReading{
						ID:          id,
						Timestamp:   tickTime,
						Temperature: temp,
						Moisture:    moist,
					}

					select {
					case <-ctx.Done():
						return
					case ch <- reading:
					}
				}
			}
		}()
		return ch, nil
	}

	batchSummarize := func(ctx context.Context, inputs []SensorReading) ([]SensorBatchSummary, error) {
		if len(inputs) == 0 {
			return nil, nil
		}

		sumTemp := 0.0
		sumMoist := 0.0
		minTemp, maxTemp := inputs[0].Temperature, inputs[0].Temperature
		minMoist, maxMoist := inputs[0].Moisture, inputs[0].Moisture

		for _, r := range inputs {
			sumTemp += r.Temperature
			sumMoist += r.Moisture
			if r.Temperature < minTemp {
				minTemp = r.Temperature
			}
			if r.Temperature > maxTemp {
				maxTemp = r.Temperature
			}
			if r.Moisture < minMoist {
				minMoist = r.Moisture
			}
			if r.Moisture > maxMoist {
				maxMoist = r.Moisture
			}
		}

		summary := SensorBatchSummary{
			StartID:     inputs[0].ID,
			EndID:       inputs[len(inputs)-1].ID,
			StartTime:   inputs[0].Timestamp,
			EndTime:     inputs[len(inputs)-1].Timestamp,
			Count:       len(inputs),
			AvgTemp:     sumTemp / float64(len(inputs)),
			AvgMoisture: sumMoist / float64(len(inputs)),
			MinTemp:     minTemp,
			MaxTemp:     maxTemp,
			MinMoisture: minMoist,
			MaxMoisture: maxMoist,
		}

		// Emit one output per batch: a summary.
		return []SensorBatchSummary{summary}, nil
	}

	batchCount := 0
	sink := func(ctx context.Context, s SensorBatchSummary) error {
		batchCount++
		fmt.Printf(
			"batch=%d ids=%d..%d count=%d avgTemp=%.2f avgMoist=%.2f minTemp=%.2f maxTemp=%.2f minMoist=%.2f maxMoist=%.2f\n",
			batchCount,
			s.StartID,
			s.EndID,
			s.Count,
			s.AvgTemp,
			s.AvgMoisture,
			s.MinTemp,
			s.MaxTemp,
			s.MinMoisture,
			s.MaxMoisture,
		)
		return nil
	}

	runnable := pipeline.New("batch-sensor", src).
		ThenBatch(batchSummarize, pipeline.BatchPolicy{Size: 10}).
		To(sink)

	res, err := runnable.Run(ctx)
	fmt.Printf("result=%s err=%v batches=%d\n", res.State(), err, batchCount)
}
