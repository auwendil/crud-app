CREATE TABLE IF NOT EXISTS public.books (
                                            id SERIAL PRIMARY KEY,
                                            name varchar(40) NOT NULL,
                                            author varchar(40) NOT NULL
);