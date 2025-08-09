Система рейтинга продуктов на gRPC

проект представляет собой систему для отслеживания просмотров продуктов и расчета рейтингов

Система состоит из 4 микросервисов:

Бэкенд-сервис: gRPC сервер, который получает события просмотров продуктов и публикует их в Kafka

Сервис рейтинга: Потребляет события из Kafka и обновляет рейтинги в Redis

Дашборд: Используется для просмотра продуктов из Redis

Веб-клиент: Простой веб-интерфейс для просмотра продуктов

--------------------------------------------------------------------------------

Стек технологий:

Go

gRPC 

Kafka 

Redis 

Docker

Docker-compose

--------------------------------------------------------------------------------

Ehv:

REDIS_PORT=

REDIS_HOST=r

KAFKA_BROKERS=

KAFKA_TOPIC=

KAFKA_GROUP_ID=

GRPC_SERVER=

GRPC_PORT=

RATING_PORT=

DASHBOARD_PORT=

WEB_PORT=

--------------------------------------------------------------------------------

Для старта:

      docker-compose up -d --build

-----------------------------------------------------------------------------------------

Product Rating System with gRPC

This project is a system for tracking product views and calculating ratings.

System Components

The system consists of 4 microservices:

Backend Service: gRPC server that receives product view events and publishes them to Kafka

Rating Service: Consumes events from Kafka and updates ratings in Redis

Dashboard: Used to view products from Redis

Web Client: Simple web interface for viewing products

Technology Stack

Go

gRPC

Kafka

Redis

Docker

Docker-compose

--------------------------------------------------------------------------------

Ehv:

REDIS_PORT=

REDIS_HOST=r

KAFKA_BROKERS=

KAFKA_TOPIC=

KAFKA_GROUP_ID=

GRPC_SERVER=

GRPC_PORT=

RATING_PORT=

DASHBOARD_PORT=

WEB_PORT=

--------------------------------------------------------------------------------

Start: 

      docker-compose up -d --build
        

