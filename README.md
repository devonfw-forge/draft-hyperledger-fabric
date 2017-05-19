## Chaincode Interface

This document outlines the interface for communicating with the PictureLicenseVerifier Chaincode.

### Init Function:
Request
```
{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "type": 1,
    "chaincodeID": {
      "path": "plv"
    },
    "ctorMsg": {
      "function": "Init"
    },
    "secureContext": "WebAppAdmin"
  },
  "id": 0
}
```
Response
```
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "1f8dde6aca14d49e2281346019227c3971d2bc0aa4217f576c270af140398ba78edd268ddf3bdf5f0c884f000490304b633a23d5e2463c76b64d5e59c4071862"
  },
  "id": 0
}
```
### Invoke Functions: 
#### Add user: 
Request
```
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": " 1f8dde6aca14d49e2281346019227c3971d2bc0aa4217f576c270af140398ba78edd268ddf3bdf5f0c884f000490304b633a23d5e2463c76b64d5e59c4071862"
    },
    "ctorMsg": {
      "function": "addUser",
      "args": ["username@capgemini.com","{\"password\":\"123456\", \"participant-type\":\"employee\"}"]
    },
    "secureContext": "WebAppAdmin"
  },
  "id": 1
}
```
Response
```
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "4f4c0db4-319d-42bc-a744-3dc46a533615"
  },
  "id": 1
}
```
#### Demand image: 
Request
```
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "1f8dde6aca14d49e2281346019227c3971d2bc0aa4217f576c270af140398ba78edd268ddf3bdf5f0c884f000490304b633a23d5e2463c76b64d5e59c4071862"
    },
    "ctorMsg": {
      "function": "DemandImage",
      "args": ["{\"id\":\"IMG1\", \"name\":\"UNDEFINED\", \"author\" : \"ildogesto\", \"url\":\"http://www.istockphoto.com/vector/flat-design-icons-for-business-and-finance-gm509786662-85956153\", \"user\": \"username@capgemini.com\", \"md5-hash\" : \"UNDEFINED\", \"remarks\": \"UNDEFINED\", \"purchase-date\" : \"UNDEFINED\", \"status\":1}"]
    },
    "secureContext": "WebAppAdmin"
  },
  "id": 2
}
```
Response
```
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "5b9177c4-778a-4cbd-a35a-2a92d16ff02b"
  },
  "id": 2
}
```
#### Deliver image:
Request
```
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "1f8dde6aca14d49e2281346019227c3971d2bc0aa4217f576c270af140398ba78edd268ddf3bdf5f0c884f000490304b633a23d5e2463c76b64d5e59c4071862"
    },
    "ctorMsg": {
      "function": "DeliverImage",
      "args": ["IMG1","search-icon.png","da39a3ee5e6b4b0d3255bfef95601890afd80709","19.05.2017"]
    },
    "secureContext": "WebAppAdmin"
  },
  "id": 3
}
```

Response

```
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "3358adf3-fa91-4bdd-9ad2-9a4a823423dd"
  },
  "id": 3
}
```

### Query Functions: 
#### Authenticate as user:
Request
```
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "1f8dde6aca14d49e2281346019227c3971d2bc0aa4217f576c270af140398ba78edd268ddf3bdf5f0c884f000490304b633a23d5e2463c76b64d5e59c4071862"
    },
    "ctorMsg": {
      "function": "AuthenticateAsUser",
      "args": ["username@capgemini.com","123456"]
    },
    "secureContext": "WebAppAdmin"
  },
  "id": 4
}
```

Response (successful Authentication)
```
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "{\"User\":{\"password\":\"123456\",\"participant-type\":\"employee\"},\"Authenticated\":true}"
  },
  "id": 4
}
```

Response (failed Authentication)
```
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "{\"User\":{\"password\":\"\",\"participant-type\":\"\"},\"Authenticated\":false}"
  },
  "id": 4
}
```

#### Get users list: 
Request

```
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "1f8dde6aca14d49e2281346019227c3971d2bc0aa4217f576c270af140398ba78edd268ddf3bdf5f0c884f000490304b633a23d5e2463c76b64d5e59c4071862"
    },
    "ctorMsg": {
      "function": "getUsers"
    },
    "secureContext": "WebAppAdmin"
  },
  "id": 5
}
```
Response
```
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "{\"users\":[{\"username\":\"username@capgemini.com\",\"password\":\"123456\",\"participant-type\":\"employee\"},{\"username\":\"username2@capgemini.com\",\"password\":\"123456\",\"participant-type\":\"employee\"},{\"username\":\"username3@capgemini.com\",\"password\":\"123456\",\"participant-type\":\"employee\"}]}"
  },
  "id": 5
}
```

#### Get image by id: 
Request
```
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "1f8dde6aca14d49e2281346019227c3971d2bc0aa4217f576c270af140398ba78edd268ddf3bdf5f0c884f000490304b633a23d5e2463c76b64d5e59c4071862"
    },
    "ctorMsg": {
      "function": "getImage",
       "args": ["IMG1"]
    },
    "secureContext": "WebAppAdmin"
  },
  "id": 6
}
```
Response
```
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "{\"id\":\"IMG1\",\"name\":\"UNDEFINED\",\"author\":\"ildogesto\",\"url\":\"http://www.istockphoto.com/vector/flat-design-icons-for-business-and-finance-gm509786662-85956153\",\"user\":\"username@capgemini.com\",\"md5-hash\":\"UNDEFINED\",\"remarks\":\"UNDEFINED\",\"purchase-date\":\"UNDEFINED\",\"status\":2}"
  },
  "id": 6
}
```
#### Get images by user
Request
```
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "1f8dde6aca14d49e2281346019227c3971d2bc0aa4217f576c270af140398ba78edd268ddf3bdf5f0c884f000490304b633a23d5e2463c76b64d5e59c4071862"
    },
    "ctorMsg": {
      "function": "GetImagesByUser",
       "args": ["username@capgemini.com"]
    },
    "secureContext": "WebAppAdmin"
  },
  "id": 7
}
```
Response
```
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "{\"images\":[{\"id\":\"IMG1\",\"name\":\"UNDEFINED\",\"author\":\"ildogesto\",\"url\":\"http://www.istockphoto.com/vector/flat-design-icons-for-business-and-finance-gm509786662-85956153\",\"user\":\"username@capgemini.com\",\"md5-hash\":\"UNDEFINED\",\"remarks\":\"UNDEFINED\",\"purchase-date\":\"UNDEFINED\",\"status\":2}]}"
  },
  "id": 7
}
```

#### Get all images 
Request
```
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "1f8dde6aca14d49e2281346019227c3971d2bc0aa4217f576c270af140398ba78edd268ddf3bdf5f0c884f000490304b633a23d5e2463c76b64d5e59c4071862"
    },
    "ctorMsg": {
      "function": "GetImages"
    },
    "secureContext": "WebAppAdmin"
  },
  "id": 8
}
```
Response
```
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "{\"images\":[{\"id\":\"IMG1\",\"name\":\"IMG1.jpge\",\"author\":\"ildogesto\",\"url\":\"http://www.istockphoto.com/vector/flat-design-icons-for-business-and-finance-gm509786662-85956153\",\"user\":\"username@capgemini.com\",\"md5-hash\":\"da39a3ee5e6b4b0d3255bfef95601890afd80709\",\"remarks\":\"UNDEFINED\",\"purchase-date\":\"19.05.2017\",\"status\":2},{\"id\":\"IMG2\",\"name\":\"UNDEFINED\",\"author\":\"erhui1979\",\"url\":\"http://www.istockphoto.com/vector/teamwork-gm517994151-49374946\",\"user\":\"username2@capgemini.com\",\"md5-hash\":\"UNDEFINED\",\"remarks\":\"UNDEFINED\",\"purchase-date\":\"UNDEFINED\",\"status\":1}]}"
  },
  "id": 8
}
```
