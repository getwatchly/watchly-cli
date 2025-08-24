# watchly-cli
CLI for connecting your workflows to Watchly

## Installation

Download `watchly-cli`:

### Linux
```bash
curl -L https://github.com/getwatchly/watchly-cli/releases/download/v0.0.5/watchly-cli_0.0.5_linux_386.tar.gz > watchly-cli.tar.gz
tar -xzf watchly-cli.tar.gz
cp watchly-cli /usr/local/bin
```

## Usage

Notify Watchly about a new deployment:
```bash
watchly-cli deployment start -k YOUR_API_KEY
```

Notify Watchly about a deployment's completion:
```bash
watchly-cli deployment finish -k YOUR_API_KEY -s successful
```