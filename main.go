package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"time"
)

const timeLayout = "15:04:05"

// Trains ...
type Trains []Train

// Train ...
type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

// UnmarshalJSON unmarshal a JSON description of Trains type. - коментарі мають починатися з назви функції чи метода.
func (st *Train) UnmarshalJSON(data []byte) error {
	var res struct { // якщо ми використовуємо цю структуру лише один раз - можна ось так зробити - створити змінну неіменованого типу
		TrainID            int
		DepartureStationID int
		ArrivalStationID   int
		Price              float32
		ArrivalTime        string
		DepartureTime      string
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	parsedArrivalTime, err := time.Parse(timeLayout, res.ArrivalTime) // ми формат часу використовуємо декілька разів, тож краще винести в timeLayout константу
	if err != nil {
		return fmt.Errorf("wrong arrival time: %w", err)
	}

	parsedDepartureTime, err := time.Parse(timeLayout, res.DepartureTime)
	if err != nil {
		return fmt.Errorf("wrong departure time: %w", err)
	}

	st.TrainID = res.TrainID
	st.DepartureStationID = res.DepartureStationID
	st.ArrivalStationID = res.ArrivalStationID
	st.Price = res.Price
	st.ArrivalTime = parsedArrivalTime
	st.DepartureTime = parsedDepartureTime

	return nil
}

func main() {
	result, err := FindTrains("1902", "1929", "price")
	if err != nil {
		log.Fatal(err) // після Fatal програма робить os.Exit, що моментально припиняє виконання програми
		// ретурн в такому разі не потрібен
	}

	if len(result) != 0 {
		fmt.Printf("%#v\n", result)
	}
}

// ReadTrainsJson ...
func ReadTrainsJson(pathJson string) (Trains, error) {
	byteValue, err := ioutil.ReadFile(pathJson)
	if err != nil {
		return nil, err
	}

	var trains Trains
	if err := json.Unmarshal(byteValue, &trains); err != nil {
		return nil, err
	}

	return trains, nil
}

// FindTrains ...
func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	trains, err := ReadTrainsJson("data.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read trains data: %w", err)
	}

	if len(trains) <= 0 {
		return nil, errors.New("not enough trains")
	}

	if len(departureStation) <= 0 {
		return nil, errors.New("empty departure station")
	}

	if len(arrivalStation) <= 0 {
		return nil, errors.New("empty arrival station")
	}

	depStName, err := strconv.Atoi(departureStation)
	if err != nil {
		return nil, errors.New("bad departure station input")
	}

	arrStName, err := strconv.Atoi(arrivalStation)
	if err != nil {
		return nil, errors.New("bad arrival station input")
	}

	var result Trains
	for _, tempTrain := range trains {
		if depStName == tempTrain.DepartureStationID && arrStName == tempTrain.ArrivalStationID {
			result = append(result, tempTrain)
		}
	}

	const bestTrainsNum = 3

	if len(result) < bestTrainsNum {
		return nil, errors.New("not enought best trains")
	}

	switch criteria {
	case "price":
		sort.SliceStable(result, func(i, j int) bool { return result[i].Price < result[j].Price })
	case "arrival-time":
		sort.SliceStable(result, func(i, j int) bool { return result[i].ArrivalTime.Before(result[j].ArrivalTime) })
	case "departure-time":
		sort.SliceStable(result, func(i, j int) bool { return result[i].DepartureTime.Before(result[j].DepartureTime) })
	default:
		return nil, errors.New("unsupported criteria")
	}

	return result[0:bestTrainsNum], nil
}
