package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
	"net/http"
	"strings"
	"time"
	"io"

	"github.com/adrg/frontmatter"
	"github.com/sirupsen/logrus"
	"alerter/models"
)

func getSubscribers(agendaID string) ([]models.AlertConfig, error) {
    // Appel HTTP GET vers ton API Config
	resp, err := http.Get(models.ConfigAPIUrl + agendaID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var config models.AlertConfig
    if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
        return nil, err
    }
    return []models.AlertConfig{config}, nil
}

func parseTemplate(path string, data interface{}) (string, string, error) {
    // Utilise "templates/" + path car embed est à la racine du package souvent
	fileContent, err := models.EmbeddedTemplates.ReadFile("templates/" + path)
	if err != nil {
		return "", "", err
	}

	var matter models.FrontMatter
	content, err := frontmatter.Parse(strings.NewReader(string(fileContent)), &matter)
	if err != nil {
		return "", "", err
	}

	tmpl, err := template.New("mail").Parse(string(content))
	if err != nil {
		return "", "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", "", err
	}

	return buf.String(), matter.Subject, nil
}

func sendMail(to, subject, body string) error {
    payload := models.MailRequest{
        Recipient: to,
        Subject:   subject,
        Content:   body,
    }

    jsonBytes, _ := json.Marshal(payload)

    req, _ := http.NewRequest("POST", models.MailAPIUrl, bytes.NewBuffer(jsonBytes))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", models.MailToken)

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    responseBody, _:= io.ReadAll(resp.Body)
    
    // Log important pour voir si l'API accepte ou rejette
    logrus.Infof("Réponse API Mail pour %s : Code=%d Body=%s", to, resp.StatusCode, string(responseBody))

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("API Mail erreur: %d - %s", resp.StatusCode, string(responseBody))
    }
    return nil
}