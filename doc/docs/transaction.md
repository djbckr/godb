# Transaction Control #

GoDB provides very flexible transaction control. Using the `Trx-xxx` headers allows you to embed transaction control
into /sql calls, but you can just as easily run separate SQL transaction statements if you find it easier to understand.

A common mistake is to commit after every operation. This is not only wasteful (uses more database resources to commit
and start many small transactions rather than fewer large ones), but circumvents the idea of keeping related operations
atomic. A trivial example is to transfer money from one account to another. This minimally involves two steps: Remove
money from one account and add money to the other account. There may be other steps involved that may include journaling
or logging. All of these operations should take place within one transaction and committed at the end. If any of the
operations fails for some reason, rolling back will return the database to its prior state before that transaction started.

## Starting a Transaction ##
Transactions can be started in one of several ways:

1. By issuing a SQL statement that writes to the database. This includes INSERT, UPDATE, DELETE, UPSERT and MERGE
statements, as well as "SELECT ... FOR UPDATE", and executing stored procedures that issue any of those statements.
The transaction will begin with read/write capabilities, isolation level at read-committed, and a system generated
transaction name.

2. By issuing a `START TRANSACTION ...` SQL command.

3. By providing a `Trx-Start` header on the /sql endpoint, optionally providing read/write or read-only,
isolation level, and transaction name. This is used with the same syntax as the "START TRANSACTION ..." SQL statement.
The transaction will start _before_ executing the SQL.

If a transaction is already running, doing any of the above will be ignored.

## Savepoint ##
A transaction savepoint can be done one of two ways:

1. By issuing a `SAVEPOINT name` SQL command.

2. By providing a `Trx-Savepoint` header on the /sql endpoint. You must provide a name. The savepoint will occur
_before_ the SQL has executed.

If a savepoint is requested but there is no active transaction, a new transaction is started, then the savepoint is
established. You can provide transaction details using the `Trx-Start` header in addition to this header.

## Rollback ##
A rollback can be done one of two ways:

1. By issuing a `ROLLBACK {TO savepoint}` SQL command.

2. By providing a `Trx-Rollback` header, optionally providing the savepoint name. By doing this, you tell GoDB that if
there is an error executing the SQL, execute the rollback. If there is no error, no rollback is performed. If the
header is not provided and an error occurs, the client must decide what to do.

If a rollback is requested but there is no active transaction, nothing happens.

Providing a `Trx-Rollback` header along with a `Trx-Savepoint` header allows batch SQL to run and rollback if an error
occurs.

## Commit ##
Commit can be done one of two ways:

1. By issuing a "COMMIT" SQL command.

2. By providing a `Trx-Commit` header. By doing this, you tell GoDB to commit _after_ the SQL has executed successfully.
If there is an error, the transaction is _not_ committed. If the `Trx-Rollback` header is not provided,
the client must decide what to do.

If a commit is issued but there is no active transaction, nothing happens.
