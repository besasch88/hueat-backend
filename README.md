# Hueat

## 📋 About
Hueat is a point-of-sale system that helps huts take orders and manage them more smoothly. Waiters can take orders directly at tables on their devices, and orders get sent to the right places: the kitchen, bar, or serving area. It helps reduce mistakes, saves time, and makes service faster. The system also tracks what's being ordered and used, helping with inventory management.

### Install GO

First of all, let's install go version `1.25.0` or higher from this link: https://go.dev/doc/install

### Check the Go version

Before proceed, ensure your version is correct. Run this command in your terminal:

```sh
go version
```

The answer should be something like this according to your installed version and arch:

```sh
go version go1.25.0 darwin/amd64
```

### Start external services

Navigate in the `build` folder and start the Postgres DB and Redis inside Docker:

```sh
cd build
docker compose up hueat-database  -d
```

It contains a PostgresQL database server mapped on the local port `54322`. Feel free to take a look to the docker-compose file to retrieve credentials if you want to use an external tool to connect with.

### Migration Tool

The Migration Tool is a command that help you in creating migrations, apply or revert thanks to migration versioning. Let's start by installing the migration tool:

```sh
brew install golang-migrate
```

and with the following command you can create your first migration:

```sh
migrate create -ext sql -dir ./scripts/migrations -seq init schema
```

Thanks to it, the tool will create two empty sql files in the `scripts/migrations` folder to apply a new changes to the Database or to revert it.
Once your migrations are defined, you can apply them locally with this command:

```sh
migrate -path "./scripts/migrations" -database "postgres://hueat:iUkcBQj2o_PpfX*uaXx7@127.0.0.1:54322/hueat?sslmode=disable" up
```

or just update the DB credentials in the file and run it as a shortcut:

```sh
./scripts/migrate-local.sh
```

Looking to the docker-compose file, you will notice that there is a dedicated service aims to apply migrations each time the project is deployed in your production environment. Basically it starts, applies all the migrations and shutdown.

### Start the webapp locally

Now we have all the migration setup, the DB running and updated and we can run your local webapp locally via this command:

```sh
go run cmd/webapp/main.go
```

If everything is fine, you will see in logs that the webapp is up and running, waiting incoming API requests.

### Start dockerized application

If you want. you can run the webapp application in docker, useful for testing/demo purposes. So from the root folder of your project run:

```sh
bash build/scripts/dev/start.sh
```

The webapp is mapped on the port `8001`.

### Test the webapp

To test the webapp, please open Postman and call this endpoint:

```
POST http://0.0.0.0:8001/api/v1/health-check
```

### Env variables

This project is configured via environment variables that are declared and expected in the repository.

Please use:

- `.env` file to change configs of the app while working natively
- Check out `docker-compose.yaml` to override configs of the app when it's run as docker container

### Commands

To see the list of available commands run the following scripts from the home directory:

```sh
go run ./cmd/cli/cli.go
```

The CLI will prompt all the available commands and you can select one of them to be run, accepting input parameters. E.g.

```sh
go run ./cmd/cli/cli.go create-user --username pippo --password pluto
```

## 📄 License

This project is licensed under the [Apache 2.0 License](LICENSE).

## 🫶 Support Us

If you find this project useful, please consider supporting us.
