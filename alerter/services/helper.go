package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
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

	var configs []models.AlertConfig
	if err := json.NewDecoder(resp.Body).Decode(&configs); err != nil {
		return nil, err
	}
	return configs, nil
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
		From:    "chloe.despesse@etu.uca.fr", // IMPORTANT : Adresse uca.fr obligatoire
		To:      []string{to},
		Subject: subject,
		Body:    body,
	}

	jsonBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", models.MailAPIUrl, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ models.MailToken) // Si besoin d'un header Auth

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API Mail erreur: %d", resp.StatusCode)
	}
	return nil
}