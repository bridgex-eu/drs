# DRS

Run your software anywhere.

- 🛰️ One-command deployment.
- 🗄️ Runs on any machine.
- 🎮 Full control at your fingertips.
- 🔒 Auto-secured with TLS certificate.
- 🗂️ Multiple apps, one machine.
- 🤩 Makes development enjoyable.

It's a new, minimal deployment tool built on SSH and Docker. You can run and share any software that can be placed into a Docker container on any machine: bare metal, AWS EC2, Hetzner Cloud VM, etc.

## Installation

See [releases](https://github.com/bridgex-eu/drs/releases) for pre-built binaries.

On Unix:

```
env CGO_ENABLED=0 go install -ldflags="-s -w" github.com/bridgex-eu/drs@latest
```

On Windows `cmd`:

```
set CGO_ENABLED=0
go install -ldflags="-s -w" github.com/bridgex-eu/drs@latest
```

On Windows powershell:

```
$env:CGO_ENABLED = '0'
go install -ldflags="-s -w" github.com/bridgex-eu/drs@latest
```

## How to use

```
drs deploy -m MACHINE_IP -p [YOUR_DOMAIN]:SERVER_PORT:CONTAINER_PORT --name APP_NAME DOCKER_IMAGE
```

> [!NOTE]
> - DRS by default uses ssh-agent to connect to the machine. Alternatively, you can use DRS Keys instead of ssh-agent. See the Key Management section below for more details.
> - Use the -r flag to pull image on a remote server (e.g. from Docker Registry) instead of transferring it via SSH.

This command will:

1. Install docker on remote machine
2. Run traefik
3. Create a docker network for drs apps
4. Send your app image through SSH or pull it from a registry (if -r flag is set)
5. Run your app with traefik labels

> [!TIP]
> To make your app accessible from the internet, you need to set up a DNS record for your domain pointing to the machine's IP address. DRS will automatically generate a TLS certificate for your domain.

### App Management

To manage your apps, use drs do. This command proxies Docker commands to the machine. Example:

```
drs do -m pluto ps # get a list of apps
drs do -m pluto myapp logs # check logs
drs do -m 125.23.43.110 myapp stop # stop an app
drs do -m 125.23.43.110 myapp rm # remove an app
```

### Machine Management

You can manage your machines with DRS:

```
drs machine add --name pluto --key pluto-rsa 123.43.10.120
```

Now you can access your machine by name in other commands: `drs do -m pluto ...`

### Key Management

DRS also assists you with your SSH keys:

```
drs key add --name pluto-rsa --passphrase passphrase ./pluto_rsa
```

To use this key to access your machines, use:

```
drs machine add --name pluto --key pluto-rsa 123.43.10.120
```
