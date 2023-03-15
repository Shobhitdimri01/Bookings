#for linux/Mac
#!/bin/bash
go build -o PROJECT_BOOKING ./cmd/web/*.go
./PROJECT_BOOKING -dbname=bookings -dbuser=postgres -dbpass=postgres -cache=false -production=false