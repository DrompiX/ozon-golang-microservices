# HW - Microservices

## Part 1.1 - Ozon Design
[Link](https://miro.com/app/board/uXjVOxOVzUE=/?moveToWidget=3458764526360704908&cot=14) to the miro for full Ozon HLD

## Part 1.2 - Detailed microservices
[Link](https://miro.com/app/board/uXjVOxOVzUE=/?moveToWidget=3458764526361236575&cot=14) to the miro for selected detailed microservices

## Что сделано
- Реализовано 3 микросервиса - создание заказа, проверка наличия на складе и оплата
- Реализована хореографическая сага
- Общение осуществляется через kafka, данные в postgresql

## Было скошено много углов, чтобы успеть, часть из них ниже
- Только один инстанс базы данных, но с разными базами внутри, чтобы не поднимать пока на каждый сервис отдельный контейнер с базой
- Предполагается, что в базе уже есть item'ы
