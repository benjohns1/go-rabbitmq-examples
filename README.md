# Go RabbitMQ Examples
Simple examples and sample implementations
Need to have Docker and Compose
## Marco Polo Example
Basic direct send using a default exchange and direct queue routing
### Run
```sh
git clone https://github.com/benjohns1/go-rabbitmq-examples.git
cd go-rabbitmq-examples
git checkout marco-polo
docker-compose build && docker-compose up
```
 - After the containers are spun up, go to http://localhost:8080/marco and watch your console  
 - You can watch the queue in the RabbitMQ admin console at http://localhost:15672, login with 'guest' 'guest'