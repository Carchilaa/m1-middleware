package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {

	rawDate := "20251228T152000Z"
       // 2006 = année ; 01 = mois ; 02 = jour ; 15 = heure ; 04 = minute ; 05 = seconde
	d, _ := time.Parse("20060102T150405Z", rawDate)

	if d.Before(time.Now()) {
		fmt.Println("Avant !")
	} else {
		fmt.Println("Après !")
	}

	fmt.Println(d)

	//resp, err := http.Get("https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=46030&projectId=3&calType=ical&nbWeeks=8&displayConfigId=128")
	resp, err := http.Get("http://localhost:8080/config/")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	lines := strings.Split(string(body), "\n")

	currentlyParsing := false
	tmpObj := map[string]interface{}{}

	for _, line := range lines {
		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			currentlyParsing = true
		} else {
			if currentlyParsing {
				if strings.HasPrefix(line, "END:VEVENT") {
					fmt.Println(tmpObj)
					tmpObj = map[string]interface{}{}
					currentlyParsing = false
				} else {
					if strings.HasPrefix(line, "SUMMARY:") {
                                                // Attention, le dernier caractère est un "carriage return" (\r). On le supprime sinon ça fait échouer toute notre logique.
						tmpObj["summary"] = strings.Replace(strings.Replace(line, "SUMMARY:", "", 1), "\r", "", 1)
					}
					if strings.HasPrefix(line, "DTSTART:") {
						tmpObj["start"], _ = time.Parse("20060102T150405Z", strings.Replace(strings.Replace(line, "DTSTART:", "", 1), "\r", "", 1))
					}
				}
			} else {
				continue
			}
		}
	}
}
