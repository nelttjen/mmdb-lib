package mmdbgeo

import (
	"fmt"
	"net"
)

type GeoData struct {
	Country    string
	CountryISO string
	City       string
	Region     string
	RegionISO  string
}

func GetDataByIP(ip string) (*GeoData, error) {
	if globalConnection == nil {
		return nil, ConnectionNotInitialized
	}

	netip := net.ParseIP(ip)
	if netip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	countryResult, err := globalConnection.Country.Reader.Country(netip)
	if err != nil {
		return nil, err
	}
	cityResult, err := globalConnection.City.Reader.City(netip)
	if err != nil {
		return nil, err
	}
	var regionIso string
	var regionEn string
	if len(cityResult.Subdivisions) > 0 {
		regionIso = cityResult.Subdivisions[0].IsoCode
		regionEn = cityResult.Subdivisions[0].Names["en"]
	}

	return &GeoData{
		Country:    countryResult.Country.Names["en"],
		CountryISO: countryResult.Country.IsoCode,
		City:       cityResult.City.Names["en"],
		RegionISO:  regionIso,
		Region:     regionEn,
	}, nil
}
