CREATE TABLE drink_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sake_id UUID NOT NULL REFERENCES sake(id) ON DELETE CASCADE,
    drank_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_drink_logs_user_id ON drink_logs(user_id);
CREATE INDEX idx_drink_logs_sake_id ON drink_logs(sake_id);
CREATE INDEX idx_drink_logs_drank_at ON drink_logs(drank_at);
