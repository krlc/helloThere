package helloThere

import (
	"net"
	"sync"

	gh "github.com/mmcloughlin/geohash"
	"github.com/oschwald/maxminddb-golang"
)

const (
	geohashPrecision = 5
)

type GeoDBReader interface {
	Lookup(net.IP, interface{}) error
	Close() error
}

type GeoDB struct {
	reader       GeoDBReader
	geoDBReqPool *sync.Pool
}

type GeoDBRequest struct {
	City struct {
		Names struct {
			En string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"city"`

	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`

	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
	} `maxminddb:"location"`
}

type GeoDBResult struct {
	country string
	city    string
	geohash string
}

func newGeoDBInstance(reader GeoDBReader) *GeoDB {
	return &GeoDB{
		reader: reader,
		geoDBReqPool: &sync.Pool{
			New: func() interface{} {
				return new(GeoDBRequest)
			},
		},
	}
}

func NewGeoDB(dbPath string) (*GeoDB, error) {
	reader, err := maxminddb.Open(dbPath)
	if err != nil {
		return nil, err
	}

	return newGeoDBInstance(reader), nil
}

func (g *GeoDB) Lookup(ip net.IP) (string, string, string, error) {
	req := g.geoDBReqPool.Get().(*GeoDBRequest)
	if err := g.reader.Lookup(ip, req); err != nil {
		return "", "", "", err
	}
	g.geoDBReqPool.Put(req)

	geohash := gh.EncodeWithPrecision(req.Location.Latitude, req.Location.Longitude, geohashPrecision)
	return req.Country.ISOCode, req.City.Names.En, geohash, nil
}

func (g *GeoDB) Close() error {
	return g.reader.Close()
}
