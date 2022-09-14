# Timed Queue (Cache) implementation
This service will store values provided via API up to certain time. If the value has been accessed, expiration time updates. Key-Value stores in binary file with [ttlStore](github.com/BON4/timedQ/tree/master/pkg/ttlstore) package.

## Install
```
> go get github.com/BON4/timedQ
```

## Run
```
> cd employees\cmd\app
timedQ\cmd\app> go build .
timedQ\cmd\app> .\app
```

## Test
```
timedQ> go test ./...
```


## Service architecture
![alt text](https://github.com/BON4/timedQ/blob/master/architecture.svg?raw=true)
