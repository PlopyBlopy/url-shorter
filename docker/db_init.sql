-- Create a counter table
CREATE TABLE IF NOT EXISTS counter (
    id INT PRIMARY KEY,
    counter BIGINT NOT NULL DEFAULT 0
);

-- Initialization of the only record in the counter
INSERT INTO counter (id, counter) VALUES (0, 0) ON CONFLICT (id) DO NOTHING;

-- Creating a url table
CREATE TABLE IF NOT EXISTS urls (
    orig_url TEXT PRIMARY KEY,
    short_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- Create a hash index for quick search by orig_url
CREATE  INDEX IF NOT EXISTS idx_orig_url_hash ON urls USING HASH(orig_url);

-- Create a unique hash index for quick short_url search
CREATE INDEX IF NOT EXISTS idx_short_url_hash ON urls USING HASH(short_url);