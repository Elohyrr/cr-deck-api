-- Royal API Personnel - Initial Schema
-- Version: 001
-- Date: 2026-01-10

-- Table: battles
-- Stores individual battle records from top players
CREATE TABLE IF NOT EXISTS battles (
    id SERIAL PRIMARY KEY,
    battle_time TIMESTAMP NOT NULL,
    player_tag VARCHAR(20) NOT NULL,
    opponent_tag VARCHAR(20),
    game_mode VARCHAR(50),
    player_crowns INT,
    opponent_crowns INT,
    deck_signature VARCHAR(255) NOT NULL,
    deck_cards JSONB NOT NULL,
    is_victory BOOLEAN,
    created_at TIMESTAMP DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT unique_player_battle UNIQUE (player_tag, battle_time)
);

-- Indexes for battles table
CREATE INDEX idx_battles_time ON battles(battle_time);
CREATE INDEX idx_battles_deck ON battles(deck_signature);
CREATE INDEX idx_battles_player ON battles(player_tag);

-- Table: meta_decks
-- Aggregated deck statistics
CREATE TABLE IF NOT EXISTS meta_decks (
    deck_signature VARCHAR(255) PRIMARY KEY,
    cards JSONB NOT NULL,
    total_games INT NOT NULL DEFAULT 0,
    wins INT NOT NULL DEFAULT 0,
    losses INT NOT NULL DEFAULT 0,
    win_rate DECIMAL(5,2),
    first_seen TIMESTAMP,
    last_seen TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for meta_decks table
CREATE INDEX idx_meta_win_rate ON meta_decks(win_rate DESC);
CREATE INDEX idx_meta_frequency ON meta_decks(total_games DESC);
CREATE INDEX idx_meta_last_seen ON meta_decks(last_seen DESC);

-- Table: collection_stats
-- Tracks collection runs for monitoring
CREATE TABLE IF NOT EXISTS collection_stats (
    id SERIAL PRIMARY KEY,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    players_processed INT DEFAULT 0,
    battles_collected INT DEFAULT 0,
    battles_stored INT DEFAULT 0,
    errors INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'running',
    error_message TEXT
);

-- Index for collection_stats
CREATE INDEX idx_collection_started ON collection_stats(started_at DESC);

-- Comments for documentation
COMMENT ON TABLE battles IS 'Individual battle records from top 1000 players';
COMMENT ON TABLE meta_decks IS 'Aggregated deck statistics calculated from battles';
COMMENT ON TABLE collection_stats IS 'Monitoring data for collection runs';

COMMENT ON COLUMN battles.deck_signature IS 'Hash of sorted card IDs (e.g., 26000000-26000001-...)';
COMMENT ON COLUMN battles.deck_cards IS 'JSON array of 8 cards with id, name, level';
COMMENT ON COLUMN meta_decks.win_rate IS 'Percentage: (wins / total_games) * 100';
