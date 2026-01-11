package collector

// TopPlayerTags contient la liste des tags de top players à tracker.
// Cette liste doit être mise à jour manuellement environ 1x/mois.
// Source: RoyaleAPI leaderboard ou équivalent.
// Dernière update: 2026-01-11
var TopPlayerTags = []string{
	"#22LV0QUQJ", // Mini Léo CoC (validé: 10668 trophées, 2026-01-11)
	"#2PP",       // Top 1 historique
	"#8L9L9GL",   // Mohamed Light (Top EU)
	"#2CCCP",     // Nova Esports clan tag (erreur: c'est un clan, pas un joueur)
	"#YC8UY",     // Chief Pat (créateur de contenu)
	"#8QVJ8PL",   // Surgical Goblin (pro player)
	"#2LGRCU",    // Boss CR
	"#9CQ2U8QJ",  // Jonas (Top EU)
	"#YV2GJC",    // CWA (OJ - créateur de contenu)
	"#8PPRR",     // Morten (Top DK)
}

// GetTopPlayerTags retourne la liste complète des top players à tracker.
// Cette fonction permet d'encapsuler l'accès à la liste et facilite
// les tests unitaires et l'évolution future (chargement depuis fichier, etc.).
func GetTopPlayerTags() []string {
	// Retourner une copie pour éviter les modifications externes
	tags := make([]string, len(TopPlayerTags))
	copy(tags, TopPlayerTags)
	return tags
}

// GetTopPlayerCount retourne le nombre de players trackés.
// Utile pour les logs et métriques.
func GetTopPlayerCount() int {
	return len(TopPlayerTags)
}
