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
drs deploy -m MACHINE_IP -p [YOUR_DOMAIN]:SERVER_PORT:CONTAINER_PORT --name APP_NAME LOCAL_DOCKER_IMAGE
```

> [!TIP]
> Use -r to pull your image on a remote server instead of sending it through SSH.

This will:

1. Install docker on remote machine
2. Run traefik
3. Create a docker network for drs apps
4. Send your app image through SSH
5. Run your app with traefik labels

> [!NOTE]
> To make your app available by your domain, you need to set A DNS record to IP of the machine.

### App management

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
drs machine add --name pluto --key-file ./pluto_rsa 123.43.10.120
```

Now you can access your machine by name in other commands: `drs do -m pluto ...`

### Key Management

DRS also assists you with your SSH keys:

```
drs key add --name pluto-key --passphrase passphrase ./pluto_rsa
```

To use this key to access your machines, use:

```
drs machine add --name pluto --key pluto-rsa 123.43.10.120
```
