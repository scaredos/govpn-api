## API Documenation


## Initialize
- This is required to use any of the other endpoints
- Endpoint (GET /init)
- HTTP Parameters
```
hubuser: API Username/Administrator Password (Commonly VPN)
hubpass: API Password/Administrator Password
```

## Create User
- Endpoint (GET /createUser)
- HTTP Parameters
```
username: User username
password: User password (base64 encoded)
sip     : SoftEther Server IP (E.x.: 127.0.0.1:443)
```

## Delete User
- Endpoint (GET /deleteUser)
- HTTP Parameters
```
username: User username
sip     : SoftEther Server IP (E.x.: 127.0.0.1:443)
```

## Change Password of User
- Endpoint (GET /changePassword)
- HTTP Parameters
```
username: User username
password: New user password (base64 encoded)
serverip: SoftEther Server IP (E.x.: 127.0.0.1:443)
```

## Set Expiration Date
- Endpoint (GET /setExpireDate)
- HTTP Parameters
```
username: User username
expdate : Expiration Date (date object) (2014-10-12T00:00:00.000Z)
sip     : SoftEther Server IP (E.x.: 127.0.0.1:443)
```


## Get User Info
- Endpoint (GET /getUser)
- HTTP Parameters
```
username: User username
sip     : SoftEther Server IP (E.x.: 127.0.0.1:443)
```
