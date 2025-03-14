CONFIG_PATH=./config/local.yaml
MAIN=./cmd/url-shortener/main.go
SERVER_PATH=http://localhost:8082

run: 
	CONFIG_PATH=$(CONFIG_PATH) go run $(MAIN) &
test:
	go test ./internal/http-server/handlers/url/save -v; 
	go test ./internal/http-server/handlers/url/redirect -v
go_to:
	curl $(SERVER_PATH)/$(ALIAS)
add:
	curl -X POST -H "Content-Type: application/json" -d '{"url":"$(URL)", "alias":"$(ALIAS)"}' -u admin:12345  $(SERVER_PATH)/url/save
delete:
	curl -X DELETE -H "Content-Type: application/json" -d '{"alias":"$(ALIAS)"}' -u admin:12345 $(SERVER_PATH)/url/delete
kill:
	killall go; killall main