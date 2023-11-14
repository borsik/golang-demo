# Example-CRUD

A RESTful API example for simple CRUD application with Go

Built with **chi**

### Build and Run

#### Using docker compose

```
docker compose up
```

### Endpoints

#### User Service

| HTTP Method | URL                                                                                        | Description                                  |
|-------------|--------------------------------------------------------------------------------------------|----------------------------------------------|
| `POST`      | http://localhost:8000/users                                                                | Create new User                              |
| `PUT`       | http://localhost:8000/users/{userId}                                                       | Update User by ID                            |
| `GET`       | http://localhost:8000/users/{userId}                                                       | Get User by ID                               |
| `DELETE`    | http://localhost:8000/users/{userId}                                                       | Delete User by ID                            |
| `GET`       | http://localhost:8000/users?name={name}&country={country}&page={page}&page_size={pageSize} | Search Users by name and country with Paging |

#### POST/PUT body

```json
{
    "first_name": "Alice",
    "last_name": "Bob",
    "nickname": "AB123",
    "password": "supersecurepassword",
    "email": "alice@bob.com",
    "country": "US"
}
   ```

### User entity JSON
```json
{
    "id": "1ed8bd54-cb7b-4a42-9fd2-ddc62835b755",
    "first_name": "Alice",
    "last_name": "Bob",
    "nickname": "AB123",
    "email": "alice@bob.com",
    "country": "US",
    "created_at": "2023-11-13T11:49:14.025679Z",
    "updated_at": "2023-11-13T11:49:14.025679Z"
}
```

#### RabbitMQ

Sends message with user id to RabbitMQ corresponding queues on every user create/update/delete event
