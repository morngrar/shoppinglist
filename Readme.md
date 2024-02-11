# shoppinglist

A shopping list app using HTMX, written to explore that technology.


## Setting up local dev environment

Start a docker container running postgres:

```command
./setup_db
```

Use the `psql` tool to run the migration scripts in order:

```command
$ psql -h localhost -U postgres 

postgres=# create database shoppinglist;

postgres=# \c shoppinglist;

postgres=# \ir migrations/0001.sql   -- will run first migration, repeat as needed
```

This tool may need to be installed on your machine, via a package. On arch this
is done with:

```command
sudo pacman -S postgresql-libs
```

## Deployment

The application is meant to be run on a linux system with an already prepared
external postgres instance. You will have to define the following environment
variables:

```
DB_USER
DB_PASSWORD
DB_HOST
DB_DATABASE
```

Or else the application will assume that it runs in a local dev environment and
use the credentials above.
