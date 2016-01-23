package main

import (
	geo "github.com/oschwald/geoip2-golang"
	"fmt"
	"net"
)

const geolocationHelp = "geolocate an ip address using maxmind. syntax: ip <ip>"

func initMaxmind() {
	var err error
	cfg.Maxmind.Reader, err = geo.Open(cfg.Maxmind.DB)
	if err != nil {
		panic(err)
	}
	cfg.Maxmind.available = true
	return
}

func geolocate(ip string) string {
	if !cfg.Maxmind.available {
		return "maxmind geolocation is not available"
	}
	record, err := cfg.Maxmind.Reader.City(net.ParseIP(ip))
	if err != nil {
		return fmt.Sprintf("maxmind geolocation failed: %v", err)
	}
	out := fmt.Sprintf("%s is located in %s, %s - lat/lon=%.6f,%.6f",
		ip, record.City.Names["en"], record.Country.Names["en"],
		record.Location.Latitude, record.Location.Longitude)
	if record.Traits.IsAnonymousProxy {
		out += " - LISTED AS ANONYMOUS PROXY"
	}
	return out
}
