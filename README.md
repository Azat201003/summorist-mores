# summorist-mores

This is microservice for [summorist](https://github.com/Azat201003/summorist-shared) application

## Run

I supposed you will run all services together, so just download [devops repo](https://github.com/Azat201003/summorist-devops) and read its instruction, but I describe how it should be ran anyway

### Env

You need to pass some environment variables:
 - `MORES_FILE_PREFIX` begin of file path (*prefix*[id]*suffix*)
 - `MORES_FILE_SUFFIX` end of file path
 - `MORES_HOST` obviously host of mores (like 127.0.0.1)
 - `MORES_HOST` just port (like 0-65535 or something)
 - `MORES_POSTGRES_DSN` dsn for postgresql database

### Database

You need to start postgresql for connection by env data, passed before, it is easy read [docs](https://www.postgresql.com)
Also you need to run migrations (psql -f *file*) or similar, also in the documentation

### Docker

After all you can just build and run container like this:
``` bash
docker build -t summorist-mores .
docker run --env-file <where are your envs> -p <port>:<port> summorist-mores
```
