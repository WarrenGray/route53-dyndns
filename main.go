package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

const (
	resolver = "http://whatismyip.akamai.com/"
)

func fatal(s string, args ...interface{}) {
	fmt.Printf(s+"\n", args...)
	os.Exit(1)
}

func args() (zone, fqdn string, ttl int64) {
	if len(os.Args) < 4 {
		fatal("usage: %s <hosted zone ID> <fqdn> <ttl>\n", os.Args[0])
	}

	t, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fatal("invalid ttl: %s", err)
	}

	return os.Args[1], os.Args[2] + ".", int64(t)
}

func main() {
	zone, fqdn, ttl := args()
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fatal("unable to load AWS config: %s", err)
	}

	resp, err := http.Get(resolver)
	if err != nil {
		fatal("could not determine IP address: %s", err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fatal("could not determine IP address: %s", err)
	}

	ip := string(b)
	fmt.Printf("%s/%s => %s ttl=%d\n", zone, fqdn, ip, ttl)
	_, err = route53.NewFromConfig(cfg).ChangeResourceRecordSets(
		ctx,
		&route53.ChangeResourceRecordSetsInput{
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name:            aws.String(fqdn),
						Type:            types.RRTypeA,
						ResourceRecords: []types.ResourceRecord{{aws.String(ip)}},
						TTL:             aws.Int64(ttl),
					},
				}},
				Comment: aws.String("updated by route53-dyndns"),
			},
			HostedZoneId: aws.String(zone),
		},
	)

	if err != nil {
		fatal("record update failed: %s", err)
	}

	fmt.Println("done")
}
