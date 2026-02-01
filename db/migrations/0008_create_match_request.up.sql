CREATE TYPE match_status AS ENUM ('PENDING', 'ACCEPTED', 'REJECTED');

CREATE TABLE match_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status match_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_match_requests_from_user_id ON match_requests(from_user_id);
CREATE INDEX idx_match_requests_to_user_id ON match_requests(to_user_id);
CREATE INDEX idx_match_requests_status ON match_requests(status);
