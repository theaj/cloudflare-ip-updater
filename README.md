# Cloudflare IP Updater
A simple tool to detect WAN IP changes and updates your cloudflare DNS entries.

## Usage
1. Create a new [Cloudflare API token](https://dash.cloudflare.com/profile/api-tokens)
2. Copy the example.env file `cp example.env .env`
3. Populate the environment variables in the `.env` file. The `ZONE_NAME` is the name of the Cloudflare zone / domain name and the `DNS_RECORD` is the name of the record you would like to change.
4. Start the docker container as a daemon: `./daemon.sh start`

The docker container is configured to restart if your computer also restarts (assuming docker starts on boot).