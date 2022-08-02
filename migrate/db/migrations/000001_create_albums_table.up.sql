CREATE TABLE IF NOT EXISTS albums(
   id VARCHAR (50) PRIMARY KEY,
   title VARCHAR (50) NOT NULL,
   artist VARCHAR (50) NOT NULL,
   price double precision NOT NULL
);