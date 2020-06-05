# Server Endpoints #

GoDB has very few http endpoints listed below.

## Interactive Database UI (`/`) ##
The `/` endpoint is an interactive web-interface to the database. (TBD more docs)

## User Authentication (`/authenticate`) ##
The `/authenticate` endpoint is used to authenticate (login) a user and start a session.
It supports POST or PUT equivalently. GET and HEAD are not supported.
The body consists of either a username/password, or `AuthToken` value, but not both.

* JSON + Username/Password:
```
{
  "username": "john doe",
  "password": "my fancy password"
}
```
* JSON + AuthToken:
```    
{
  "AuthToken": "8d8d766c-4fc4-493b-b82b-e83096b7e8d8"
}
```
* XML + Username/Password:
```    
<?xml version="1.0" encoding="UTF-8"?>
<godb>
  <username>john doe</username>
  <password>my fancy password</password>
</godb>
```    
* XML + AuthToken:
```    
<?xml version="1.0" encoding="UTF-8"?>
<godb>
  <AuthToken>8d8d766c-4fc4-493b-b82b-e83096b7e8d8</AuthToken>
</godb>
```
Note that all XML is in a root `<godb>` element.

If the login is successful, the server responds with a 200 and the `Authorization` header you will use for subsequent
requests to the server. There is no body in the response.

If the login is unsuccessful, the server will return a 401 Unauthorized response. There is no body in the response.

Providing the `AuthToken` begins a new session that is functionally separate from the existing session, as if two
different users were logged-in.

## Execute SQL Statements (`/sql`) ##
The `/sql` endpoint allows one to execute any valid SQL statement.
The request supports POST or PUT equivalently. GET and HEAD are not supported.
The request body consists of a `sql` or `sqlid`, `meta`, and `data` section. Each of these sections are optional depending on what you
are doing. The server responds with `code`, `message`, `meta` and `data` sections as appropriate.

Note: Both the JSON and XML formats are intended to be "streamed" - this means that `data` should always appear last in
the document. This allows the metadata to be parsed so that both the server and client can correctly interpret the
potentially large amounts of data passed from one to the other. It is not outside the realm of possibilities to send
gigabytes of data at a time. Doing so should of course use the chunking protocol available in HTTP.

We will go through several examples with JSON and XML.

### EXAMPLE: Change my Password ###
* JSON
```
{
  "sql": "alter user set password q[my new fancy password]"
}
```
* XML
```
<?xml version="1.0" encoding="UTF-8"?>
<godb>
   <sql>alter user set password "my new fancy password"</sql>
</godb>
```

Note that string delimiters can be one of several types, both for flexibility and safety. These will be noted in the
String Constants chapter.

The server will respond to the above with a `message` section:
* JSON
```
{
  "code": 0,
  "message": "Password successfully changed"
}
```
* XML
```
<?xml version="1.0" encoding="UTF-8"?>
<godb>
   <code>0</code>
   <message>Password successfully changed</message>
</godb>
```

There is no need for a `meta` or `data` in either the request or response.

### EXAMPLE: Run a Query ###
* JSON
```
{
  "sql": "select * from some_table where updated >= :updated",
  "meta": {
    "binds": [
      {
        "name": "updated",
        "type": "timestamp",
        "format": "YYYY-MM-DDTHH:mm:SS"
      }
    ]
  },
  "data": {
    "updated": "2020-01-01T00:00:00"
  }
}
```
* XML
```
<?xml version="1.0" encoding="UTF-8"?>
<godb>
   <sql>select * from some_table where updated &gt;= :updated</sql>
   <meta>
      <binds>
         <updated type="timestamp" format="YYYY-MM-DDTHH:mm:SS" />
      </binds>      
   </meta>
   <data>
      <updated>2020-01-01T00:00:00</updated>
   </data>
</godb>
```
Note a few things:

1. The bind variable `:updated` can alternatively be a question mark and be referenced positionally, starting at 1. In that
case it would be referenced as `"1"` or `<1>` in the JSON or XML as appropriate.

2. In the XML version, the `>=` must be escaped using `&gt;=` __OR__ if you wish to avoid doing that, you can embed
a CDATA inside the `<sql>` tag with unescaped SQL inside.

3. In the case of a query as above, you can only have one data section. If you were to execute an INSERT statement
multiple times, the JSON version can have an array of data objects, and the XML version can have multiple `<data>`
elements. GoDB will execute the INSERT statement as many times as there are elements in the array, or as many times
as there are `<data>` sections.

The return from the server will be as follows:

* JSON
```
{
  "code": 0,
  "message": "Success",
  "meta": {
    "fields": [
      {
        "name": "field1",
        "type": "string"
      },
      {
        "name": "field2",
        "type": "number"
      },
      {
        "name": "updated",
        "type": "timestamp",
        "format": "YYYY-MM-DDTHH:mm:ss"
      },
      {
        "name": "fieldX",
        "type": "enum"
      }
    ],
    "sqlid": "84c8e7f9-b6d5-4a9b-a1ba-c7608efacd8b"
  },
  "data": [
    {
      "field1": "this is a string",
      "field2": 234567.293,
      "updated": "2020-02-02T11:23:33",
      "fieldX": "png"
    },
    {
      "field1": "another value",
      "field2": 3.14,
      "updated": "2020-03-02T11:23:33",
      "fieldX": "jpeg"
    }
  ]
}
```
* XML
```
<?xml version="1.0" encoding="UTF-8"?>
<godb>
   <code>0</code>
   <message>Success</message>
   <meta sqlid="84c8e7f9-b6d5-4a9b-a1ba-c7608efacd8b">
      <fields>
         <field name="field1" type="string" />
         <field name="field2" type="number" />
         <field name="updated" type="timestamp" format="YYYY-MM-DDTHH:mm:SS" />
         <field name="fieldX" type="enum" />
      </fields>
   </meta>
   <data>
      <field1>this is a string</field1>
      <field2>234567.293</field2>
      <updated>2020-02-02T11:23:33</updated>
      <fieldX enum="2">png</fieldX>
   </data>
   <data>
      <field1>another value</field1>
      <field2>3.14</field2>
      <updated>2020-03-02T11:23:33</updated>
      <fieldX enum="3">jpeg</fieldX>
   </data>
</godb>
```
In subsequent executions, you can provide `sqlid` instead of the `sql` text. This allows the server to bypass the
parse and plan stages of execution, reducing overhead.


## Administrative Tools (`/admin`) ##
TBD (99% of all work will be done in SQL)

