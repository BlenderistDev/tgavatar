# tgavatar
Go service for telegram avatar updating automation

## Description
Service load image for current hour and change it. Every hour telegram avatar updates with new image.

## Installation
````
git clone https://github.com/BlenderistDev/tgavatar
````
## Envariment variables
APP_ID - telegram app id. https://core.telegram.org/api/obtaining_api_id for more information
APP_HASH - telegram app hash. https://core.telegram.org/api/obtaining_api_id for more information
HOST - host for auth web server.
TIMEZONE - timezone for avatar generation.

## Launch
### Docker
#### Build
```
docker build -t tgavatar .
docker run -d -p 8081:8081 -v ~/storage:/app/storage --name tgavatar --restart always tgavatar
```

## Authorization
To launch service you need to auth in telegram. Here is a web service for auth. You need to enter phone and auth code.
Session is stored in storege/session. You need to make a volume for this folder if you are using docker for launching.

