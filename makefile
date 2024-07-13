build:
	@docker build --platform linux/amd64 -t lambda:test .

run:
	@docker run --env-file .env -p 8080:8080 lambda:test

test-book-order:
	@curl -X POST "http://localhost:8080/bookOrder" -d '{"order_id": "1234", "rider_id": "5678"}'

test-insert-order:
	@curl -X POST "http://localhost:8080/api/orderQueue" -d '{"order_id": "1234", "rider_id": null}'

test-delete-order:
	@curl -X DELETE "http://localhost:8080/api/orderProcess?order_id=1234"