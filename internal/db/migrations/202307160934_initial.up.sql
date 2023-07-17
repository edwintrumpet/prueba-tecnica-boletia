BEGIN;

CREATE TABLE requests (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    requested_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    last_updated_at TIMESTAMP WITHOUT TIME ZONE,
    request_duration DECIMAL NOT NULL,
    response_status VARCHAR(3),
    is_ok BOOLEAN NOT NULL,
    error_msg TEXT
);

CREATE TABLE currencies (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(5) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    request_id uuid NOT NULL,

    CONSTRAINT fk_request FOREIGN KEY(request_id) REFERENCES requests(id)
);

END;
