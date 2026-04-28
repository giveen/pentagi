# How to Connect to Hack The Box

## Connect

1. Download your Hack The Box OpenVPN configuration file from the Hack The Box VPN access page.
2. Place it somewhere local, for example:

```bash
/tmp/htb.ovpn
```

3. Start the VPN connection:

```bash
sudo openvpn --allow-compression yes --config /tmp/htb.ovpn
```

> **Note:** The `--allow-compression yes` flag is required because the HTB config includes
> `compress lzo`. OpenVPN 2.6+ blocks compression by default (VORACLE mitigation) and will
> refuse to connect without explicitly allowing it.

4. Keep that terminal open while you are working. A successful connection usually ends with output similar to:

```text
Initialization Sequence Completed
```

## Disconnect

If OpenVPN is running in the foreground, stop it with:

```bash
Ctrl+C
```

If you need to stop it from another terminal:

```bash
sudo pkill -f "openvpn --config /tmp/htb.ovpn"
```

## Verify Connection

To confirm the VPN tunnel exists:

```bash
ip addr show tun0
```

To check whether OpenVPN is still running:

```bash
pgrep -af openvpn
```