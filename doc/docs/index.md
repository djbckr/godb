# GoDB #
GoDB is a SQL RDBMS written in Go, and is unique in that all interaction is performed through HTTPS.

## The Basics ##
All interaction with GoDB is done through HTTPS (port 443). GoDB supports keep-alive connections as well as stand-alone
requests. Except for the initial login, all requests must have an Authorization header with two pieces of content:
`AuthToken` and `SessionID`. These values are separated by space and are alphanumeric in type.

    Authorization: 'AuthToken xxxxxxx SessionID xxxxxxxx'

If using keep-alive connections, it is valid to have various session/auth tokens in each request. The only thing
keep-alive does is reduce connection and handshake overhead (which can be substantial).

The server will respond with the following possible HTTP response codes:

- 200 OK - The call completed successfuly with no errors.
- 400 Bad Request - This is typically in response to invalid SQL. The response body will contain a message indicating what the problem might be.
- 401 Unauthorized - This response means that the user needs to login, either because they never logged-in, their session expired, or their session was terminated by an administrator. Or if logging in, an invalid login combination was specified.
- 403 Forbidden - This means that the user tried to access a resource (table, view, etc) that they are not allowed to access.
- 404 Not Found - A request to a non-existent end-point.
- 405 Method not allowed - Example: A GET operation was requested, but the endpoint only allows PUT/POST. 
- 409 Conflict - More than one /sql request was made on the same SessionID. This is not allowed.
- 415 Unsupported Media Type - GoDB currently only supports application/json and application/XML media types.
- 422 Unprocessable Entity - The content of the request body expected certain data, but it was not there. The response body will contain a message indicating what the problem might be.
- 500 Internal Server Error - this will happen when GoDB has an internal error. The response body should contain a message indicating what the problem might be.

In addition to the above HTTP response codes, the server will always respond with `code` and `message` element/attribute.
You can look up the error codes [here](err/index.md)

It is important to note the difference between `AuthToken` and `SessionID`. The `AuthToken` is a user's login authorization:
What the user is allowed to do when using the database. This determines whether the user is allowed to query a certain
table or view, execute a certain stored procedure, or update/insert a particular table. The `SessionID` is one user's
particular login session. One session is limited to one executing routine at one time. That is, you cannot have two
https connections make simultaneous requests using the same SessionID. Doing so could violate the integrity of
transactions, so this is not allowed. You can, however, use an AuthToken to create a new session. This will be described below.

### `AuthToken` Lifetime ###
`AuthToken` lifetime is determined by the GoDB administrator. It's generally recommended limiting an AuthToken to 10 hours
and this is the default. When an `AuthToken` expires, a new login is required. Batch processes and application server
processes typically last longer than this, and they should be granted unlimited `AuthToken` lifetime. However, when dealing
with a live user, the `AuthToken` and `SessionID` should be directly related to their activity.

### `SessionID` Lifetime ###
The `SessionID` lifetime is determined by the last activity recorded. Keeping sessions open while inactive is a waste of resources.
If there has been no activity on the session for a given time period, the session is closed and must be reopened by /login again.
The session lifetime is set by the GoDB administrator. The default is 10 minutes.

## JSON and/or XML ##
GoDB supports JSON and XML based on the `Accept` and `Content-Type` headers. For example, if your request contains:

    Content-Type: application/json
    Accept: application/xml

You would send the body of the request in JSON, and the response from the server would be in XML format. Likewise:

    Content-Type: application/xml
    Accept: application/json

The request body would be in XML, and the server response would be JSON. As of this writing, these are the only supported
formats. We will consider binary formats in the future.

Both the JSON and XML formats are intended to stream. This means there is a certain order of the data structure that is
required. See the [/sql endpoint](ep/index.md#execute-sql-statements-sql) for more information on how that works.

Note that all XML is in a root `<godb>` element. JSON does not use a root element.
