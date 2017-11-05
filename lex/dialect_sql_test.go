package lex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlDialectInit(t *testing.T) {
	// Make sure we can init more than once, see if it panics
	SqlDialect.Init()
	for _, stmt := range SqlDialect.Statements {
		assert.NotEqual(t, "", stmt.String())
	}
}

func TestLexSqlDescribe(t *testing.T) {
	/*
		describe myidentity
	*/
	verifyTokens(t, `DESCRIBE mytable;`,
		[]Token{
			tv(TokenDescribe, "DESCRIBE"),
			tv(TokenIdentity, "mytable"),
		})
	verifyTokens(t, `DESC mytable;`,
		[]Token{
			tv(TokenDesc, "DESC"),
			tv(TokenIdentity, "mytable"),
		})
}

func TestLexSqlShow(t *testing.T) {
	/*
		show myidentity
	*/
	verifyTokens(t, `SHOW mytable;`,
		[]Token{
			tv(TokenShow, "SHOW"),
			tv(TokenIdentity, "mytable"),
		})
	verifyTokens(t, `SHOW CREATE TRIGGER mytrigger;`,
		[]Token{
			tv(TokenShow, "SHOW"),
			tv(TokenCreate, "CREATE"),
			tv(TokenIdentity, "TRIGGER"),
			tv(TokenIdentity, "mytrigger"),
		})
	verifyTokenTypes(t, "SHOW FULL TABLES FROM `ourschema` LIKE '%'",
		[]TokenType{TokenShow,
			TokenFull, TokenTables,
			TokenFrom, TokenIdentity, TokenLike, TokenValue,
		})

	/*
	   SHOW TABLES
	   FROM `<yourdbname>`
	   WHERE
	       `Tables_in_<yourdbname>` LIKE '%cms%'
	       OR `Tables_in_<yourdbname>` LIKE '%role%';
	*/
	// SHOW TRIGGERS [FROM db_name] [like_or_where]
	verifyTokens(t, `SHOW TRIGGERS FROM mydb LIKE "tr*";`,
		[]Token{
			tv(TokenShow, "SHOW"),
			tv(TokenIdentity, "TRIGGERS"),
			tv(TokenFrom, "FROM"),
			tv(TokenIdentity, "mydb"),
			tv(TokenLike, "LIKE"),
			tv(TokenValue, "tr*"),
		})
	verifyTokens(t, "SHOW TRIGGERS FROM mydb WHERE `Triggers_in_mydb` LIKE 'tr*';",
		[]Token{
			tv(TokenShow, "SHOW"),
			tv(TokenIdentity, "TRIGGERS"),
			tv(TokenFrom, "FROM"),
			tv(TokenIdentity, "mydb"),
			tv(TokenWhere, "WHERE"),
			tv(TokenIdentity, "Triggers_in_mydb"),
			tv(TokenLike, "LIKE"),
			tv(TokenValue, "tr*"),
		})
	// SHOW INDEX FROM tbl_name [FROM db_name]
	verifyTokens(t, `SHOW INDEX FROM mydb LIKE "idx*";`,
		[]Token{
			tv(TokenShow, "SHOW"),
			tv(TokenIdentity, "INDEX"),
			tv(TokenFrom, "FROM"),
			tv(TokenIdentity, "mydb"),
			tv(TokenLike, "LIKE"),
			tv(TokenValue, "idx*"),
		})
}

func TestLexSqlCreate(t *testing.T) {
	// CREATE {DATABASE | SCHEMA} [IF NOT EXISTS] db_name
	// [create_specification] ...
	verifyTokens(t, `CREATE SCHEMA IF NOT EXISTS mysource 
		WITH stuff = "hello";
		`,
		[]Token{
			tv(TokenCreate, "CREATE"),
			tv(TokenSchema, "SCHEMA"),
			tv(TokenIf, "IF"),
			tv(TokenNegate, "NOT"),
			tv(TokenExists, "EXISTS"),
			tv(TokenIdentity, "mysource"),
			tv(TokenWith, "WITH"),
			tv(TokenIdentity, "stuff"),
			tv(TokenEqual, "="),
			tv(TokenValue, "hello"),
		})

	verifyTokens(t, `CREATE SCHEMA mysource WITH stuff = "hello";`,
		[]Token{
			tv(TokenCreate, "CREATE"),
			tv(TokenSchema, "SCHEMA"),
			tv(TokenIdentity, "mysource"),
			tv(TokenWith, "WITH"),
			tv(TokenIdentity, "stuff"),
			tv(TokenEqual, "="),
			tv(TokenValue, "hello"),
		})

	verifyTokens(t, `CREATE SOURCE mysource WITH stuff = "hello";`,
		[]Token{
			tv(TokenCreate, "CREATE"),
			tv(TokenSource, "SOURCE"),
			tv(TokenIdentity, "mysource"),
			tv(TokenWith, "WITH"),
			tv(TokenIdentity, "stuff"),
			tv(TokenEqual, "="),
			tv(TokenValue, "hello"),
		})

	verifyTokens(t, `CREATE VIEW mysource WITH stuff = "hello";`,
		[]Token{
			tv(TokenCreate, "CREATE"),
			tv(TokenView, "VIEW"),
			tv(TokenIdentity, "mysource"),
			tv(TokenWith, "WITH"),
			tv(TokenIdentity, "stuff"),
			tv(TokenEqual, "="),
			tv(TokenValue, "hello"),
		})

	verifyTokens(t, `CREATE TABLE articles 
		 (
		  ID int(11) NOT NULL AUTO_INCREMENT,
		  Email char(150) NOT NULL DEFAULT '',
		  PRIMARY KEY (ID),
		  CONSTRAINT emails_fk FOREIGN KEY (Email) REFERENCES Emails (Email)
		) ENGINE=InnoDB AUTO_INCREMENT=4080 DEFAULT CHARSET=utf8
	WITH stuff = "hello";`,
		[]Token{
			tv(TokenCreate, "CREATE"),
			tv(TokenTable, "TABLE"),
			tv(TokenIdentity, "articles"),
			tv(TokenLeftParenthesis, "("),
			tv(TokenIdentity, "ID"),
			tv(TokenTypeInteger, "int"),
			tv(TokenLeftParenthesis, "("),
			tv(TokenInteger, "11"),
			tv(TokenRightParenthesis, ")"),
			tv(TokenNegate, "NOT"),
			tv(TokenNull, "NULL"),
			tv(TokenIdentity, "AUTO_INCREMENT"),
			tv(TokenComma, ","),
			tv(TokenIdentity, "Email"),
			tv(TokenTypeChar, "char"),
			tv(TokenLeftParenthesis, "("),
			tv(TokenInteger, "150"),
			tv(TokenRightParenthesis, ")"),
			tv(TokenNegate, "NOT"),
			tv(TokenNull, "NULL"),
			tv(TokenDefault, "DEFAULT"),
			tv(TokenValue, ""),
			tv(TokenComma, ","),
			tv(TokenPrimary, "PRIMARY"),
			tv(TokenKey, "KEY"),
			tv(TokenLeftParenthesis, "("),
			tv(TokenIdentity, "ID"),
			tv(TokenRightParenthesis, ")"),
			tv(TokenComma, ","),
			tv(TokenConstraint, "CONSTRAINT"),
			tv(TokenIdentity, "emails_fk"),
			tv(TokenForeign, "FOREIGN"),
			tv(TokenKey, "KEY"),
			tv(TokenLeftParenthesis, "("),
			tv(TokenIdentity, "Email"),
			tv(TokenRightParenthesis, ")"),
			tv(TokenReferences, "REFERENCES"),
			tv(TokenIdentity, "Emails"),
			tv(TokenLeftParenthesis, "("),
			tv(TokenIdentity, "Email"),
			tv(TokenRightParenthesis, ")"),
			tv(TokenRightParenthesis, ")"),
			tv(TokenEngine, "ENGINE"),
			tv(TokenEqual, "="),
			tv(TokenIdentity, "InnoDB"),
			tv(TokenIdentity, "AUTO_INCREMENT"),
			tv(TokenEqual, "="),
			tv(TokenInteger, "4080"),
			tv(TokenDefault, "DEFAULT"),
			tv(TokenIdentity, "CHARSET"),
			tv(TokenEqual, "="),
			tv(TokenIdentity, "utf8"),
			tv(TokenWith, "WITH"),
			tv(TokenIdentity, "stuff"),
			tv(TokenEqual, "="),
			tv(TokenValue, "hello"),
		})
}
func TestLexSqlDrop(t *testing.T) {
	// DROP {DATABASE | SCHEMA | SOURCE | TABLE} [IF EXISTS] db_name
	verifyTokens(t, `DROP SCHEMA IF EXISTS myschema;`,
		[]Token{
			tv(TokenDrop, "DROP"),
			tv(TokenSchema, "SCHEMA"),
			tv(TokenIf, "IF"),
			tv(TokenExists, "EXISTS"),
			tv(TokenIdentity, "myschema"),
		})
	verifyTokens(t, `DROP TABLE IF EXISTS mytable;`,
		[]Token{
			tv(TokenDrop, "DROP"),
			tv(TokenTable, "TABLE"),
			tv(TokenIf, "IF"),
			tv(TokenExists, "EXISTS"),
			tv(TokenIdentity, "mytable"),
		})
	verifyTokens(t, `DROP SOURCE IF EXISTS mysource;`,
		[]Token{
			tv(TokenDrop, "DROP"),
			tv(TokenSource, "SOURCE"),
			tv(TokenIf, "IF"),
			tv(TokenExists, "EXISTS"),
			tv(TokenIdentity, "mysource"),
		})
	verifyTokens(t, `DROP DATABASE IF EXISTS mydb;`,
		[]Token{
			tv(TokenDrop, "DROP"),
			tv(TokenDatabase, "DATABASE"),
			tv(TokenIf, "IF"),
			tv(TokenExists, "EXISTS"),
			tv(TokenIdentity, "mydb"),
		})
	verifyTokens(t, `DROP DATABASE mydb;`,
		[]Token{
			tv(TokenDrop, "DROP"),
			tv(TokenDatabase, "DATABASE"),
			tv(TokenIdentity, "mydb"),
		})
}
