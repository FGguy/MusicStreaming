CREATE TABLE IF NOT EXISTS Users (
    username VARCHAR(30),
    password VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,

    --Roles
    scrobblingEnabled BOOLEAN NOT NULL DEFAULT FALSE,
    ldapAuthenticated BOOLEAN NOT NULL DEFAULT FALSE,
    adminRole BOOLEAN NOT NULL DEFAULT FALSE,
    settingsRole BOOLEAN NOT NULL DEFAULT TRUE,
    streamRole BOOLEAN NOT NULL DEFAULT TRUE,
    jukeboxRole BOOLEAN NOT NULL DEFAULT FALSE,
    downloadRole BOOLEAN NOT NULL DEFAULT FALSE,
    uploadRole BOOLEAN NOT NULL DEFAULT FALSE,
    playlistRole BOOLEAN NOT NULL DEFAULT FALSE,
    coverArtRole BOOLEAN NOT NULL DEFAULT FALSE,
    commentRole BOOLEAN NOT NULL DEFAULT FALSE,
    podcastRole BOOLEAN NOT NULL DEFAULT FALSE,
    shareRole BOOLEAN NOT NULL DEFAULT FALSE,
    videoConversionRole BOOLEAN NOT NULL DEFAULT FALSE,
    musicFolderId TEXT,
    maxBitRate INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY(username)
);

CREATE INDEX IF NOT EXISTS idx_users_username
ON Users(username);

CREATE TABLE IF NOT EXISTS Covers(
    cover_id TEXT,
    path TEXT NOT NULL,
    PRIMARY KEY(cover_id)
);

CREATE TABLE IF NOT EXISTS Artists (
    artist_id SERIAL,
    name TEXT NOT NULL,
    cover_art TEXT,
    album_count INTEGER,
    PRIMARY KEY(artist_id)
);

CREATE TABLE IF NOT EXISTS Albums (
    album_id SERIAL,
    artist_id INTEGER,
    name TEXT NOT NULL,
    cover_art TEXT,
    song_count INTEGER,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    duration INTEGER,
    artist TEXT,
    PRIMARY KEY(album_id),
    FOREIGN KEY (artist_id) REFERENCES Artists(artist_id)
);

CREATE TABLE IF NOT EXISTS Songs (
    song_id SERIAL,
    album_id INTEGER,
    title TEXT NOT NULL,
    album TEXT,
    artist TEXT,
    is_dir BOOLEAN,
    cover_art TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    duration INTEGER,
    bit_rate INTEGER,
    size INTEGER,
    suffix TEXT,
    content_type TEXT,
    is_video BOOLEAN,
    path TEXT NOT NULL,
    PRIMARY KEY(song_id),
    FOREIGN KEY (album_id) REFERENCES Albums(album_id)
);