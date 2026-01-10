package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"net/http"
	"strings"
	"time"
	"scheduler/internal/models"

	"github.com/nats-io/nats.go"
	"github.com/zhashkevych/scheduler"
)

var jsc nats.JetStreamContext

func main() {

	nc, err := nats.Connect(nats.DefaulURL)
	if err != nil{
		log.Fatal("Impossible de se connecter a NATS:", err)
	}
	defer nc.Close()

	jsc, err = nc.JetStream()

	if err != nil{
		log.Fatal(err)
	}

	_, _ = jsc.AddStream(&nats.StreamConfig{
		Name: "COURSES",
		Subjects: []string{"COURSES.>"},
	})


	ctx := context.Background()
	sc := scheduler.NewScheduler()

	sc.Add(ctx, runJob, time.Minute*1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	
	<-quit
	sc.Stop()

}

func runJob(ctx context.Context){
	agendas, err := fetchAgendas("http://localhost:8080/agendas/")

	if err != nil {
		fmt.Printf("Erreur lors de la recuperation des agendas: %v\n", err)
		return
	}

	for _, a := range agendas{
		url := fmt.Sprintf("https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=%d&projectId=3&calType=ical&nbWeeks=8&displayConfigId=128", a.UcaId)

		courses, err := parseICal(url, a.Name)

		if err != nil{
			fmt.Printf("Erreur pour l'agenda %s: %v\n", a.Name, err)
			continue
		}

		for _, c := range courses{
			publishToNats(c)
		}



	}
}


func publishToNats(course models.Course){
	messageBytes, err := json.Marshal(course)
	if err != nil{
		fmt.Println("Erreur JSON:", err)
		return
	}


	_, err = jsc.PublishAsync("COURSES.new", messageBytes)
	if err != nil {
		fmt.Println("Erreur publication NATS:", err)
	}
}

func fetchAgendas(url string)([]models.Agenda, error){
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var agendas []models.Agenda
	err = json.NewDecoder(resp.Body).Decode(&agendas)
	return agendas, err
}

func parseICal(url string, agendaName string)([]models.Course, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	lines := strings.Split(string(body), "\n")

	var courses []models.Course
	var currentCourse models.Course
	currentlyParsing := false
	tmpObj := map[string]interface{}{}

	for _, line := range lines {
		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			currentlyParsing = true
			currentCourse = models.Course{AgendaName: agendaName}
		} else {
			if currentlyParsing {
				if strings.HasPrefix(line, "END:VEVENT") {
					courses = append(courses, currentCourse)
					currentlyParsing = false
				} else {
					if strings.HasPrefix(line, "SUMMARY:") {
                                                // Attention, le dernier caractère est un "carriage return" (\r). On le supprime sinon ça fait échouer toute notre logique.
						tmpObj["summary"] = strings.Replace(strings.Replace(line, "SUMMARY:", "", 1), "\r", "", 1)
						currentCourse.Summary = strings.TrimPrefix(line, "SUMMARY:")
					}
					if strings.HasPrefix(line, "DTSTART:") {
						tmpObj["start"], _ = time.Parse("20060102T150405Z", strings.Replace(strings.Replace(line, "DTSTART:", "", 1), "\r", "", 1))
						val := strings.TrimPrefix(line, "DISTART:")
						currentCourse.Start, _ = time.Parse("20060102T150405Z", val)
					}
				}
			} else {
				continue
			}
		}

	}
	return courses, nil
}