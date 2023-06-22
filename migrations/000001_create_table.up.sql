CREATE TABLE
    IF NOT EXISTS messages (
        id SERIAL PRIMARY KEY,
        client_id VARCHAR(255) NOT NULL,
        title  VARCHAR(255) NOT NULL,
        payload  JSONB NOT NULL,
        cta_url  VARCHAR(255) NULL,
        "read" BOOLEAN NOT NULL DEFAULT(false),
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now (),
        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now (),
    );