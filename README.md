# IP2Location Service (using maxmind database)

## How to Use

1. pull docker image

```
docker pull moinologics/ip2location:latest
```

2. run docker container

```
docker run --rm \
  -e MAXMIND_ACCOUNT_ID=your-maxmind-user-id \
  -e MAXMIND_LICENSE_KEY=your-maxmind-license-key \
  -e MAXMIND_EDITION_IDS=GeoLite2-City \
  -p 8080:8080 \
  moinologics/ip2location:latest
```

3. get location for ipv4

```
curl --location 'http://localhost:8080/location/38.152.13.107
```

4. get location for ipv6

```
curl --location 'http://localhost:8080/location/2001:4860:4860::8888'
```

## Docker Container Envirolment Variables

| ENV                 | DESCRIPTION                                                                                                                        | Required |
| ------------------- | ---------------------------------------------------------------------------------------------------------------------------------- | -------- |
| MAXMIND_ACCOUNT_ID  | your maxmind account id<br />(can we found on [license key page](https://www.maxmind.com/en/accounts/current/license-key))           | YES      |
| MAXMIND_LICENSE_KEY | your maxmind license key<br />(can we found/generate on [license key page](https://www.maxmind.com/en/accounts/current/license-key)) | YES      |
| MAXMIND_EDITION_IDS | space seprated GEOIP EDITION IDs<br />for example **GeoLite2-City**                                                         | YES      |
| ALLOWED_API_KEY     | API key which should we provided in API-KEY header for Authentication<br />if not provided, no authentication will be used         | NO       |

## HTTP request Authentication

if you provided ALLOWED_API_KEY env in docker container with a non empty string, than you have to send this api key in header API-KEY with same value

example request

```
curl --location 'http://localhost:8080/location/38.152.13.107' --header 'API-KEY: yourapikey'
```

## Automatic database update

one cronjob is also scheduled for running every night at 12:00 (UTC)

# Developer Contact

[mohin.ahmad.dev@gmail.com](mailto:mohin.ahmad.dev@gmail.com)
