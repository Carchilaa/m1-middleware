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
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Impossible de se connecter a NATS:", err)
	}
	defer nc.Close()

	jsc, err = nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

    // Nommage dans Timetable : Stream=SCHEDULER, Sujet=SCHEDULER.>
	_, _ = jsc.AddStream(&nats.StreamConfig{
		Name:     "SCHEDULER",
		Subjects: []string{"SCHEDULER.>"},
	})

	ctx := context.Background()
	sc := scheduler.NewScheduler()

    // On lance le job
	sc.Add(ctx, runJob, time.Minute*1)

    // Log pour dire que ça tourne
    fmt.Println("Scheduler démarré. En attente de jobs.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	sc.Stop()
}

func runJob(ctx context.Context) {
    fmt.Println("Lancement du job de récupération")

	agendas, err := fetchAgendas("http://localhost:8080/agendas/")
	if err != nil {
		fmt.Printf("Erreur lors de la recuperation des agendas: %v\n", err)
		return
	}

	for _, a := range agendas {
		url := fmt.Sprintf("https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=%d&projectId=3&calType=ical&nbWeeks=8&displayConfigId=128", a.UcaId)

		courses, err := parseICal(url, a.Name)
		if err != nil {
			fmt.Printf("Erreur pour l'agenda %s: %v\n", a.Name, err)
			continue
		}

		for _, c := range courses {
			publishToNats(c)
		}
	}
}

func publishToNats(course models.Course) {
	messageBytes, err := json.Marshal(course)
	if err != nil {
		fmt.Println("Erreur JSON:", err)
		return
	}

    // Sujet Timetable: SCHEDULER.events
	_, err = jsc.PublishAsync("SCHEDULER.events", messageBytes)
	if err != nil {
		fmt.Println("Erreur publication NATS:", err)
	} else {
        fmt.Printf("Envoyé: %s (%s)\n", course.Name, course.Uid)
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

func parseICal(url string, agendaName string) ([]models.Course, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	lines := strings.Split(string(body), "\n")

	var courses []models.Course
	var currentCourse models.Course
	currentlyParsing := false

    // format de date standard iCal (UTC)
    layout := "20060102T150405Z"

	for _, line := range lines {
        // on enlève les retours chariot \r et les espaces
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			currentlyParsing = true
			currentCourse = models.Course{}
		} else if strings.HasPrefix(line, "END:VEVENT") {
			if currentlyParsing {
                // on n'ajoute que si on a au moins un nom et un UID
                if currentCourse.Uid != "" && currentCourse.Name != "" {
				    courses = append(courses, currentCourse)
                }
				currentlyParsing = false
			}
		} else if currentlyParsing {
			if strings.HasPrefix(line, "SUMMARY:") {
				currentCourse.Name = strings.TrimPrefix(line, "SUMMARY:")
			}
			if strings.HasPrefix(line, "UID:") {
				currentCourse.Uid = strings.TrimPrefix(line, "UID:")
			}
			if strings.HasPrefix(line, "LOCATION:") {
				currentCourse.Location = strings.TrimPrefix(line, "LOCATION:")
			}
			if strings.HasPrefix(line, "DTSTART:") {
				val := strings.TrimPrefix(line, "DTSTART:")
				t, err := time.Parse(layout, val)
                if err == nil { currentCourse.Start = t }
			}
			if strings.HasPrefix(line, "DTEND:") {
				val := strings.TrimPrefix(line, "DTEND:")
				t, err := time.Parse(layout, val)
                if err == nil { currentCourse.End = t }
			}
		}
	}
	return courses, nil
}
