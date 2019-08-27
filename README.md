# Hypernova Proxy
Hypernova Proxy is an Reverse Proxy whick look in the hosts responses for [Hypernova Directives](https://github.com/marconi1992/hypernova-blade-directive) in order to inject the components rendered by [Hypernova Server](https://github.com/airbnb/hypernova).

## Environment Variables

```env
HYPERNOVA_BATCH=http://hypernova:3000/batch
CONFIG_FILE=config.json
```

## Configuration File

Hypernova Proxy needs a configuration file in order to setup the location for the reverse proxy, you can also specify which location needs post-process with the Hypernova views.

```json
{
  "locations": [
    {
      "path": "/",
      "host": "http://blog:8000",
      "modifyResponse": true
    },
    {
      "path": "/public/client.js",
      "host": "http://hypernova:3000"
    }
  ]
}
```

## Using Hypernova Proxy with Docker

```Dockerfile
FROM araframework/nova-proxy:1.0.5

COPY config.json config.json
```