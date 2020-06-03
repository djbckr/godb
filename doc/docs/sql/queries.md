# Queries #

The most basic query that most people might start with is:
```sql
SELECT 'Hello World'
```

This will return a single row with one column, and the contents "Hello World" in it.

GoDB offers some flexibility on how you might structure a query. Typical SQL is in line with
the English-like syntax that SQL was originally designed.
```sql
SELECT field1, field2 FROM my_table
```
An extension to this reverses these clauses to "behave more like a computer" where the source is brought
into scope before accessing the data in it.
```sql
FROM my_table SELECT field1, field2
```

If you wish to always have a one-row source table in your queries, use the § table.
Note the table name is a UTF8 character. This table has no columns so you have to provide something.
```sql
FROM § SELECT 'Hello World'
```
Either formats are equivalent and acceptable to GoDB.

-- TODO -- lots more about select/with/from
