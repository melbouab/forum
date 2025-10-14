-- insert a random data to the data base for test 
INSERT OR IGNORE INTO categories (name) VALUES ('Sport');
INSERT OR IGNORE INTO categories (name) VALUES ('Politic');
INSERT OR IGNORE INTO categories (name) VALUES ('Economic');
INSERT OR IGNORE INTO categories (name) VALUES ('Music');
INSERT OR IGNORE INTO categories (name) VALUES ('Education');

INSERT OR IGNORE INTO users (user_name, email, password)  VALUES ('ilyass', 'ilyass@gmail.com', 'ilyasshashed');
INSERT OR IGNORE INTO users (user_name, email, password)  VALUES ('mohamed', 'mohamed@gmail.com', 'mohamedhashed');

INSERT OR IGNORE INTO posts (user_id, title, content, category_id, image_path) VALUES (1, 'post 1', 'Hello from ilyass!', 1, '');
INSERT OR IGNORE INTO posts (user_id, title, content, category_id, image_path)  VALUES (2, 'post 2', 'Hi from mohamed', 2, '');