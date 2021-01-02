## API Documenation


## Initialize
- This is required to use any of the other endpoints
- Endpoint (/init)
- Parameters
```
hubuser: API Username/Administrator Password (Commonly VPN)
hubpass: API Password/Administrator Password
GET /init?hubuser=hubuser&hubpass=hubpass
```

## Create User
- Endpoint (/createUser)
- Parameters
```
username: User username
password: User password (base64 encoded)
serverip: SoftEther Server IP (E.x.: 127.0.0.1:443)
```

## Delete User
- Endpoint (/deleteUser)
- Parameters
```
username: User username
serverip: SoftEther Server IP (E.x.: 127.0.0.1:443)
```

## Change Password of User
- Endpoint (/changePassword)
- Parameters
```
username: User username
password: New user password (base64 encoded)
serverip: SoftEther Server IP (E.x.: 127.0.0.1:443)
```

## Get User Info
- Endpoint (/getUser)
- Parameters
```
username: User username
serverip: SoftEther Server IP (E.x.: 127.0.0.1:443)
```
