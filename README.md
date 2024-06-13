# oauth2-server

## Run the Server
```bash
    make run-local
```

## Request a token

```bash
    curl -X POST -u client-id:client-secret http://localhost:8080/token
```

response will be smth like this:

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTgyODI1ODMsImlzcyI6InNvbWVfc2VydmVyX2lkIiwic3ViIjoiY2xpZW50LWlkIn0.NCsbibM1A8T4V7uagX2vB1jHbv1Y_yLw0JTqhypZpcggIVIpWlLwzy2NRPkGm0UrnX6VIdRmAhoaY6Le2OOz5KirHpcO9dK6hKJUxgm9m5aklCbExaIkDLsBuejHF1yJbY2vvs-V9XnfOq53gf1griZx9LcZ-bETD98r1gcjPOaiHjBmrSyAE7UqYnoT4iyPibb09L6cIJZx8rYi0M2j7MiImJF5B6t6gct9f0SZoEhTGytx7s3sL8qB72k7IuJ6hkBcCJkQI8ItzicZ349piK7zV8jAgPAcoLxyMYn0KIuYP52WxrQ9OKOIQGQ-kj-6wWwPXiOn7R4enA1p1WBz2A",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

## Validate token

```bash
    curl --location 'localhost:8080/protected' \
--header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTgyODI1ODMsImlzcyI6InNvbWVfc2VydmVyX2lkIiwic3ViIjoiY2xpZW50LWlkIn0.NCsbibM1A8T4V7uagX2vB1jHbv1Y_yLw0JTqhypZpcggIVIpWlLwzy2NRPkGm0UrnX6VIdRmAhoaY6Le2OOz5KirHpcO9dK6hKJUxgm9m5aklCbExaIkDLsBuejHF1yJbY2vvs-V9XnfOq53gf1griZx9LcZ-bETD98r1gcjPOaiHjBmrSyAE7UqYnoT4iyPibb09L6cIJZx8rYi0M2j7MiImJF5B6t6gct9f0SZoEhTGytx7s3sL8qB72k7IuJ6hkBcCJkQI8ItzicZ349piK7zV8jAgPAcoLxyMYn0KIuYP52WxrQ9OKOIQGQ-kj-6wWwPXiOn7R4enA1p1WBz2A'
```
> Note: any access_token from POST token should work here.
