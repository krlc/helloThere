package helloThere

import (
	"errors"
	"net"
	"testing"
)

const (
	testCountry = "LV"
	testCity    = "Riga"
	testLat     = 56.946285
	testLong    = 24.105078
	testGeohash = "ud15u"
)

type dbMock struct{}

func (db *dbMock) Lookup(_ net.IP, req interface{}) error {
	req.(*GeoDBRequest).Country.ISOCode = testCountry
	req.(*GeoDBRequest).City.Names.En = testCity
	req.(*GeoDBRequest).Location.Latitude = testLat
	req.(*GeoDBRequest).Location.Longitude = testLong
	return nil
}

func (db *dbMock) Close() error {
	return nil
}

type dbMockErr struct{}

func (db *dbMockErr) Lookup(_ net.IP, _ interface{}) error {
	return errors.New("error")
}

func (db *dbMockErr) Close() error {
	return errors.New("error")
}

func TestGeoIP_Lookup(t *testing.T) {
	geoip := newGeoDBInstance(new(dbMock))

	country, city, geohash, err := geoip.Lookup(nil)
	if err != nil {
		t.Error(err)
	}

	if country != testCountry || city != testCity || geohash != testGeohash {
		t.Errorf("Lookup() got result = {%s, %s, %s}, want {%s, %s, %s}",
			testCountry, testCity, testGeohash, country, city, geohash)
	}
}

func TestGeoIP_LookupError(t *testing.T) {
	geoipErr := newGeoDBInstance(new(dbMockErr))

	_, _, _, err := geoipErr.Lookup(nil)
	if err == nil {
		t.Error("Lookup() got a nil error")
	}

	if err := geoipErr.Close(); err == nil {
		t.Error("Close() got a nil error")
	}
}

func TestGeoIP_createGeoIP(t *testing.T) {
	if _, err := NewGeoDB(""); err == nil {
		t.Errorf("NewGeoDB() got error = %v, want nil", err)
	}
}

func BenchmarkGeoIP_Lookup(b *testing.B) {
	geoip := newGeoDBInstance(new(dbMock))

	for n := 0; n < b.N; n++ {
		if _, _, _, err := geoip.Lookup(nil); err != nil {
			b.Error(err)
		}
	}
}
