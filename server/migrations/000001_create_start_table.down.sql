-- Откат начальных таблиц приложения

DROP TABLE IF EXISTS bytes; 
DROP TABLE IF EXISTS texts; 
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users; 

