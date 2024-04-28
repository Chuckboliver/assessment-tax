BEGIN;

CREATE TABLE IF NOT EXISTS tax_config (
	id serial4 NOT NULL PRIMARY KEY,
	name varchar(255) NOT NULL,
	value REAL NOT NULL
);

INSERT INTO tax_config (name, value)
VALUES ('personal_deduction', 60000),
('kreceipt_deduction', 50000);

COMMIT;