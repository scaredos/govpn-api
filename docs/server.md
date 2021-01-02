## Web Management // Server Documentation

## Basic Usage

### To run the web management server on default port, run the following command
- Note: The default port is 1337

```./server run```

### To run the web management server on port `1234`, run the following command
- Note: In order for the web management to work, the REST-API is also automatically deployed on said port

```./server run --port 1234```

### To run **ONLY** the REST-API service on port `1234`, run the following command
- Note: This makes the web management unaccesiable

```./server runapi --port 1234```
