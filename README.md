# GO & Django Communication using gRPC

This project is a showcase on how gRPC can be used for comunication between serverice. In this case a service built with golang and another built with Django.

![django-go drawio](https://github.com/helissonomc/go-django-communication/assets/60279210/35e6139b-6515-45d0-ad05-4f92c8b8bd59)

We have a Restful API on go server, where we can perform a CRUD of an User

On Django Side we have a table called externaluser where it is a replica of the table user on Go database, it does not store the password though.

## How to run the project
run:
```
docker-compose up
```

wait a few seconds because the database containers takes a few seconds to be ready, then check if all the servers are up and running

run:
```
curl --location 'http://localhost:8080/users' \                                                                                                   ─╯
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "test@example.com",
    "name": "test name ",
    "password": "test123"
}'
```
This will create the user in the golang service where you can check it by running:
```
curl --location 'http://localhost:8080/users'
```
and  also it will create a gRPC call to the Django gRPC server and it will also persist the data in the table `ExternalUser`
You can check it in `http://localhost:8081/admin/users/externaluser/` user: `admin` password: `adminpassword`
![image](https://github.com/helissonomc/go-django-communication/assets/60279210/082f86e5-ec12-4b67-9999-dbc5d2de01d9)
