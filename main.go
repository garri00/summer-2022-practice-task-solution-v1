package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"time"
)

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

// Переписуємо метод UnmarshalJSON для структури Trains для того щоб спарсити час

func (st *Train) UnmarshalJSON(data []byte) error {
	// Створюємо додаткову структуру для перезапису ArrivalTime та DepartureTime у вигляді рядка
	type parseType struct {
		TrainID            int
		DepartureStationID int
		ArrivalStationID   int
		Price              float32
		ArrivalTime        string
		DepartureTime      string
	}
	var res parseType
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	parsedArrivalTime, err := time.Parse("15:04:05", res.ArrivalTime)
	if err != nil {
		return fmt.Errorf("wrong arrival time: %w", err)
	}
	parsedDepartureTime, err := time.Parse("15:04:05", res.DepartureTime)
	if err != nil {
		return fmt.Errorf("wrong departure time: %w", err)
	}

	// Записуємо в вихідну структуру наш час
	st.TrainID = res.TrainID
	st.DepartureStationID = res.DepartureStationID
	st.ArrivalStationID = res.ArrivalStationID
	st.Price = res.Price
	st.ArrivalTime = parsedArrivalTime
	st.DepartureTime = parsedDepartureTime

	return nil
}

func main() {

	//var departureStation, arrivalStation, criteria string
	//fmt.Println("Введіть номер станції відправлення :")
	//fmt.Scan(&departureStation)
	//fmt.Println("Введіть номер станції прибуття :")
	//fmt.Scan(&arrivalStation)
	//fmt.Println("Введіть ритерій, по котрому треба відсортувати потяги (price, arrival-time, departure-time):")
	//fmt.Scan(&criteria)

	//test cases
	result, err := FindTrains("1902", "1929", "price")

	//result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(result) != 0 {
		fmt.Printf("%#v\n", result)
	}
}

func ReadTrainsJson(pathJson string) Trains {
	var trains Trains
	byteValue, err := ioutil.ReadFile(pathJson)
	if err != nil {
		fmt.Errorf("%w", err)
		return nil
	}
	if err := json.Unmarshal(byteValue, &trains); err != nil {
		fmt.Errorf("%w", err)
		return nil
	}
	return trains
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {

	var trains Trains
	const pathJson = "data.json"
	trains = ReadTrainsJson(pathJson)

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

	var bestTrains Trains
	for _, tempTrain := range trains {
		if depStName == tempTrain.DepartureStationID && arrStName == tempTrain.ArrivalStationID {
			bestTrains = append(bestTrains, tempTrain)
		}
	}

	if len(bestTrains) < 3 {
		return nil, errors.New("not enought best trains")
	}

	switch criteria {
	case "price":
		sort.SliceStable(bestTrains, func(i, j int) bool { return bestTrains[i].Price < bestTrains[j].Price })
	case "arrival-time":
		sort.SliceStable(bestTrains, func(i, j int) bool { return bestTrains[i].ArrivalTime.Before(bestTrains[j].ArrivalTime) })
	case "departure-time":
		sort.SliceStable(bestTrains, func(i, j int) bool { return bestTrains[i].DepartureTime.Before(bestTrains[j].DepartureTime) })
	default:
		return nil, errors.New("unsupported criteria")
	}

	const bestTrainsNum = 3
	return bestTrains[0:bestTrainsNum], nil // маєте повернути правильні значення
}
