## Web Management // Server Documentation

## Basic Usage

### To run the web management server on default port, run the following command
- Note: The default port is 1337

```./server run```

### To run the web management server on port `1234`, run the following command
- Note: In order for the web management to work, the RestAPI is also deployed

```./server run --port 1234```

### To run **ONLY** the REST-API service on port `1234`, run the following command
- Note: The web interface will no longer be avaialble, only the RestAPI 

```./server runapi --port 1234```
