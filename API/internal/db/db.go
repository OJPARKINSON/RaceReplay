package db

import (
	"crypto/tls"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type DB struct {
	Client influxdb2.Client
}

func (db *DB) ConnClient() {
	db.Client = influxdb2.NewClientWithOptions("http://localhost:8086", "",
		influxdb2.DefaultOptions().
			SetUseGZip(true).
			SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true,
			}))
}
