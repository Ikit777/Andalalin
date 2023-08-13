package utils

import (
	"embed"
	"time"
)

var zoneInfoFS embed.FS // get it from GOROOT/lib/time/zoneInfo.zip

func GetLocation(name string) (loc *time.Location) {
	bs, err := zoneInfoFS.ReadFile("zoneinfo/" + name)
	if err != nil {
		panic(err)
	}
	loc, err = time.LoadLocationFromTZData(name, bs)
	if err != nil {
		panic(err)
	}
	return loc
}
