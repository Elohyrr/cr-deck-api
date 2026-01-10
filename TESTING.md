# 🧪 Guide de Tests Manuels - Royal API Personnel MVP

**Version**: 1.0.0-mvp  
**Date**: 2026-01-10

---

## 📋 Prérequis

- Docker & Docker Compose installés
- Clé API Supercell valide ([obtenir ici](https://developer.clashroyale.com))
- Fichier `.env.local` configuré

---

## 🚀 Test 1: Démarrage des services

**Objectif**: Vérifier que tous les services démarrent correctement.

**Commandes**:
```bash
cd royal-api-personal
cp .env.example .env.local
# Éditer .env.local avec vos vraies valeurs
docker-compose up -d
```

**Vérifications**:
```bash
docker-compose ps
```

**Résultat attendu**:
```
NAME                    STATUS
royal-api-postgres      Up (healthy)
royal-api-server        Up
royal-api-collector     Up
```

**Logs**:
```bash
docker-compose logs -f api
docker-compose logs -f collector
```

---

## 🩺 Test 2: Health Check

**Objectif**: Vérifier que l'API répond et que la DB est connectée.

**Commande**:
```bash
curl http://localhost:8080/health
```

**Résultat attendu**:
```json
{
  "status": "healthy",
  "database": "connected",
  "total_battles": 0,
  "total_decks": 0
}
```

**Critères de succès**:
- ✅ Status code 200
- ✅ `database: "connected"`
- ✅ Pas d'erreur dans les logs

---

## 🔒 Test 3: Authentification

**Objectif**: Vérifier que le Bearer token est requis.

**Test 3.1: Sans token (doit échouer)**:
```bash
curl http://localhost:8080/decks/meta
```

**Résultat attendu**:
- Status code: **401 Unauthorized**
- Body: `{"error": "unauthorized", "message": "missing Bearer token"}`

**Test 3.2: Token invalide (doit échouer)**:
```bash
curl -H "Authorization: Bearer wrong_token" \
  http://localhost:8080/decks/meta
```

**Résultat attendu**:
- Status code: **401 Unauthorized**
- Body: `{"error": "unauthorized", "message": "invalid token"}`

**Test 3.3: Token valide (doit réussir)**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/decks/meta
```

**Résultat attendu**:
- Status code: **200 OK**
- Body JSON avec clés `decks` et `metadata`

---

## 📊 Test 4: Collecte manuelle

**Objectif**: Déclencher une collecte et vérifier l'insertion en DB.

**Commande**:
```bash
docker-compose exec api ./royal-api -command collect
```

**Résultat attendu** (dans les logs):
```
Starting collection for top 1000 players
Fetched 1000 top players
Progress: 100/1000 players processed
Progress: 200/1000 players processed
...
Collected X raw battles from 1000 players
Filtered to Y PvP Ladder battles
Parsed Z valid battles
Stored Z battles in database
Recalculated meta deck statistics
Purged N old battles (7+ days)
Collection completed in 2m34s
```

**Vérification DB**:
```bash
docker-compose exec postgres psql -U royale -d royale_api -c "SELECT COUNT(*) FROM battles"
docker-compose exec postgres psql -U royale -d royale_api -c "SELECT COUNT(*) FROM meta_decks"
```

**Critères de succès**:
- ✅ `battles` > 0
- ✅ `meta_decks` > 0
- ✅ Pas d'erreurs fatales (404 acceptables)
- ✅ Durée < 10 minutes

---

## 🎮 Test 5: Endpoint GET /decks/meta

**Objectif**: Récupérer les decks méta avec différents filtres.

**Test 5.1: Par défaut (top 50 par winrate)**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/decks/meta
```

**Résultat attendu**:
```json
{
  "decks": [
    {
      "signature": "26000000-26000001-...",
      "cards": [
        {"id": "26000000", "name": "Knight", "level": 14},
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

**Test 5.2: Par fréquence**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/decks/meta?sort=frequency&limit=10"
```

**Critères de succès**:
- ✅ Deck[0].stats.total_games ≥ Deck[1].stats.total_games

**Test 5.3: Avec min_games**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/decks/meta?min_games=50"
```

**Critères de succès**:
- ✅ Tous les decks ont `total_games ≥ 50`

---

## 🔍 Test 6: Endpoint GET /decks/{signature}

**Objectif**: Récupérer les détails d'un deck spécifique.

**Prérequis**: Copier une `signature` depuis le Test 5.

**Commande**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/decks/26000000-26000001-26000003-..."
```

**Résultat attendu**:
```json
{
  "deck": {
    "signature": "26000000-26000001-...",
    "cards": [...],
    "stats": {...}
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

**Test 6.2: Signature invalide (404)**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/decks/invalid-signature"
```

**Résultat attendu**:
- Status code: **404 Not Found**
- Body: `{"error": "deck not found"}`

---

## 📈 Test 7: Endpoint GET /stats/summary

**Objectif**: Récupérer les statistiques globales.

**Commande**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/stats/summary
```

**Résultat attendu**:
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
    "win_rate": 71.43
  },
  "top_cards": []
}
```

**Critères de succès**:
- ✅ `total_battles` > 0
- ✅ `total_decks` > 0
- ✅ `top_deck.total_games` est le maximum
- ✅ `best_deck.win_rate` est élevé (> 50%)

---

## 🔄 Test 8: Collecte automatique (24h loop)

**Objectif**: Vérifier que le collector tourne en boucle.

**Vérification**:
```bash
docker-compose logs -f collector
```

**Résultat attendu** (dans les logs):
```
Starting collection loop (every 24 hours)...
Starting collection run...
Collection completed: 12543 battles stored in 2m34s
Sleeping for 24 hours...
```

**Critères de succès**:
- ✅ Message "Sleeping for 24 hours..." présent
- ✅ Container `royal-api-collector` status = `Up`
- ✅ Pas de restart en boucle

---

## 💾 Test 9: Persistance des données

**Objectif**: Vérifier que les données survivent à un restart.

**Test 9.1: Avant restart**:
```bash
docker-compose exec postgres psql -U royale -d royale_api \
  -c "SELECT COUNT(*) FROM battles"
# Noter le nombre
```

**Test 9.2: Restart**:
```bash
docker-compose restart
docker-compose ps
```

**Test 9.3: Après restart**:
```bash
docker-compose exec postgres psql -U royale -d royale_api \
  -c "SELECT COUNT(*) FROM battles"
# Le nombre doit être identique
```

**Critères de succès**:
- ✅ Même nombre de battles avant/après
- ✅ Volume `postgres_data` existe
- ✅ API répond immédiatement

---

## 🧹 Test 10: Purge automatique (7 jours)

**Objectif**: Vérifier que les vieilles données sont supprimées.

**Simulation** (pour test rapide):
```bash
docker-compose exec postgres psql -U royale -d royale_api -c "
  UPDATE battles 
  SET battle_time = NOW() - INTERVAL '8 days' 
  WHERE id IN (SELECT id FROM battles LIMIT 100)
"
```

**Déclencher purge**:
```bash
docker-compose exec api ./royal-api -command collect
```

**Vérification**:
```bash
docker-compose exec postgres psql -U royale -d royale_api -c "
  SELECT COUNT(*) FROM battles 
  WHERE battle_time < NOW() - INTERVAL '7 days'
"
# Doit retourner 0
```

**Critères de succès**:
- ✅ Battles > 7 jours = 0
- ✅ Message "Purged X old battles" dans les logs

---

## 🚨 Test 11: Gestion des erreurs

**Test 11.1: API Supercell down (simulation)**:
- Mettre une clé API invalide dans `.env.local`
- Redémarrer: `docker-compose restart collector`
- Logs doivent montrer: "collection failed: ..."
- ✅ Pas de crash du container

**Test 11.2: PostgreSQL down**:
```bash
docker-compose stop postgres
curl http://localhost:8080/health
```
- Status: `database: "disconnected"`
- ✅ API reste up

**Test 11.3: Rate limit 429** (si atteint):
- Logs doivent montrer: "rate limit exceeded"
- ✅ Retry automatique avec backoff

---

## ✅ Checklist finale MVP

**Infrastructure**:
- [ ] 3 containers démarrent (postgres, api, collector)
- [ ] Healthchecks OK
- [ ] Logs propres sans erreurs critiques

**API REST**:
- [ ] GET /health répond 200
- [ ] Auth Bearer token fonctionne (401 sans token)
- [ ] GET /decks/meta retourne des decks
- [ ] GET /decks/{signature} retourne détails
- [ ] GET /stats/summary retourne stats

**Collecte**:
- [ ] Collecte manuelle réussit (< 10 min)
- [ ] Données insérées en DB (battles + meta_decks)
- [ ] Collecte automatique tourne en boucle 24h
- [ ] Purge supprime battles > 7 jours

**Persistance**:
- [ ] Données survivent au restart containers
- [ ] Volume postgres_data créé

**Documentation**:
- [ ] README.md complet
- [ ] TESTING.md (ce fichier) validé
- [ ] .env.example fourni

---

## 🐛 Troubleshooting

**Problème**: "connection refused" vers postgres
- **Solution**: Attendre que postgres soit `healthy`
- **Commande**: `docker-compose logs postgres`

**Problème**: 401 Unauthorized
- **Solution**: Vérifier `API_TOKEN` dans `.env.local`
- **Commande**: `echo $API_TOKEN`

**Problème**: Collecte échoue avec 404
- **Solution**: Normal pour certains joueurs (bannis/inactifs)
- **Action**: Vérifier que > 90% des joueurs réussissent

**Problème**: Rate limit 429
- **Solution**: Réduire `TOP_PLAYERS_LIMIT` à 500
- **Action**: Attendre 1h et réessayer

**Problème**: Collecte trop lente (> 10 min)
- **Solution**: Problème réseau ou API Supercell
- **Action**: Vérifier logs pour timeout

---

## 📊 Métriques de succès MVP

**Performance**:
- ✅ Collecte top 1000 < 5 minutes
- ✅ API response time < 200ms
- ✅ Pas de memory leak après 24h

**Données**:
- ✅ > 10,000 battles après 1ère collecte
- ✅ > 200 decks uniques
- ✅ Win rates cohérents (40-70%)

**Stabilité**:
- ✅ Pas de crash pendant 7 jours
- ✅ Collecte automatique fonctionne sans intervention
- ✅ Pas de blocage rate limit (< 1001 req/jour)

---

**Validation complète** : Tous les tests passent ✅  
**MVP prêt pour production** : OUI / NON

---

**Validé par**: _________  
**Date**: __/__/____
