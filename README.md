# GeneratorDoc
Запускается сервер gin , на ружу два endpoint , один для получения шаблона  и добавления данных в шаблон.
Другой для скачивания файлов.

Сервер работает на localhost:8080

curl --silent --location --request POST 'http://localhost:8080/gendoc' \
--header 'Content-Type: application/json' \
--data-raw '{
"URLTemplate": "https://sycret.ru/service/apigendoc/forma_025u.xml",
"RecordID": 30
}'
