# go-chat-bot
## Table of contents
* [General info](#general-info)
* [Features](#features)
* [Technologies](#technologies)
* [Prerequisites](#prerequisites)
* [Setup](#setup)
* [Usage](#usage)

---

## General info
Browser-based chat application in Go. This app allows several users to talk in multiple chatrooms and also to get stock quotes from an API using a specific command.

---

## Features
* Allow registered users to log in and talk with other users in a chatroom.
* Allow users to post messages as commands into the chatroom with the following format **/stock=stock_code**.
* Decoupled bot that calls an API using the stock_code as a parameter (https://stooq.com/q/l/?s=aapl.us&f=sd2t2ohlcv&h&e=csv, here `aapl.us` is the stock_code).
* The bot parses the received CSV file and then send a message back into the chatroom using RabbitMQ. The message is a stock quote
with the following format: “APPL.US quote is $93.42 per share”. The post owner of the message is the bot.
* Chat messages ordered by their timestamps and show only the last 50 messages.
* Have more than one room.
* Unit tests for the `bot`.
* Messages that are not understood or any exceptions raised within the bot are handled.

---

## Technologies
The project is created with or uses:

* HTML
* JS
* Go
* Mariadb
* RabbitMQ
* Docker

---

## Prerequisites
* Docker Desktop (Download from [here](https://www.docker.com/products/docker-desktop))

## Setup
Follow the next steps to run this project locally:

1. Open the project solution on your IDE of preference, look inside the **GoChallenge** folder for the `.env`and update the connection string if need it.

2. Open Powershell or Bash and run the next command to start the containers.
```
docker-compose up -d
```
---

## Usage
Once the application is running, you just need to register as a user and login into the app to access the chatroom.

The application will be serving on http://localhost:8080.

