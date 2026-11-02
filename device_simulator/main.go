package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"time"
	"uuid"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

type Reading struct {
	Schema_version uint8     `json:"schema_version"`
	Device_id      uuid.UUID `json:"device_id" binding:"required"`
	Metric         string    `json:"metric" binding:"required"`
	Value          float32   `json:"value"`
	Timestamp      time.Time `json:"timestamp" binding:"required"`
	Reliability    uint8     `json:"reliability"`
}

var uuid_devices = [15]uuid.UUID{
	uuid.MustParse("2e91d5a1-4bc7-489f-8c38-0cdb2477ba77"),
	uuid.MustParse("a280262f-3c7a-4603-9452-221f647029cc"),
	uuid.MustParse("376c8797-42e4-412f-81d2-5b2cc8f33316"),
	uuid.MustParse("46a186fe-2097-4cf2-bfe8-855662d2172f"),
	uuid.MustParse("57f5ec1b-2e10-4e67-84ca-6316de7b8fd2"),
	uuid.MustParse("99b2a7da-3f66-42c0-8e69-cc0adb329637"),
	uuid.MustParse("725e16e5-d79e-4f79-8f72-31c6459a4d38"),
	uuid.MustParse("8908dbbd-72a3-42b4-8085-9cb2cdb19d00"),
	uuid.MustParse("95a8c86a-f570-430b-8a23-6c1e5ca0c595"),
	uuid.MustParse("6e5ab456-a208-4ee4-af21-1903711849dc"),
	uuid.MustParse("04c42ee2-b3f9-47bf-a7f1-5a0a4c8d35b0"),
	uuid.MustParse("3f626137-7468-493e-a810-9fa0be7d1aa5"),
	uuid.MustParse("91197646-79fc-474a-bcba-2268694bd079"),
	uuid.MustParse("94df4127-cb7e-4533-a36d-d90c5f974245"),
	uuid.MustParse("6539a354-f905-4bb3-9b6b-56d3fdea3b9e"),
}

var metrics = [3]string{
	"temperature",
	"vibration",
	"current_draw",
}

func roundToPlaces(val float32, places int) float32 {
	shift := math.Pow(10, float64(places))
	return float32(math.Round(float64(val)*shift) / shift)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
	}

	app_address := os.Getenv("APP_ADDRESS")

	devices, err := strconv.Atoi(os.Getenv("DEVICES"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if devices > 15 || devices <= 0 {
		fmt.Println("Max number of devices is 15 and must be grater than 0")
		return
	}

	timeout, err := strconv.Atoi(os.Getenv("DELAY_IN_SECONDS"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if timeout <= 0 {
		fmt.Println("Timeout must be grater than 0")
		return
	}

	for {
		g, ctx := errgroup.WithContext(context.Background())

		var readings []Reading
		timestamp := time.Now()

		for range devices {
			metric := metrics[rand.IntN(3)]
			var value float32

			switch metric {
			case "temperature":
				min, max := 10, 50
				value = roundToPlaces(rand.Float32(), 2) + float32(rand.IntN(max-min+1))
			case "vibration":
				min, max := 0, 100
				value = roundToPlaces(rand.Float32(), 2) + float32(rand.IntN(max-min+1))
			case "current_draw":
				min, max := 0, 220
				value = roundToPlaces(rand.Float32(), 2) + float32(rand.IntN(max-min+1))
			}

			readings = append(
				readings,
				Reading{
					Schema_version: 1,
					Device_id:      uuid_devices[rand.IntN(devices)],
					Metric:         metric,
					Value:          value,
					Timestamp:      timestamp,
					Reliability:    uint8(rand.IntN(100)),
				},
			)
		}

		for i := 0; i < len(readings); i++ {
			reading := readings[i]
			g.Go(func() error {
				jsonData, err := json.Marshal(reading)
				if err != nil {
					return err
				}

				req, err := http.NewRequestWithContext(ctx, "POST", app_address, bytes.NewReader(jsonData))
				if err != nil {
					return err
				}
				req.Header.Set("Content-Type", "application/json")

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				if resp.StatusCode >= 400 {
					return fmt.Errorf("device %s: server returned %d", reading.Device_id, resp.StatusCode)
				}

				return nil
			})
		}

		err := g.Wait()
		if err != nil {
			fmt.Printf("Batch error: %v\n", err.Error())
		}

		time.Sleep(time.Duration(timeout) * time.Second)
	}
}
