
CREATE TABLE
    IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(50) UNIQUE,
        icon TEXT UNIQUE
    );

CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        firstname TEXT NOT NULL ,
        lastname TEXT NOT NULL ,
        email TEXT NOT NULL UNIQUE,
        gender TEXT NOT NULL ,
        age INTEGER NOT NULL ,
        password TEXT NOT NULL,
        session TEXT DEFAULT NULL
    );

CREATE TABLE
    IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(30),
        title VARCHAR(255),
        description TEXT,
        time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        userID INTEGER  NOT NULL,
        image_path TEXT,
        FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        postID INTEGER,
        username VARCHAR(30),
        comment TEXT,
        time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (postID) REFERENCES posts (id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS likes (
        userID INTEGER,
        postID INTEGER,
        value VARCHAR(2),
        PRIMARY KEY (userID, postID),
        FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE,
        FOREIGN KEY (postID) REFERENCES posts (id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS categories_post (
        categoryID INTEGER,
        postID INTEGER,
        PRIMARY KEY (categoryID, postID),
        FOREIGN KEY (categoryID) REFERENCES categories (id) ON DELETE CASCADE,
        FOREIGN KEY (postID) REFERENCES posts (id) ON DELETE CASCADE
    );



CREATE TRIGGER IF NOT EXISTS post_cleanup_trigger
AFTER DELETE ON posts
BEGIN
    DELETE FROM commentsLikes WHERE commentID IN (SELECT id FROM comments WHERE postID = OLD.id);

    DELETE FROM comments WHERE postID = OLD.id;

    DELETE FROM likes WHERE postID = OLD.id;

    DELETE FROM categories_post WHERE postID = OLD.id;
END;


CREATE TRIGGER IF NOT EXISTS user_cleanup_trigger
AFTER DELETE ON users
BEGIN
    DELETE FROM posts WHERE userID = OLD.id;

END;