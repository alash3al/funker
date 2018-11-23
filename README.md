Funker
=======
> Funker is a platform for function as a service based javascript.

Features
========
- Standalone
- Using Redis for storing `funktions`
- Implements a caching mechanism based on redis
- Includes `lodash` javascript library
- Includes multiple modules like `crypto`, `localStorage`, `aes`, ... etc
- Internal `javascript` runtime pool

Available Modules
==================
- `exports`
- `fetch`
- `localStorage`
- `crypto`
- `uniqid`
- `base64`
- `env`
- `lodash`

API Docs
========
> You can view it [here](https://documenter.getpostman.com/view/2408647/RzfZQDJF)

Installation
=============
> `$ go get github.com/alash3al/funker`

Modules RFCs
============

## # exports
> exports is a global variable represents what we want to send to the client i.e `browser`, its all about a response object contains the `type` of the response `DocType` (`json` or `html`), the response status code `Status`, the response `Headers` and the `Body`.
```js
exports = {
    DocType: "json", // json|html
    Status: 200, // status code
    Headers: {
        "Key": "Value"
    },
    Body: {
        message: "Hello World"
    }
}
```

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

## # localStorage
> emulates the browser's `localStorage` but the backend is `redis`, it has the following methods

```js

Void localStorage.set(String namespace, String key, Any data)
Any localStorage.get(String namespace, String key)
localStorage.delete(String namespace, String key)
localStorage.delete(String namespace)
localStorage.getAll(String namespace)
localStorage.incr(String namespace, String key, Integer factor)

```

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

## # env
> request environment
```js
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
```