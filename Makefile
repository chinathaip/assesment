run:
	DATABASE_URL=postgres://hhbonlde:FHtwcrPFGiSIONbiqMkSjAuHVFXIWmmX@floppy.db.elephantsql.com/hhbonlde PORT=:2565 go run .


docker :
	docker run -e DATABASE_URL=postgres://hhbonlde:FHtwcrPFGiSIONbiqMkSjAuHVFXIWmmX@floppy.db.elephantsql.com/hhbonlde -e PORT=:2565 -p 2565:2565 kbtg-go-assesment

unittest :
	go clean -testcache && go test -v --tags=unit ./...

ittest:
	docker-compose -f docker-compose.yml down && docker-compose -f docker-compose.yml up --build --abort-on-container-exit --exit-code-from integration_tests