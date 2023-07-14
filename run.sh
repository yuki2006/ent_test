docker-compose down
docker-compose up -d --wait
go build -o a.out .
./a.out
