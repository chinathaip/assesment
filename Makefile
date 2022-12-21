run:
	DATABASE_URL=postgres://hhbonlde:FHtwcrPFGiSIONbiqMkSjAuHVFXIWmmX@floppy.db.elephantsql.com/hhbonlde PORT=:2565 go run .


docker :
	docker run -e DATABASE_URL=postgres://hhbonlde:FHtwcrPFGiSIONbiqMkSjAuHVFXIWmmX@floppy.db.elephantsql.com/hhbonlde -e PORT=:2565 -p 2565:2565 kbtg-go-assesment

it test:
	go clean -testcache && go test -v --tags=integration ./...