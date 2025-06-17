INSERT INTO projects (name)
SELECT 'first'
WHERE (SELECT COUNT(*) FROM projects) = 0;