## Swapp-Static

Simple Go API for handling and serving static files


### .env

```
JWTSECRET=xxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Usage

All "POST" routes are authenticated, the "GET"s can be anonymous

`/api/upload` - Multipart form with ("Images" and "dir") keys. "dir" is the directory to which you're uploading the files.

#### Example response

```js
{
  "Message": "2 out of 2 files were uploaded",
  "FileUrls": [
    "/static/products/upload-585578885.png",
    "/static/products/upload-152937248.png"
  ]
}
```

