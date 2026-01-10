import (
    "bytes"
    "embed"
    "html/template"
    "strings"
    "github.com/adrg/frontmatter"
)

//go:embed templates
var embeddedTemplates embed.FS

// ParseTemplate charge le fichier, exécute les variables et extrait le sujet
func ParseTemplate(templatePath string, data interface{}) (string, string, error) {
    // 1. Charger le template
    tmplContent, err := embeddedTemplates.ReadFile("templates/" + templatePath)
    if err != nil {
        return "", "", err
    }

    // 2. Extraire le Frontmatter (Sujet) et le contenu
    var matter FrontMatter
    content, err := frontmatter.Parse(strings.NewReader(string(tmplContent)), &matter)
    if err != nil {
        return "", "", err
    }

    // 3. Exécuter le template (remplacer {{ .EventName }} etc.)
    tmpl, err := template.New("mail").Parse(string(content))
    if err != nil {
        return "", "", err
    }

    var finalBody bytes.Buffer
    if err := tmpl.Execute(&finalBody, data); err != nil {
        return "", "", err
    }

    return finalBody.String(), matter.Subject, nil
}
