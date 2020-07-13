# Quid

A Json web tokens server

## Install

Clone the repository

### Database

Create the Postgresql database:

   ```bash
   sudo su postgres
   psql
   create database quid;
   \c quid;
   GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public to pguser;
   \q
   exit
   ```

Replace `pguser` by your database user

### Configuration

Create a config file:

   ```bash
   go run main.go -conf
   ```

Edit the config file to provide your database credentials. Initialize the database and create an admin user:

   ```bash
   go run main.go -init
   ```

## Run

Run the backend:

   ```bash
   go run main.go
   ```

Run the frontend:

   ```bash
   cd quidui
   npm install
   npm run serve
   ```

Go to `localhost:8080` to login into the admin interface

![Screenshot](doc/img/screenshot.png)