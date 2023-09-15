# Startup

Below targets should start all required containers: 

make start-with-mongo

or

make start-with-postgres


## Test usage

### Retrieve all available books

curl http://localhost:3000/book

### Retrieve one book

curl http://localhost:3000/book/{id}

### Create new book

curl -X POST http://localhost:3000/book -d '{"name":"Example Book","author":"Some Author"}'

### Update book

curl -X PUT http://localhost:3000/book/{id} -d '{"name":"Example Book","author":"Some Author"}'

### Delete book

curl -X DELETE http://localhost:3000/book/{id}

### Delete all books

curl -X DELETE http://localhost:3000/book