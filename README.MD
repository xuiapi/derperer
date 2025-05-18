# Derperer

Derperer is mother of [derper](https://pkg.go.dev/tailscale.com/cmd/derper).

## Quick Start

Serve A derperer server, which auto collect DERP nodes , speedtest to them and automatically updates Tailscale Access Controls.

First, you need to have a [fofa](https://fofa.info/) account and a [Tailscale-API-access tokens](https://login.tailscale.com/admin/settings/keys)

And fill your credentials in the `derperer.yaml`
```yaml
config:
  # Fofa client configuration
  FofaClient:
    Email:  # fofa email
    Key:    # fofa key
  AdminToken:           # admin token
  Address: :8080               # Service listening address (default: ":8080")
  DataPath: /data               # Data storage path (default: "./data")
  Account:                     # Tailscale account
  ApiKey:                      # Tailscale apikey
  FetchBatch: 100                # Batch fetch count (default: 100)
  FetchInterval: 6h0m0s          # Fetch interval (default: 24h0m0s)
  UpdateInterval: 4h0m0s         # Tailscale ACL update interval
  DERPMapPolicy:
    BaselineBandwidth: 2         # Bandwidth limit, unit: Mbps (default: 2)
    CheckDuration: 10s           # Check duration (default: 5s)
    RecheckInterval: 1h0m0s      # Recheck interval (default: 1h0m0s)
    TestConcurrency: 1           # Test concurrency (default: 4)
```

Then start by

```bash
docker compose up -d
```

There are three main routes:

- `/derp.json` return a DERP map like https://controlplane.tailscale.com/derpmap/default
- `/swagger/index.html` swagger ui
- `/admin` Admin page, available when `--config.AdminToken` is set


## How it works

The derperer will search DERP nodes from [fofa](https://fofa.info/) by fingerprint `"<h1>DERP</h1>"`.
Then it will speedtest to each DERP node using official [derp](https://pkg.go.dev/tailscale.com/derp) binary.

The speedtest result will be saved in `data` directory, and the DERP map will be served at `/derp.json`.

The nodes are sorted by bandwidth, and Tailscale Access Controls are automatically updated.

## Reference
[https://github.com/yoshino-s/derperer](https://github.com/yoshino-s/derperer)