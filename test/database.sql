CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,  -- الاسم الأول ما خاصوش يكون خاوي
    last_name TEXT NOT NULL,   -- الاسم العائلي ما خاصوش يكون خاوي
    email VARCHAR UNIQUE       -- الإيميل، Unique باش مايتعاودش
);