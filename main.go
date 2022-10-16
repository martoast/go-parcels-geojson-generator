package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type Property struct {
}

type Feature struct {
	Type       string   `json:"type"`
	Geometry   Geometry `json:"geometry"`
	Properties Property `json:"properties"`
}

type Geojson struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Polygon struct {
	Id        string `gorm:""json:"id"`
	Parcel_id string `json:"parcel_id"`
	Geojson   string `json:"geojson"`
}

func main() {
	err := godotenv.Load()
	if err == nil {
		log.Println("Working with local env")
	}
	dbUser := os.Getenv("DB_User")
	dbPass := os.Getenv("DB_Password")
	dbName := os.Getenv("DB_Name")
	dbHost := os.Getenv("DB_Host")
	dbPort := os.Getenv("DB_Port")

	var DSN = " host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort

	log.Println(DSN)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	} else {
		log.Println("DB Connected")
	}

	var Polygons []Polygon

	db.Select("ST_AsGeoJSON(polygons) as geojson, parcel_id as parcel_id, id as id").Limit(10).Find(&Polygons)

	// json.Marshal(Polygons)

	file, err := os.Create("test.geojson")

	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	geojsons := []Geojson{}

	for key, _ := range Polygons {

		data := Geometry{}

		json.Unmarshal([]byte(Polygons[key].Geojson), &data)

		for _index1, _key1 := range data.Coordinates {
			fmt.Println(_key1)
			for _index2 := range _key1 {
				data.Coordinates[_index1][_index2][0], data.Coordinates[_index1][_index2][1] = data.Coordinates[_index1][_index2][1], data.Coordinates[_index1][_index2][0]
			}
		}

		geojson := Geojson{
			Type: "FeatureCollection",
			Features: []Feature{
				{
					Type: "Feature",
					Geometry: Geometry{
						Type:        data.Type,
						Coordinates: data.Coordinates,
					},
					Properties: Property{},
				},
			},
		}

		geojsons = append(geojsons, geojson)

	}

	jsonData, _ := json.MarshalIndent(geojsons, "", "  ")
	file.Write(jsonData)
}