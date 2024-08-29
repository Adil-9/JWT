# JWT token auth using gin framework

Foobar is a Python library for dealing with word pluralization.

## Deployment


```bash
docker compose up
```
there could be a trouble after deploying that application exits before it could make connection to the db. In that case start application container 
```bash
docker start jwt-app
```


## Usage

connect to 
### localhost:8081/access 
to get jwt token and 
### localhost:8081/refresh 
to refresh token
```
