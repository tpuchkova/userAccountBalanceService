run:
	docker-compose up todo-app

migrate:
	 migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable'
