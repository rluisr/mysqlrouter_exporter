If you want to run InnoDB Cluster with http and https REST API, run this command.

```shell
$ git clone git@github.com:rluisr/mysqlrouter-go.git
$ cd mysqlrouter-go/test
$ docker compose up --build --force-recreate --always-recreate-deps --renew-anon-volumes -d;
```

### Basic

root / mysql