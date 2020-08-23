# TWljaGFsLUtvbW9y
Zadanie rekrutacyjne GWP

folder item:

item.go - struktury i funkcje wspólne dla serwera i workera

item_test.go - plik testowy

folder server:

server.go - obsługa serwera

server_test.go - plik testowy

folder worker:

worker.go - obsługa workera

worker_test.go - plik testowy

przykładowe komendy z poziomu konsoli:

curl -si 127.0.0.1:8080/api/fetcher -X POST -d "{\"url\": \"https://httpbin.org/range/1\",\"interval\":10}"

curl -si 127.0.0.1:8080/api/fetcher -X POST -d "{\"url\": \"https://httpbin.org/delay/4\",\"interval\":10}"

curl -si 127.0.0.1:8080/api/fetcher -X POST -d "{\"url\": \"https://httpbin.org/delay/8\",\"interval\":10}"

curl -si 127.0.0.1:8080/api/fetcher

curl -si 127.0.0.1:8080/api/fetcher/2

curl -si 127.0.0.1:8080/api/fetcher/2/history

curl -si 127.0.0.1:8080/api/fetcher/2 -X DELETE

curl -si 127.0.0.1:8080/api/fetcher/3 -X PUT -d "{\"url\": \"https://httpbin.org/delay/3\",\"interval\":8}"




