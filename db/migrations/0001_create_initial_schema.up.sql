CREATE TYPE sake_category AS ENUM ('BEER', 'SAKE', 'WINE', 'WHISKEY', 'SHOCHU', 'COCKTAIL', 'OTHER');

CREATE TABLE sakes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    category sake_category NOT NULL,
    subcategory TEXT,
    alcohol_percentage NUMERIC(4, 1) NOT NULL CHECK (alcohol_percentage >= 0 AND alcohol_percentage <= 100),
    volume INTEGER NOT NULL CHECK (volume > 0),
    origin TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    stock INTEGER NOT NULL CHECK (stock >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_sakes_category ON sakes(category);
CREATE INDEX idx_sakes_name ON sakes(name);
CREATE INDEX idx_sakes_origin ON sakes(origin);