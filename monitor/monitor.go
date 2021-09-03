package monitor

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog/log"
)

const WhatIsMyIPAddressURL = "https://ifconfig.me/ip"

func Start() {
	ctx := context.Background()

	cfKey := os.Getenv("CF_API_KEY")
	if cfKey == "" {
		log.Fatal().Msg("CF_API_KEY cannot be empty.")
	}

	api, err := cloudflare.NewWithAPIToken(cfKey)
	if err != nil {
		log.Fatal().AnErr("error", err).Msgf("Failed to start monitoring")
	}

	zoneName := os.Getenv("ZONE_NAME")
	if zoneName == "" {
		log.Fatal().Msg("ZONE_NAME cannot be empty.")
	}

	dnsRecordName := os.Getenv("DNS_RECORD")
	if dnsRecordName == "" {
		log.Fatal().Msg("DNS_RECORD cannot be empty.")
	}

	checkInterval := 60 * time.Second
	if interval, exists := os.LookupEnv("CHECK_INTERVAL"); exists {
		i, err := strconv.ParseInt(interval, 10, 16)
		if err != nil {
			log.Fatal().AnErr("error", err).Msgf("Could not parse environment variable: CHECK_INTERVAL")
		}
		checkInterval = time.Duration(i) * time.Second
	}

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
		time.Sleep(checkInterval)
	}
}
