# Gator

Gator is a boot.dev guided project. It is a CLI tool for rss feed aggregation. Gator is Operating System Agnostic.

#### Requirements

##### Golang Toolchain 1.23+

https://go.dev/dl/

##### PostgresSQL 17+

[https://www.postgresql.org/download/](https://www.postgresql.org/download/)

**Download and install both for your Operating System**

**Clone this repository to your machine you can find inactions on how to do this here.**

[https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository)

### **Setup**

#### **PostgreSQL**

Use psql command line tool to connect to Postgres (in windows go to start and search for sql shell (psql))

Create database called gator

Command Line Inactions

CREATE DATABASE gator;

Connect to gator.

\c gator

Set DB user password 'postgres'

ALTER USER postgres PASSWORD 'postgres';

If you suck with defaults this, is your database connection string.

postgres://postgres:postgres@localhost:5432/gator?sslmode=disable

if you did not use defaults below is a string that you need to modify for your use case.

postgres://%username%:%password%@%ipaddress%:%port%/%databasename%?sslmode=disable

Open the shell you are using; navigate to the repository that you have cloned.

We need to configure the database, don’t worry don’t need to do this manually.

Install goose its part of the go tool chain

Enter the fallowing into shell “go install github.com/pressly/goose/v3/cmd/goose@latest“

Once installed navigate to %repo%\sql\schema

Run this command

goose postgres "%your database connection string" up

goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up

go back to root of the repository open and look at  .gatorconfig.json modify the {"db_url" value if necessary. Copy .gatorconfig.json to home directly for windows its %userprofile%.

Now simply type “go install”

### Usage

1.      First register yourself as a user.
>		gator register %name%

2.      Login as the user you just registered
>		gator login %name%

3.      Add some RSS feeds.
>		gator addfeed “name” “feed url”

4.      Aggregate the feeds: note that the feeds will be continually polled unit stopped “CTRL+C”. Polling interval valid arguments are 1s, 1m, 1h.
> 		gator agg 1s

5.      Browse your feed browse has an optional argument of number of feeds you’d like to see.
> 		gator browse 5

Remember if you want to clear the database use the “gator reset”, there is no going back after reset. The data in data base will be erased.


### **Full list of commands**

##### **register**

Registers a new user, **gator register user**

##### **login**

Login as registered user, **gator login user**

##### **Addfeed**

Adds a new feed and subscribes the current user to it, **gator  “feed name” “feed url”**

##### **feeds**

Shows a list of available feeds, **gator feeds**

##### **follow**

Subscribes current user to a registered feed, **gator follow “url”**

##### **following**

Get a list of feeds you are subscribing to, **gator following**

##### **users**

Get a list of users, **gator users**

##### **agg**

Aggregate feeds must be done before browsing, **gator agg 1s, or 1m, or 1h**

##### **unfollow**

Unsubscribe from a feed, **gator unfollow “url”**

##### **browse**

return a list of most recently updated articles, **gator browse #***(number of articles you want to see)*

##### **reset**

wipe the database, **gator reset**


