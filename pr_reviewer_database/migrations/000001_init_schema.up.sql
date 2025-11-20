CREATE TABLE teams (
    name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    team_name VARCHAR(255) NOT NULL,

    CONSTRAINT fk_users_team FOREIGN KEY (team_name)
        REFERENCES teams(name) ON DELETE RESTRICT
);

CREATE INDEX idx_users_team_name ON users(team_name);

CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE pull_requests (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) NOT NULL, 
    status pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    merged_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_pr_author FOREIGN KEY (author_id)
        REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX idx_pull_requests_author_id ON pull_requests(author_id);

CREATE TABLE pull_requests_reviewers (
    pull_request_id VARCHAR(255) NOT NULL,
    reviewer_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (pull_request_id, reviewer_id),

    CONSTRAINT fk_pr_reviewers_pr FOREIGN KEY (pull_request_id)
        REFERENCES pull_requests(id) ON DELETE CASCADE,
    CONSTRAINT fk_pr_reviewers_user FOREIGN KEY (reviewer_id)
        REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX idx_pr_reviewers_reviewer_id ON pull_requests_reviewers(reviewer_id);