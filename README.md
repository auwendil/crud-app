# Startup

## Run go app

`make run` will start application locally

These targets should start all required containers: 

`make start-with-mongo`

or

`make start-with-postgres`


## Test usage

### Retrieve all available books

`curl http://localhost:3000/book`

### Retrieve one book

`curl http://localhost:3000/book/{id}`

### Create new book

`curl -X POST http://localhost:3000/book -d '{"name":"Example Book","author":"Some Author"}'`

### Update book

`curl -X PUT http://localhost:3000/book/{id} -d '{"name":"Example Book","author":"Some Author"}'`

### Delete book

`curl -X DELETE http://localhost:3000/book/{id}`

### Delete all books

`curl -X DELETE http://localhost:3000/book`


## Improvements

### More tests

There may be more negative case tests to better check for edge case scenarios.

### Swagger API

Manual testing API endpoints will be easier with swagger.

### Configuration files

Use config files to provide application configuration


## Alternatives

### internal.repository

With bigger repositories maybe will be better to put each repo into separate packages. 

### API routes

API routes may be setup in one separate *routes.go* file or in multiple files for each path/subpath.
