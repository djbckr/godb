package main

/*

	- pages: 8K or possibly multi-size capable
	- numbers: up to 64000 decimal places
  - text/varchar/char: up to 64000 characters
  - LOB: "unlimited" size
  - cache using LIRS

DML:
  - query: SELECT/WITH/FROM...
  - INSERT...
  - UPDATE...
  - UPSERT/MERGE...
  - DELETE...

TRX CTL:
  - START...
  - COMMIT...
  - SAVEPOINT...
  - ROLLBACK...

CODE BLOCKS:
  - BEGIN/DECLARE...END;

DDL:
  - COMMENT
  - CREATE
  - ALTER
  - DROP
  - GRANT
  - REVOKE

Object Types:
  ? CONTEXT
  ? AGGREGATE
  DOMAIN
  DATABASE
  FUNCTION
  INDEX
  MVIEW
  PACKAGE [body]
  PROCEDURE
  ROLE
  SCHEMA
  SEQUENCE
  [session]
  [system]
  SYNONYM
  TABLE
  TABLESPACE
  TRIGGER
  TYPE (SqlType)
  USER
  VIEW

*/
