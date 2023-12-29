# sack
An in-memory key/value store that can be used with clients that supports redis. This is not a production ready application but more like redis. This is me doing some heavy stuffs in Golang :) 

# Installation Instructions
New released version is pushed to dockerhub with tags. Just pull that image from dockerhub and run the container.

``` sh
docker pull neymarsabin/sack:latest
docker run -p 6379:6379 neymarsabin/sack:latest
```
The port exposed is *6379*, so make sure you have the redis server stopped. You can use the redis-cli to connect to `sackDB` server. 

# Supported Commands
Only basic commands supported at the moment. 
 - SET -> set a key and value 
 ```sh 
 SET name neymar
 ```
 - GET -> get a value when key passed as arguments
 ```sh
 GET name 
 ```
 - HSET -> set a key to object value
 ```sh
 HSET users u1 neymar
 HSET users u1 sabin
 ```
*`HSET` does not support nested objects at the moment.*
 - HGET -> get value from the objects
 ```sh
 HGET users u1
 HGET users u2
 ```
 
 More commands and optimizations are on the way. 
