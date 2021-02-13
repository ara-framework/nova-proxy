# Nova Proxy
Hypernova Proxy is an Reverse Proxy whick look in the hosts responses for [Hypernova Directives](https://github.com/marconi1992/hypernova-blade-directive) in order to inject the components rendered by [Hypernova Server](https://github.com/airbnb/hypernova).

## Environment Variables

```env
HYPERNOVA_BATCH=http://hypernova:3000/batch
CONFIG_FILE=config.json
```

## Configuration File

Nova Proxy needs a configuration file:

```json
//nova-proxy.json

{
  "locations": [
    {
      "path": "/",
      "host": "http://blog:8000",
      "modifyResponse": true,
      "proxyPreserveHost": true
    }
  ]
}
```

The `locations` items require the `path` and `host` to let know to Nova Proxy which the application is responsible to serve the requested page. By default the path `/` passes all the requests to the declared host.

The `modifyResponse` enable the serve-side includes to that location.

The `proxyPreserveHost` is used to preserve and retain the original Host: header from the client browser when constructing the proxied request to send to the target server.

## Using Nova Proxy with [Ara CLI](https://github.com/ara-framework/ara-cli)

Before to run the command we need to set the `HYPERNOVA_BATCH` variable using the Nova service endpoint.

```shell
export HYPERNOVA_BATCH=http://localhost:3000/batch
```

The command uses a configuration file named `nova-proxy.json` in the folder where the command is running, otherwise you need to pass the `--config` parameter with a different path.
```
ara run:proxy --config ./nova-proxy.json
