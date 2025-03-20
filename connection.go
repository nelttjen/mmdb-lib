package mmdbgeo

import (
	"context"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"time"
)

var globalConnection *FullConnection
var ConnectionNotInitialized = fmt.Errorf("MMDB connection not initialized")

type Connection struct {
	Reader *geoip2.Reader

	reconectInterval time.Duration
	connectedAt      time.Time
}

func (c *Connection) Init(ctx context.Context, path string, interval time.Duration) {
	c.reconectInterval = interval
	err := c.Connect(path)
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize ASN connection: %v\n", err)
	}

	go func() {
		ticker := time.NewTicker(c.reconectInterval)
		defer ticker.Stop()

	reconnectloop:
		for {
			select {
			case <-ticker.C:
				err := c.Connect(path)
				if err != nil {
					log.Printf("[ERROR] Failed to reconnect connection: %v\n", err)
				}
			case <-ctx.Done():
				break reconnectloop
			}
		}
		log.Println("Exiting connection reconnector")
	}()
}

func (c *Connection) Connect(path string) error {
	if c.Reader != nil {
		c.Reader.Close()
		c.Reader = nil
	}

	db, err := geoip2.Open(path)
	if err != nil {
		log.Printf("[ERROR] Failed to open GeoLite2 database (%s): %v\n", path, err)
		return err
	}
	c.Reader = db
	c.connectedAt = time.Now()
	return nil
}

type FullConnection struct {
	Asn     *Connection
	Country *Connection
	City    *Connection
}

func (fc *FullConnection) Init(ctx context.Context, interval time.Duration) {
	asnPath := fmt.Sprintf("%s/%s/%s", rootDir, DBsPath, MMDBASNName)
	countryPath := fmt.Sprintf("%s/%s/%s", rootDir, DBsPath, MMDBCountryName)
	geoPath := fmt.Sprintf("%s/%s/%s", rootDir, DBsPath, MMDBCityName)

	fc.Asn = &Connection{}
	fc.Country = &Connection{}
	fc.City = &Connection{}

	fc.Asn.Init(ctx, asnPath, interval)
	fc.Country.Init(ctx, countryPath, interval)
	fc.City.Init(ctx, geoPath, interval)
}

func GlobalInit(ctx context.Context, interval time.Duration) {
	globalConnection = &FullConnection{}
	globalConnection.Init(ctx, interval)
}
