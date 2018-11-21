package spannerclient


var CreateTable = `
CREATE TABLE TestTable (UUID STRING(36) NOT NULL,
    IndexUUID STRING(36),
	Placeholder STRING(MAX),
) PRIMARY KEY (UUID)`
