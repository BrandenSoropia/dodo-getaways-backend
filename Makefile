check_mongo:
	ps aux | grep -v grep | grep mongod

dev_start: mongo_start server_start

server_start:
	go run main.go

mongo_start:
	brew services start mongodb-community@4.2

mongo_stop:
	brew services stop mongodb-community@4.2

load_mock_data:
	mongoimport --db test --collection islands --drop --file ./db/__mocks__/islands.json
	mongoimport --db test --collection owners --drop --file ./db/__mocks__/owners.json

	