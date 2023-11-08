# Example-CRUD

A RESTful API example for simple CRUD application with Go

Built with **chi** and **gorm**

### Build and Run

#### Using docker compose

```
docker compose up
```

### Endpoints

#### User Service

| HTTP Method | URL                                                                                       | Description                                  |
|-------------|-------------------------------------------------------------------------------------------|----------------------------------------------|
| `POST`      | http://localhost:8000/users                                                               | Create new User                              |
| `PUT`       | http://localhost:8000/users/{userId}                                                      | Update User by ID                            |
| `GET`       | http://localhost:8000/users/{userId}                                                      | Get User by ID                               |
| `DELETE`    | http://localhost:8000/users/{userId}                                                      | Delete User by ID                            |
| `GET`       | http://localhost:8000/users?name={name}&country={country}&page={page}&pageSize={pageSize} | Search Users by name and country with Paging |

#### POST/PUT body

```json
{
  "first_name": "Alice",
  "last_name": "Bob",
  "nickname": "AB123",
  "password": "supersecurepassword",
  "email": "alice@bob.com",
  "country": "uk"
}
   ```

#### RabbitMQ

Sends message with user id to RabbitMQ corresponding queues on every user create/update/delete event
