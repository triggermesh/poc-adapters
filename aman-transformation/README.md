# aman-transformation

## Running the App

- provide PORT enviroment variable

```cmd
 export PORT=8080
```

- Run the following command to run main.go

```cmd
go run .
```

or

```cmd
go run main.go
```

- Now fire up seprate terminal for client and execute the following

```cmd
curl -X POST localhost:<port-number>/<method name> \
   -H 'Content-Type: application/json' \
   -d '{"name":"bob","address":"my_password"}'
   ```