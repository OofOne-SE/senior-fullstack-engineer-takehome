CREATE TABLE IF NOT EXISTS weather (
                                       id SERIAL PRIMARY KEY,
                                       timestamp TIMESTAMPTZ NOT NULL,
                                       temperature DOUBLE PRECISION,
                                       humidity DOUBLE PRECISION
);

SELECT create_hypertable('weather', 'timestamp', if_not_exists => TRUE);
