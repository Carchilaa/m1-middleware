# Projet Middleware : Système d'Alerte d'Emploi du Temps

Equipe : Archila César; Despesse Chloé; Duchosal Jolan

Ce projet met en œuvre une architecture Microservices en Go, orchestrée autour d'un bus d'événements **NATS**.  
L'objectif est de synchroniser des emplois du temps universitaires (format **iCal / UCA**), de détecter les changements via un processus de comparaison (**Diffing**) et d'alerter les étudiants abonnés.

---

## Architecture du Système

Le système est découpé en **4 composants majeurs** interagissant via **APIs REST** et **messages asynchrones NATS**.

---

### 1. Service Config (API REST)

**Rôle :** Point d'entrée pour la configuration des préférences utilisateurs.

**Fonctionnalités :**
- CRUD des Agendas (Association ID Groupe ↔ URL iCal)
- CRUD des Alertes (Abonnements des utilisateurs aux changements d'un agenda)
- **Stockage :** Base de données SQLite dédiée

---

### 2. Service Timetable (API REST + Consumer)

Ce service est le cœur de la gestion des données de cours.

**API REST :**
- CRUD sur les ressources Évènements (cours stockés)

**Consumer (Worker) :**
- Reçoit les flux de données brutes envoyés par le Scheduler
- Compare avec la base SQLite locale
- Détecte les modifications (Ajout, Suppression, Modification)
- Publie un événement d’alerte si différence détectée

---

### 3. Service Scheduler

**Rôle :** Récupérateur de données (Fetcher)

**Fonctionnement :**
- Exécutions périodiques
- Interroge les serveurs UCA (.ics)
- Parse les données
- Envoie les données via NATS au Timetable

---

### 4. Service Alerter (Consumer)

**Rôle :** Gestionnaire de notifications

**Fonctionnement :**
- Écoute les événements de modification envoyés par Timetable
- Interroge l’API Config pour connaître les abonnés concernés
- Envoie les emails via une API externe

---

## Prérequis

- Go (v1.20+)
- NATS Server (JetStream activé)
- Accès internet (API Mail + UCA)

---

## Installation et Lancement

Assurez-vous que votre serveur NATS est lancé localement :

```bash
nats-server -js
```

### Exécution des Microservices

Chaque microservice peut être lancé depuis son répertoire racine.  
Ouvrez un terminal distinct pour chaque service :

---

### 1. Lancer l'API Config

```bash
cd API_Middleware/api-config
go run cmd/main.go
```

---

### 2. Lancer le Service Timetable (API + Consumer)

```bash
cd API_Middleware/timetable
go run cmd/main.go
```

---

### 3. Lancer l'Alerter

```bash
cd API_Middleware/alerter
go run cmd/main.go
```

---

### 4. Lancer le Scheduler (en dernier)

```bash
cd API_Middleware/scheduler
go run cmd/main.go
```

---

## Pistes d'Amélioration (Post-Mortem)

Bien que le système soit fonctionnel, plusieurs optimisations seraient pertinentes pour une V2.

---

### 1. Gestion des Fuseaux Horaires

Actuellement le parsing et stockage sont en UTC.  
Ajout recommandé :
- Conversion automatique vers Europe/Paris
- Évite les décalages dans les notifications

---

### 2. Découpage du Scheduler (Single Responsibility)

Le scheduler concentre trop de logique.  
Refactorisation recommandée en modules séparés :
- Cron Job
- Fetcher (HTTP client UCA)
- Parser (iCal → structures Go)
- Publisher (NATS)

Améliore testabilité et lisibilité.

---

### 3. Gestion de la Configuration

Externaliser les secrets :
- Tokens API Mail
- URLs
- Credentials DB

→ via variables d’environnement (`.env`) au lieu de constantes dans le code.

---
