package session

import (
	"context"
	"fmt"
	"math"
	"ojparkinson/RaceReplay/internal/db"
	"strconv"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"golang.org/x/net/websocket"
)

type SessionReplay struct {
	SessionID       string
	SpeedMultiplier float64
}

func NewSessionReplay(SessionID string) *SessionReplay {
	return &SessionReplay{
		SessionID:       SessionID,
		SpeedMultiplier: 1.0,
	}
}

type DataPoint struct {
	Speed     int
	Gear      int64
	RPM       int
	Voltage   float64
	FuelLevel float64
	TickTime  time.Time
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

	go s.loadSubSets(conn, dataChannel)

	s.StreamData(ws, dataChannel)
}

func (s *SessionReplay) loadSubSets(conn db.DB, dataChannel chan<- DataSubSet) {
	defer close(dataChannel)

	queryAPI := conn.Client.QueryAPI("myorg")

	for lapID := 1; lapID <= 1; lapID++ {
		subset := s.loadSubSet(queryAPI, strconv.Itoa(lapID))

		dataChannel <- subset
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

		speed := int(math.Round(result.Record().ValueByKey("speed").(float64)))
		gear := result.Record().ValueByKey("gear").(int64)
		rpm := int(math.Round(result.Record().ValueByKey("rpm").(float64)))
		// fmt.Printf("recv: %+v", result.Record().Values())
		volatage := result.Record().Values()["voltage"].(float64)
		fuelLevel := result.Record().Values()["fuel_level"].(float64)
		tickTime := result.Record().Time()
		data = append(data, DataPoint{Speed: speed, TickTime: tickTime, Gear: gear, RPM: rpm, Voltage: volatage, FuelLevel: fuelLevel})
	}

	if result.Err() != nil {
		fmt.Printf("Query processing error for lap %s: %s\n", lapID, result.Err().Error())
	}

	fmt.Printf("Loaded subset %s with %d data points\n", lapID, len(data))
	return DataSubSet{LapID: lapID, Data: data}
}

func (s *SessionReplay) StreamData(ws *websocket.Conn, dataChannel <-chan DataSubSet) {
	var lastTickTime *time.Time

	for subset := range dataChannel {
		for i, dataPoint := range subset.Data {
			err := websocket.JSON.Send(ws, dataPoint)
			// ws.Write([]byte(strconv.Itoa(dataPoint.Speed)))
			if err != nil {
				fmt.Println("Failed to write data: " + err.Error())
			}

			var sleepDuration time.Duration

			if i < len(subset.Data)-1 {
				nextTick := subset.Data[i+1].TickTime
				currenTick := dataPoint.TickTime
				sleepDuration = nextTick.Sub(currenTick)
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
				time.Sleep(sleepDuration)
			}

			lastTickTime = &dataPoint.TickTime
		}
	}
}
