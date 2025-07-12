package session

import (
	"context"
	"fmt"
	"math"
	"ojparkinson/RaceReplay/internal/db"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

type SessionReplay struct {
	SessionID string
}

func (s *SessionReplay) Start(ws *websocket.Conn) {
	println(s.SessionID)

	conn := db.DB{}

	conn.ConnClient()

	queryAPI := conn.Client.QueryAPI("myorg")

	result, err := queryAPI.Query(context.Background(), `from(bucket: "telemetry_Spa")|> range(start: -365d)|> filter(fn: (r) => r._measurement == "telemetry_ticks")|> filter(fn: (r) => r.lap_id == "1")|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")|> sort(columns: ["session_time"])`)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				// fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// fmt.Printf("row: %s\n", result.Record().String())
			speed := int(math.Round(result.Record().ValueByKey("speed").(float64)))
			fmt.Printf("speed: %d\n", speed)
			ws.Write([]byte(strconv.Itoa(speed)))
			time.Sleep(20 * time.Millisecond)
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
}
