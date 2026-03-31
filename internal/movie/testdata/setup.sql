-- Setup schema
CREATE TABLE movies (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    year integer NOT NULL,
    runtime integer NOT NULL,
    genres text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);

ALTER TABLE movies ADD CONSTRAINT movies_runtime_check CHECK (runtime >= 0);
ALTER TABLE movies ADD CONSTRAINT movies_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
ALTER TABLE movies ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);

-- Seed data
INSERT INTO movies (id, title, year, runtime, genres, version) 
VALUES (
    '0ee58bdc-e2de-454f-ad22-e02caa53cc31', 
    'Moana', 
    2016, 
    107, 
    '{animation, adventure}', 
    1
);

INSERT INTO movies (id, title, year, runtime, genres, version) 
VALUES (
    '6350d919-27a6-4b7b-ae8d-4f744bdf0282', 
    'Black Panther', 
    2018, 
    134, 
    '{action, adventure}', 
    1
);

INSERT INTO movies (id, title, year, runtime, genres, version) 
VALUES (
    '72dda126-efbc-41f0-8c1d-14141484f540', 
    'Deadpool', 
    2016, 
    108, 
    '{action, comedy}', 
    1
);

INSERT INTO movies (id, title, year, runtime, genres, version) 
VALUES (
    'c9396c5f-86d8-459a-b20a-5359c1bd15e7', 
    'The Breakfast Club', 
    1986, 
    96, 
    '{drama}', 
    1
);
