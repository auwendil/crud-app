make start-with-mongo

make start-with-postgres


curl http://localhost:3000/book

curl http://localhost:3000/book/{id}

curl -X POST http://localhost:3000/book -d '{"name":"Example Book","author":"Some Author"}'

curl -X PUT http://localhost:3000/book/{id} -d '{"name":"Example Book","author":"Some Author"}'

curl -X DELETE http://localhost:3000/book/{id}

curl -X DELETE http://localhost:3000/book