## To Run the App

- provide a Path for environment variable

``` 
 export PORT=<port no.>
 ```

- Run the following command to run main.go
```
go run .
```
or
```
go run main.go
```

- Now fire up seprate terminal for client and execute the following
```
curl -X POST localhost:<port-number>/<method name> \
   -H 'Content-Type: application/json' \
   -d '{"name":"bob","address":"my_password"}'
   ```