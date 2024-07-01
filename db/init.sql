-- Table to store user information
CREATE TABLE IF NOT EXISTS users
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    email      TEXT UNIQUE NOT NULL,
    username   TEXT UNIQUE NOT NULL,
    password   TEXT        NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table to store posts
CREATE TABLE IF NOT EXISTS posts
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL,
    title      TEXT    NOT NULL,
    content    TEXT    NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Table to store comments
CREATE TABLE IF NOT EXISTS comments
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id    INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    content    TEXT    NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Table to store categories
CREATE TABLE IF NOT EXISTS categories
(
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

-- Table to associate posts with categories (many-to-many relationship)
CREATE TABLE IF NOT EXISTS post_categories
(
    post_id     INTEGER NOT NULL,
    category_id INTEGER NOT NULL,

    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (category_id) REFERENCES categories (id),
    PRIMARY KEY (post_id, category_id)
);

-- Table to store likes and dislikes
CREATE TABLE IF NOT EXISTS likes
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL,
    post_id    INTEGER,
    comment_id INTEGER,
    is_like    BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (comment_id) REFERENCES comments (id),

    CHECK (post_id IS NOT NULL OR comment_id IS NOT NULL)
);

-- Table to store user sessions
CREATE TABLE IF NOT EXISTS sessions
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER   NOT NULL,
    session_token TEXT      NOT NULL,
    expires_at    TIMESTAMP NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Insert some example categories
INSERT OR IGNORE INTO categories (name)
VALUES ('General'),
       ('Technology'),
       ('Health'),
       ('Entertainment');
