## Swapp-Static

Simple Go API for handling and serving static image files


### .env

```
JWTSECRET=xxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Usage

All "POST" routes are authenticated, the "GET"s can be anonymous

#### Create directory/folder


`/api/create-dir`

```js
// Example Body 

{
	"DirName": "products"
}
```

#### Upload file(s)

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

#### Get or view the file

`/static/{folder}/{filename}` - example (/static/products/upload-152937248.png). The server return a file or a 404 if the image isn't found


#### Cleanup 
`/api/purge` - removes a directory and it's children

```
{
	"FileDir": "profiles"
}
```


