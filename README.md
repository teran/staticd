# staticd

[![Layers size](https://images.microbadger.com/badges/image/teran/staticd.svg)](https://hub.docker.com/r/teran/staticd/)
![Recent build commit](https://images.microbadger.com/badges/commit/teran/staticd.svg)
[![Docker Automated build](https://img.shields.io/docker/automated/teran/staticd.svg)](https://hub.docker.com/r/teran/staticd/)
![License](https://img.shields.io/github/license/teran/staticd.svg)

Web server for static using Amazon S3 compatible API as a backend

## Usage env vars

 * `STATICD_ALLOWDELETE` - default:`false`, whether we should handle DELETE requests
 * `STATICD_ALLOWGET` - default:`true`, whether we should handle GET requests
 * `STATICD_ALLOWPUT` - default:`false`, whether we should handle PUT requests
 * `STATICD_DEBUG` - Set log verbosing to debug(not implemented yet)
 * `STATICD_LISTEN` - default:`:8080`, address to listen on
 * `STATICD_MAXUPLOADSIZE` - default:`1024`, size in MBytes with max allowed file size to upload
 * `STATICD_S3ACCESSKEY` - S3 access key
 * `STATICD_S3BUCKETNAME` - S3 bucket name
 * `STATICD_S3ENDPOINT` - S3 endpoint URI
 * `STATICD_S3MODE` - `<proxy|redirect>`, default:`proxy` - file GET method, proxy uses simple proxying, redirect generates presigned URL and sends client to it via HTTP redirect
 * `STATICD_S3SECRETKEY` - S3 secret key
 * `STATICD_S3REGION` - S3 region
 * `STATICD_S3REDIRECTURLTTL` - default:`1800s`, TTL for presigned URL, i.e. how long it would stay valid
 * `STATICD_S3USESSL` - `<true|false>`, default:`true` - should we use SSL for backend or not
