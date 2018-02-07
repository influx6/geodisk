package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"sort"

	"github.com/influx6/faux/flags"
)

const (
	earthRadius                 = 6371 // in kilometers
	housingAnywhereGeoLatitude  = 51.925146
	housingAnywhereGeoLongitude = 4.478617
)

var (
	housingAnywhereGeoLongitudeRadians = toRadians(4.478617)
	housingAnywhereGeoLatitudeRadians  = toRadians(51.925146)
)

// errors ...
var (
	ErrInvalidCSVFormat = errors.New("csv data has invalid format, expects 3 per line")
	ErrInvalidGeoHeader = errors.New("csv has invalid geo header or has no header")
)

//**************************************************************
// CSV Geo-Record Methods
//**************************************************************

// GeoRecords defines a slice type for GeoRecords.
type GeoRecords []GeoRecord

// Swap implements the sort.Swap interface and swaps giving
// index within slice.
func (g GeoRecords) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

// Len returns the length of the slice.
func (g GeoRecords) Len() int {
	return len(g)
}

// Less implements hte sort.LessThan interface and swaps
// validates less than state of two records at index i and j by
// checking their distance to each other.
func (g GeoRecords) Less(i, j int) bool {
	return g[i].Dist < g[j].Dist
}

// GeoRecord embodies data stored in expected csv where it contains
// data in format of `"id","lat","long"`. Where the ID represents the giving
// associated ID of geographical location with respective geographical
// coordinates.
type GeoRecord struct {
	ID   string
	Lat  float64
	Long float64
	Dist float64
}

// distanceWithCSVFile attempts to load csv file from provided target path
// calculating distance of each record from giving geo-coordinates of
// latitude and longitude pairs (which must be in radians).
func distanceWithCSVFile(target string, targetLat float64, targetLong float64) ([]GeoRecord, error) {
	targetFile, err := os.Open(target)
	if err != nil {
		return nil, err
	}

	defer targetFile.Close()

	return distanceWithCSVReader(bufio.NewReader(targetFile), targetLat, targetLong)
}

// distanceWithCSVReader attempts to load csv data from provided io.Reader,
// calculating distance of each record from giving geo-coordinates of
// latitude and longitude pairs (which must be in radians).
//
// If CSV file is read and headers are validated to be correct, then
// code moves to read file line by line, it expects each line to have a
// maximum of 3 items, which it then converts into GeoRecord struct.
// If any line contains more than wanted max length or if a lines latitdude
// or longitude is not a valid float64 number, then it stops
// returning all collected records and error.
func distanceWithCSVReader(target io.Reader, targetLat float64, targetLong float64) ([]GeoRecord, error) {
	csvReader := csv.NewReader(target)

	// Read header of csv, if things fail, then return error
	// validate header matches "
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	// if headers are not 3 in total, then return format error.
	if len(header) != 3 {
		return nil, ErrInvalidCSVFormat
	}

	// if headers don't match expected, then return invalid header error.
	if header[0] != "id" || header[1] != "lat" || header[2] != "lng" {
		return nil, ErrInvalidGeoHeader
	}

	var records []GeoRecord

	for {
		line, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return records, err
		}

		if len(line) != 3 {
			return records, ErrInvalidCSVFormat
		}

		// parse latitude value which are expected to
		// be in degrees to radians.
		lat, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return records, err
		}

		// parse longitude value which are expected to
		// be in degrees to radians.
		long, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return records, err
		}

		var record GeoRecord
		record.ID = line[0]
		record.Lat = toRadians(lat)
		record.Long = toRadians(long)
		record.Dist = greatCircleDistance(record.Lat, record.Long, targetLat, targetLong)

		records = append(records, record)
	}

	return records, nil
}

// greatCircleDistance calculates the great-circle distance over a spherical domain (eg earth)
// for the distance between two points on the sphere. It uses the haversine method.
func greatCircleDistance(lat1, long1, lat2, long2 float64) float64 {
	latDiff := lat2 - lat1
	longDiff := long2 - long1
	latDiffMid := latDiff / 2
	longDiffMid := longDiff / 2

	latMidSin := math.Sin(latDiffMid)
	longMidSin := math.Sin(longDiffMid)

	a := (latMidSin * latMidSin) +
		(math.Cos(lat1)*math.Cos(lat2))*(longMidSin*longMidSin)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func toRadians(t float64) float64 {
	return (t * math.Pi) / 180
}

//**************************************************************
// CLI methods
//**************************************************************

func geoDistanceWithDB(ctx flags.Context) error {
	fmt.Fprintln(os.Stderr, "DB command not available yet.")
	return nil
}

func geoDistanceWithCSV(ctx flags.Context) error {
	csvFile, _ := ctx.GetString("file")
	if csvFile == "" {
		// if arguments is not empty, then take first value has file name
		// else return error.
		args := ctx.Args()
		if len(args) == 0 {
			return errors.New("require csv file path, see geodisk csv help")
		}

		csvFile = args[0]
	}

	csvFile = filepath.Clean(csvFile)
	records, err := distanceWithCSVFile(csvFile, housingAnywhereGeoLatitudeRadians, housingAnywhereGeoLongitudeRadians)
	if err != nil {
		return err
	}

	sort.Sort(GeoRecords(records))

	var top5, bottom5 []GeoRecord

	recLen := len(records)
	if recLen <= 5 {
		top5 = records
		bottom5 = records
	}

	if recLen > 5 {
		top5 = records[:5]
		bottom5 = records[recLen-5:]
	}

	fmt.Fprintln(os.Stdout, "Top 5 Locations closest to Housing Anywhere:")
	for _, rec := range top5 {
		fmt.Fprintf(os.Stdout, "\tLocationID: %s (%.6f kilometers)\n", rec.ID, rec.Dist)
	}

	fmt.Println("")
	fmt.Fprintln(os.Stdout, "Top 5 Locations farthest to Housing Anywhere:")
	for _, rec := range bottom5 {
		fmt.Fprintf(os.Stdout, "\tLocationID: %s (%.6f kilometers)\n", rec.ID, rec.Dist)
	}

	fmt.Println("")
	return nil
}

func main() {
	flags.Run("geodisk",
		flags.Command{
			Name:      "csv",
			Action:    geoDistanceWithCSV,
			Usages:    []string{"geodisk -csv.file=./static/geoData.csv csv", "geodisk csv ./static/geoData.csv"},
			ShortDesc: "Calculate geo distance from csv file.",
			Desc:      "Calculates geo-distance of HousingAnywhere from a database of geo coordinates",
			Flags: []flags.Flag{
				&flags.StringFlag{
					Name:    "file",
					Desc:    "csvfile to be used for calculation",
					Default: "",
				},
			},
		},
		flags.Command{
			Name:      "db",
			ShortDesc: "Calculate geo distance from a db.",
			Action:    geoDistanceWithDB,
			Desc:      "Calculates geo-distance of HousingAnywhere from a database (mongo or sql) of geo coordinates.",
			Flags: []flags.Flag{
				&flags.StringFlag{
					Name:    "config",
					Desc:    "config.yaml file that contains database configuration values",
					Default: "config.yaml",
				},
			},
		})
}
