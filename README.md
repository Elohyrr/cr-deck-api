# 🏆 Royal API Personnel

Backend service pour collecter et analyser les données de méta-jeu Clash Royale des top 1000 joueurs mondiaux.

## 📋 Description

Royal API Personnel collecte quotidiennement les battlelogs des top 1000 joueurs via l'API officielle Supercell, analyse les decks utilisés, calcule les statistiques de méta (winrate, fréquence), et expose ces données via une API REST sécurisée.

**Stack**: Go 1.21+ • PostgreSQL 16 • Docker Compose

## 🚀 Installation

### Prérequis

- Go 1.21+
- PostgreSQL 16+
- Docker & Docker Compose (pour déploiement)
- Clé API Supercell ([créer ici](https://developer.clashroyale.com))

### Configuration

1. Cloner le projet:
```bash
git clone https://github.com/leopoldhub/royal-api-personal.git
cd royal-api-personal
```

2. Créer le fichier `.env.local`:
```bash
cp .env.example .env.local
# Éditer .env.local avec vos vraies valeurs
```

3. Obtenir une clé API Supercell:
   - Créer un compte sur https://developer.clashroyale.com
   - Créer une clé API
   - Whitelister votre IP (ou celle de votre VPS)
   - Copier la clé dans `SUPERCELL_API_KEY`

4. Générer un token Bearer sécurisé:
```bash
openssl rand -base64 32
# Copier dans API_TOKEN
```

## 🏃 Utilisation

### Avec Docker Compose (recommandé)

```bash
docker-compose up -d
```

Services démarrés:
- `postgres`: PostgreSQL 16
- `api`: API REST (port 8080)
- `collector`: Collecteur automatique (cron 24h)

### En local (développement)

**1. Démarrer PostgreSQL**:
```bash
docker run -d \
  --name postgres-royale \
  -e POSTGRES_DB=royale_api \
  -e POSTGRES_USER=royale \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  postgres:16-alpine
```

**2. Build le projet**:
```bash
go build -o royal-api cmd/royal-api/main.go
```

**3. Lancer la collecte manuelle**:
```bash
./royal-api collect
```

**4. Lancer l'API REST**:
```bash
./royal-api serve
```

## 📡 API Endpoints

**Base URL**: `http://localhost:8080`

**Authentication**: Tous les endpoints (sauf `/health`) nécessitent un header:
```
Authorization: Bearer <API_TOKEN>
```

### GET `/health`

Health check du service.

**Response** (200 OK):
```json
{
  "status": "healthy",
  "database": "connected",
  "last_collection": "2026-01-10T22:00:00Z",
  "total_battles": 15234,
  "total_decks": 387
}
```

### GET `/decks/meta`

Liste des decks méta triés par winrate ou fréquence.

**Query Parameters**:
- `limit` (default: 50): Nombre de decks à retourner
- `sort` (default: win_rate): `win_rate` ou `frequency`
- `min_games` (default: 10): Minimum de parties jouées

**Example**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/decks/meta?limit=20&sort=win_rate&min_games=20"
```

**Response** (200 OK):
```json
{
  "decks": [
    {
      "signature": "26000000-26000001-26000003-26000010-26000014-26000027-26000042-28000000",
      "cards": [
        {"id": "26000000", "name": "Knight"},
        {"id": "26000001", "name": "Archers"},
        ...
      ],
      "stats": {
        "total_games": 143,
        "wins": 89,
        "losses": 54,
        "win_rate": 62.24,
        "last_seen": "2026-01-10T21:45:00Z"
      }
    }
  ],
  "metadata": {
    "total_decks": 20,
    "last_updated": "2026-01-10T22:00:00Z"
  }
}
```

### GET `/decks/{signature}`

Détails d'un deck spécifique avec exemples de battles récents.

**Example**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/decks/26000000-26000001-..."
```

**Response** (200 OK):
```json
{
  "deck": {
    "signature": "26000000-26000001-...",
    "cards": [...],
    "stats": {
      "total_games": 143,
      "win_rate": 62.24
    }
  },
  "recent_battles": [
    {
      "battle_time": "2026-01-10T21:30:00Z",
      "player_tag": "#2PP",
      "opponent_tag": "#ABC",
      "is_victory": true
    }
  ]
}
```

### GET `/stats/summary`

Statistiques globales de la collection.

**Response** (200 OK):
```json
{
  "collection": {
    "total_battles": 15234,
    "total_decks": 387,
    "last_collection": "2026-01-10T22:00:00Z",
    "players_tracked": 1000
  },
  "top_deck": {
    "signature": "...",
    "total_games": 234
  },
  "best_deck": {
    "signature": "...",
    "win_rate": 71.43,
    "min_games": 20
  },
  "top_cards": [
    {"id": "26000000", "name": "Knight", "usage_rate": 45.2},
    {"id": "26000042", "name": "Hog Rider", "usage_rate": 38.7}
  ]
}
```

## 🐳 Docker Compose

**Fichier `docker-compose.yml`** inclus avec 3 services:
- `postgres`: Base de données PostgreSQL
- `api`: Serveur API REST
- `collector`: Job de collecte quotidienne

```bash
docker-compose up -d       # Démarrer
docker-compose logs -f     # Logs en temps réel
docker-compose down        # Arrêter
```

## 🛠️ Développement

### Structure du projet

```
royal-api-personal/
├── cmd/royal-api/          # Entry point CLI
├── internal/
│   ├── api/                # HTTP server + handlers
│   ├── collector/          # Business logic collecte
│   ├── config/             # Configuration
│   ├── database/           # DB operations
│   ├── errors/             # Custom error types
│   └── models/             # Data structures
├── pkg/
│   ├── supercell/          # Supercell API client
│   └── utils/              # Utilities
├── migrations/             # SQL migrations
└── .context/               # Documentation projet
```

### Tests

```bash
go test ./...              # Tous les tests
go test -v ./pkg/supercell # Tests d'un package
go test -cover ./...       # Coverage
```

### Commandes disponibles

```bash
# Collecte manuelle
./royal-api collect

# Collecte en boucle (24h)
./royal-api collect-loop

# Démarrer l'API REST
./royal-api serve
```

## 📊 Fonctionnement

1. **Collecte quotidienne** (automated):
   - Fetch top 1000 joueurs via `/locations/global/rankings/players`
   - Pour chaque joueur: fetch 25 derniers combats via `/players/{tag}/battlelog`
   - Filtre: combats PvP Ladder uniquement
   - Insertion en DB avec batch insert
   - Recalcul des statistiques méta

2. **Purge automatique**:
   - Suppression des combats > 7 jours
   - Exécuté après chaque collecte

3. **API REST**:
   - Authentification Bearer token
   - Requêtes SQL optimisées avec indexes
   - Responses JSON

## ⚠️ Limitations MVP

- Pas d'authentification multi-utilisateurs
- Pas de frontend UI
- Pas de notifications temps réel
- Pas de monitoring/alerting avancé

Ces features sont prévues pour les phases 2 et 3.

## 🔒 Sécurité

- Token Bearer pour authentification API
- Token JWT Supercell avec IP whitelisting
- Variables d'environnement pour secrets
- PostgreSQL avec credentials sécurisés

## 📝 License

MIT

## 🤝 Contribution

Contributions bienvenues ! Créer une issue ou une PR.

## 📧 Contact

[@leopoldhub](https://github.com/leopoldhub)

---

**Made with ❤️ for Clash Royale meta analysis**
