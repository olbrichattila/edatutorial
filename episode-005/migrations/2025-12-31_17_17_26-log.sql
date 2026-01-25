CREATE TABLE logs (
    action_id CHAR(36) NOT NULL,
    level CHAR(16) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    message TEXT NOT NULL,
    INDEX idx_action_id (action_id)
);