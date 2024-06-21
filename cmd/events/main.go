package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/juancassiano/golang/api-sell/internal/events/domain"
	httpHandler "github.com/juancassiano/golang/api-sell/internal/events/infra/http"
	"github.com/juancassiano/golang/api-sell/internal/events/infra/repository"
	"github.com/juancassiano/golang/api-sell/internal/events/infra/service"
	"github.com/juancassiano/golang/api-sell/internal/events/usecase"
)

type Data struct {
	Events []domain.Event `json:"events"`
	Spots  []domain.Spot  `json:"spots"`
}

func loadDataFromJSON(filename string) (*Data, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	var data Data
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON: %v", err)
	}

	return &data, nil
}
func main() {
	// db, err := sql.Open("mysql", "test_user:test_password@tcp(localhost:3306/test_db)")
	data, err := loadDataFromJSON("/home/juan/Documentos/projetos/fullcycle/go/data.json")
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	// eventRepo, er := repository.NewMysqlEventRepository(db)
	// if err != nil {
	// 	panic(er)
	// }

	eventRepo := repository.NewInMemoryEventRepository()
	for _, event := range data.Events {
		eventRepo.CreateEvent(&event)
	}
	for _, spot := range data.Spots {
		eventRepo.CreateSpot(&spot)
	}

	partnerBaseURLs := map[int]string{
		1: "http://localhost:9080/api1",
		2: "http://localhost:9080/api2",
	}

	partnerFactory := service.NewPartnerFactory(partnerBaseURLs)

	listEventsUseCase := usecase.NewListEventsUseCase(eventRepo)
	getEventsUseCase := usecase.NewGetEventUseCase(eventRepo)
	listSpotsUseCase := usecase.NewListSpotsUseCase(eventRepo)
	buyTicketsUseCase := usecase.NewBuyTicketsUseCase(eventRepo, partnerFactory)

	eventsHandler := httpHandler.NewEventHandler(
		listEventsUseCase, getEventsUseCase, buyTicketsUseCase, listSpotsUseCase,
	)

	r := http.NewServeMux()
	r.HandleFunc("/events", eventsHandler.ListEvents)
	r.HandleFunc("/events/{eventID}", eventsHandler.GetEvent)
	r.HandleFunc("/events/{eventID}/spots", eventsHandler.ListSpots)
	r.HandleFunc("POST /checkout", eventsHandler.BuyTickets)

	http.ListenAndServe(":8080", r)
}
