# Aviasales Autocomplete Stub

## Main purpose

The following service acts as a cache for aviasales autocomplete. It does:

- Proxies all search queries right to [aviasales autocomplete](https://places.aviasales.ru/v2/places.json)
  example: 
  ```json
    {
        "slug": "MOW",
        "title": "Moscow",
        "subtitle": "Russia"
      }
   ```
- Converts responses to the format required by web widgets
- Stores converted responses in a local bolt DB storage to avoid delays on service-to-service communication 

## How to run service

### Docker compose

The easiest way to run service is docker-compose:

```bash
git clone https://github.com/titusjaka/autocomplete
cd autocomplete
docker-compose up
```

To configure service, just change environment variables in [docker-compose.yml](docker-compose.yml). 

```yaml
    environment:
      - AVS_STUB_LISTEN=:8080 # host:port to run service on
      - AVS_STUB_VERBOSE=true # remove to disable debug output
      - AVS_STUB_DBPATH=autocomplete.db # bolt DB file name
      - AVS_AUTOCOMPLETE_URL=https://places.aviasales.ru # URL for aviasales autocomplete service
      - AVS_AUTOCOMPLETE_TIMEOUT=3s # timeout for connection to aviasales autocomplete service
```

### Build and run

You can also build service by yourself out of source code. Go [1.11+](https://golang.org/dl/) required.

```bash
git clone https://github.com/titusjaka/autocomplete
cd autocomplete
make vendor
make build
./build/autocomplete
```

You can change default configuration with following arguments:
```
Usage:
  autocomplete [OPTIONS]

Application Options:
      --http.listen=              HTTP Listen address (default: :8080) [$AVS_STUB_LISTEN]
      --verbose                   Verbose output [$AVS_STUB_VERBOSE]
      --db.path=                  Path to bolt DB file (default: autocomplete.db) [$AVS_STUB_DBPATH]
      --avs.autocomplete.url=     URL of Aviasales autocomplete service (default: https://places.aviasales.ru) [$AVS_AUTOCOMPLETE_URL]
      --avs.autocomplere.timeout= Timeout for connection to Aviasales autocomplete service (default: 3s) [$AVS_AUTOCOMPLETE_TIMEOUT]

Help Options:
  -h, --help                      Show this help message
```

## Usage

Just open the following URL in your browser or [Postman](https://www.getpostman.com/):

[http://127.0.0.1:8080/search?term=mow&locale=us&type%5B%5D%3Dcity%26type%5B%5D%3Dairport](http://127.0.0.1:8080/search?term=mow&locale=us&type[]=city&type[]=airport)

All [original](https://places.aviasales.ru/v2/places.json) query parameters are supported 😊 (you know, it only proxies requests).
