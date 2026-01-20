DROP MATERIALIZED VIEW IF EXISTS positions_mv;
DROP TRIGGER IF EXISTS trigger_record_price_change ON players;
DROP FUNCTION IF EXISTS record_player_price_change;


DROP TABLE IF EXISTS player_nba_mapping;
DROP TABLE IF EXISTS nba_weekly_stats;
DROP TABLE IF EXISTS nba_season_stats;
DROP TABLE IF EXISTS nba_career_stats;
DROP TABLE IF EXISTS player_price_history;
DROP TABLE IF EXISTS holdings;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS nba_players;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS users;