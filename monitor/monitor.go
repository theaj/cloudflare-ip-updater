package monitor

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog/log"
)

const CheckInterval = 5 * time.Second
const WhatIsMyIPAddressURL = "https://ifconfig.me/ip"

func Start() {
	ctx := context.Background()
	api, err := cloudflare.NewWithAPIToken(os.Getenv("CF_API_KEY"))
	if err != nil {
		log.Fatal().AnErr("error", err).Msgf("Failed to start monitoring")
	}

	zoneName := os.Getenv("ZONE_NAME")
	dnsRecordName := os.Getenv("DNS_RECORD")

	log.Info().Msgf("Retrieving current IP for DNS entry...")

	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("Could not get Zone ID")
	}

	log.Info().Msgf("Zone ID: %s", zoneID)
	records, err := api.DNSRecords(ctx, zoneID, cloudflare.DNSRecord{Type: "A", Name: dnsRecordName})
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("Could not get DNS records")
	}

	if len(records) > 1 {
		log.Fatal().Msg("More than one DNS entry found. Do not know which one to use")
	}

	dnsRecord := records[0]
	currentIP := dnsRecord.Content

	log.Info().Msgf("DNS IP found: %s", currentIP)

	for {
		log.Info().Msgf("Checking IP...")
		resp, err := http.Get(WhatIsMyIPAddressURL)
		if err != nil {
			log.Err(err).Msg("Could not get IP address")
		} else {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Err(err).Msg("Could not read HTTP body")
			} else {
				wanIP := string(data)
				if wanIP != currentIP {
					log.Info().Msgf("Updating DNS record with new IP: %s", wanIP)
					err = api.UpdateDNSRecord(ctx, zoneID, dnsRecord.ID, cloudflare.DNSRecord{Type: "A", Name: dnsRecordName, Content: wanIP})
					if err != nil {
						log.Err(err).Msg("Could not update DNS record")
					} else {
						currentIP = wanIP
					}
				}
			}
		}
		time.Sleep(CheckInterval)
	}
}
