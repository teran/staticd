# staticd

Web server for static using Amazon S3 compatible API as a backend

## Usage

 * `STATICD_DEBUG` - Set log verbosing to debug(not implemented yet)
 * `STATICD_LISTEN` - address to listen on
 * `STATICD_S3ACCESSKEY` - S3 access key
 * `STATICD_S3BUCKETNAME` - S3 bucket name
 * `STATICD_S3ENDPOINT` - S3 endpoint URI
 * `STATICD_S3MODE` - `<proxy|redirect>` - file GET method, proxy uses simple proxying, redirect generates presigned URL and sends client to it via HTTP redirect
 * `STATICD_S3SECRETKEY` - S3 secret key
 * `STATICD_S3REGION` - S3 region
 * `STATICD_S3REDIRECTURLTTL` - `default:1800s` - TTL for presigned URL, i.e. how long it would stay valid
 * `STATICD_S3USESSL` - `<true|false>` - should we use SSL for backend or not
