IF db_id('dbname') IS NULL
    CREATE DATABASE ideone

GO

CREATE TABLE ideone.ie (
  id BIGINT AUTO_INCREMENT NOT NULL primary key,
	fullurl VARCHAR(60),
	codedate VARCHAR(60),
	codekey VARCHAR(30),
	size BIGINT,
	lines BIGINT,
	language VARCHAR(30),
	status VARCHAR(30),
  txt LONGTEXT,
  UNIQUE(fullurl, codedate, codekey, size, lines, language, status)
);