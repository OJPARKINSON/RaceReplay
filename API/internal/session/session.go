package session

import (
	"context"
	"fmt"
	"ojparkinson/RaceReplay/internal/db"
	"strconv"
	"sync"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"golang.org/x/net/websocket"
)

type SessionReplay struct {
	SessionID       string
	SpeedMultiplier float64
	stopChan        chan struct{}
	mu              sync.Mutex
	stopped         bool
}

func NewSessionReplay(SessionID string) *SessionReplay {
	return &SessionReplay{
		SessionID:       SessionID,
		SpeedMultiplier: 1.0,
		stopChan:        make(chan struct{}),
		stopped:         false,
	}
}

type DataPoint struct {
	Speed          int
	Gear           int
	RPM            int
	WaterTemp      int
	Voltage        float64
	FuelLevel      float64
	LapID          int
	CurrentLapTime string
	DeltaToBestLap string
	LastLapTime    float64
	LFpressure     float64
	RFpressure     float64
	RRpressure     float64
	LRpressure     float64
	LFtempCM       float64
	RFtempCM       float64
	LRtempCM       float64
	RRtempCM       float64
	TickTime       time.Time
}

type DataSubSet struct {
	LapID string
	Data  []DataPoint
}

func (s *SessionReplay) Start(ws *websocket.Conn) {
	println(s.SessionID)

	conn := db.DB{}
	conn.ConnClient()

	dataChannel := make(chan DataSubSet, 2)

	go s.monitorConnection(ws)

	go s.loadSubSets(conn, dataChannel)

	s.StreamData(ws, dataChannel)
}

func (s *SessionReplay) loadSubSets(conn db.DB, dataChannel chan<- DataSubSet) {
	defer close(dataChannel)

	queryAPI := conn.Client.QueryAPI("myorg")

	for lapID := 1; lapID <= 1; lapID++ {
		subset := s.loadSubSet(queryAPI, strconv.Itoa(lapID))

		select {
		case dataChannel <- subset:
			// Successfully sent
		case <-s.stopChan:
			fmt.Println("Stopping data loading due to stop signal")
			return
		}
	}
}

func (s *SessionReplay) loadSubSet(queryAPI api.QueryAPI, lapID string) DataSubSet {
	query := fmt.Sprintf(`from(bucket: "telemetry_Spa")
		|> range(start: -365d)
		|> filter(fn: (r) => r._measurement == "telemetry_ticks")
		|> filter(fn: (r) => r.lap_id == "%s")
		|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
		|> sort(columns: ["session_time"])`, lapID)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Query error for lap %s: %s\n", lapID, err.Error())
	}

	var data []DataPoint
	for result.Next() {
		if result.TableChanged() {
			// Handle table metadata
		}

		record := result.Record()
		values := record.Values()

		rpm := GetIntValue(values["rpm"])
		gear := GetIntValue(values["gear"])
		speed := GetIntValue(values["speed"])
		lapId := GetIntValue(values["lap_id"])
		volatage := GetFloatValue(values["voltage"], 1)
		waterTemp := GetIntValue(values["waterTemp"])
		fuelLevel := GetFloatValue(values["fuel_level"], 1)
		currentLapTime := GetTimeFormattedWithMillis(values["lapCurrentLapTime"])
		deltaToBestLap := GetTimeFormattedWithMillis(values["lapDeltaToBestLap"])
		lastLapTime := GetFloatValue(values["lapLastLapTime"], 2)
		lFpressure := GetPressureInBar(values["lFpressure"], 2)
		rFpressure := GetPressureInBar(values["rFpressure"], 2)
		rRpressure := GetPressureInBar(values["rRpressure"], 2)
		lRpressure := GetPressureInBar(values["lRpressure"], 2)
		lFtempCM := GetFloatValue(values["lFtempCM"], 0)
		rFtempCM := GetFloatValue(values["rFtempCM"], 0)
		lRtempCM := GetFloatValue(values["lRtempCM"], 0)
		rRtempCM := GetFloatValue(values["rRtempCM"], 0)

		tickTime := record.Time()

		dataPoint := DataPoint{
			Speed:          speed,
			TickTime:       tickTime,
			Gear:           gear,
			RPM:            rpm,
			Voltage:        volatage,
			FuelLevel:      fuelLevel,
			WaterTemp:      waterTemp,
			LapID:          lapId,
			CurrentLapTime: currentLapTime,
			DeltaToBestLap: deltaToBestLap,
			LastLapTime:    lastLapTime,
			LFpressure:     lFpressure,
			RFpressure:     rFpressure,
			RRpressure:     rRpressure,
			LRpressure:     lRpressure,
			LFtempCM:       lFtempCM,
			RFtempCM:       rFtempCM,
			LRtempCM:       lRtempCM,
			RRtempCM:       rRtempCM,
		}

		data = append(data, dataPoint)
	}

	if result.Err() != nil {
		fmt.Printf("Query processing error for lap %s: %s\n", lapID, result.Err().Error())
	}

	fmt.Printf("Loaded subset %s with %d data points\n", lapID, len(data))
	return DataSubSet{LapID: lapID, Data: data}
}

func (s *SessionReplay) StreamData(ws *websocket.Conn, dataChannel <-chan DataSubSet) {
	defer s.Stop()

	var lastTickTime *time.Time

	for subset := range dataChannel {
		if s.isStopped() {
			fmt.Println("Stopping stream due to disconnection")
			return
		}

		fmt.Printf("Streaming lap %s with %d data points\n", subset.LapID, len(subset.Data))

		for i, dataPoint := range subset.Data {
			if s.isStopped() {
				fmt.Println("Stopping stream due to disconnection")
				return
			}

			done := make(chan error, 1)
			go func() {
				done <- websocket.JSON.Send(ws, dataPoint)
			}()

			select {
			case err := <-done:
				if err != nil {
					fmt.Printf("Failed to write data: %v\n", err)
					s.Stop()
					return
				}
			case <-time.After(5 * time.Second):
				fmt.Println("Send timeout - client may be unresponsive")
				s.Stop()
				return
			case <-s.stopChan:
				fmt.Println("Stopping stream due to stop signal")
				return
			}

			var sleepDuration time.Duration

			if i < len(subset.Data)-1 {
				nextTick := subset.Data[i+1].TickTime
				currentTick := dataPoint.TickTime
				sleepDuration = nextTick.Sub(currentTick)
			} else if lastTickTime != nil {
				sleepDuration = dataPoint.TickTime.Sub(*lastTickTime)
			}

			if sleepDuration > 5*time.Second {
				sleepDuration = 5 * time.Second
			} else if sleepDuration < time.Millisecond {
				sleepDuration = time.Millisecond
			}

			sleepDuration = time.Duration(float64(sleepDuration) / s.SpeedMultiplier)

			if sleepDuration > 0 {
				select {
				case <-time.After(sleepDuration):
				case <-s.stopChan:
					fmt.Println("Sleep interrupted due to stop signal")
					return
				}
			}

			lastTickTime = &dataPoint.TickTime
		}

		fmt.Printf("Finished streaming lap %s\n", subset.LapID)
	}

	fmt.Println("Stream completed normally")
}

func (s *SessionReplay) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.stopped {
		s.stopped = true
		close(s.stopChan)
	}
}

func (s *SessionReplay) isStopped() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopped
}

func (s *SessionReplay) monitorConnection(ws *websocket.Conn) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if s.isStopped() {
				return
			}

			if err := websocket.Message.Send(ws, `{"type":"ping"}`); err != nil {
				fmt.Printf("Connection lost (ping failed): %v\n", err)
				s.Stop()
				return
			}
		}
	}
}
