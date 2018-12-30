Funker
=======
> Funker is a platform for function as a service based javascript.

Features
========
- Standalone.
- Using Redis/Redix for storing `funktions`.
- Includes multiple modules like `crypto`, `redis`, `base64`, `uniqid`, ... etc.
- Internal `javascript` runtime pool to boost the performance.

Available Modules
==================
- `fetch`
- `redis`
- `crypto`
- `uniqid`
- `base64`
- `validator`

API Docs
========
> You can view it [here](https://documenter.getpostman.com/view/2408647/RzfZQDJF)

Installation
=============
> `$ go get github.com/alash3al/funker`

Quick Overview
===============
> a Funk should be in the following style
```js
function(){
    // your code goes here
    // `this` includes many useful helpers like the following
    
    // this.request: 
    {
        "uri":         String,
        "proto":       String,
        "method":      String,
        "path":        String,
        "host":        String,
        "https":       Boolean,
        "query":       Object,
        "body":        Object,
        "remote_addr": String,
        "real_ip":     String,
        "headers":     Object,
    }

    // this.response
    Response this.response.send(Any output)
    Response this.response.status(Integer statusCode)
    Response this.response.headers(Object headers)
    Response this.response.type(String) // supported types: ['json', 'html']
    
    // this.module
    // load a module, supported modules are: ['fetch', 'crypto', 'localStorage', 'uniqid', 'base64']
    // example:
    uniqid = this.module("uniqid")
    this.response.status(200).type("json").send({theUniqueID: uniqid(50)})

    // or
    this.response.send(this.request)
}
```

Modules RFCs
============
## # fetch
> a broser like `fetch` function that implements a `http` client
```js
options = {
    method: String,
    headers: Object,
    body: Any,
    redirects: Integer,
    timeout: Integer,
    proxy: String
}
responseObject = fetch(String url, [Object options])
/**

    {
        code: Integer,
        headers: Object,
        body: Any
    }

*/
```

## # Redis
> uses [RedisClient](https://godoc.org/github.com/go-redis/redis#Client), just go there and read the docs


## # crypto
> a crypto module for hashing & encrypting data
```js

String crypto.md5(String)
String crypto.sha1(String)
String crypto.sha256(String)
String crypto.sha512(String)

String crypto.bcrypt.hash(String password)
String crypto.bcrypt.check(String hashedString, String password)

String crypto.aes.encrypt(String data, String a32bitString)
String crypto.aes.decypt(String encryptedData, String a32bitString)

```

## # uniqid
> generate a unique random string with a custom length
```js
String uniqid([Integer length = 15])
```

## # base64
> encode/decode to and from base64
```js
String base64.encode(String)
String base64.decode(String)
```

## # validator
> validates an input against some rules, it uses [govalidator](https://github.com/asaskevich/govalidator).
```js
Object validator.validate(Object data, Object rules)

// Example
data = {
    "name": "",
    "email": "this@is.mail"
}
rules = {
    "name": ["required", "stringlength:5,10"],
    "email": ["email"]
}

result = this.module("validator").validate(data, rules)

// errors count
result.errors

// errors messages
result.fields
```
