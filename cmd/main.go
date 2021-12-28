package main

import (
	"fmt"
	"os"

	"github.com/hupe1980/fakedns"
	"github.com/spf13/cobra"
)

const (
	version    = "dev"
	defaultTTL = 60
)

func main() {
	var opts struct {
		addr            string
		net             string
		upstream        string
		ttl             uint32
		ipsV4           []string
		rebindV4        string
		ipsV6           []string
		rebindV6        string
		rebindThreshold int
	}

	rootCmd := &cobra.Command{
		Use:     "fakedns [domains]",
		Version: version,
		Short:   "Tiny DNS proxy for Penetration Testers and Malware Analysts",
		Args:    cobra.MaximumNArgs(1),
		Example: `IPV4: fakedns example.org --ipv4 127.0.0.1
Wildcards: fakedns example.* --ipv4 127.0.0.1
RoundRobin: fakedns example.org --ipv4 127.0.0.1,10.10.10.10
Rebind: fakedns example.org --ipv4 127.0.0.1 --rebind-v4 10.10.10
Upstream: fakedns example.org --ipv4 127.0.0.1 --upstream 8.8.8.8`,
		RunE: func(cmd *cobra.Command, args []string) error {
			options := &fakedns.Options{
				TTL:     opts.ttl,
				Domains: args,
				IPsV4:   opts.ipsV4,
				IPsV6:   opts.ipsV6,
				Rebind:  fakedns.NewRebind(opts.rebindV4, opts.rebindV6, opts.rebindThreshold),
			}

			if opts.upstream != "" {
				options.FallbackDNSResolver = opts.upstream
			}

			fakeDNS, err := fakedns.New(options)
			if err != nil {
				return err
			}

			return fakeDNS.ListenAndServe(opts.addr, opts.net)
		},
	}

	rootCmd.Flags().StringVarP(&opts.addr, "addr", "a", "0.0.0.0:53", "fakeDNS address")
	rootCmd.Flags().StringVarP(&opts.net, "net", "n", "udp", "fakeDNS network protocol")
	rootCmd.Flags().StringVarP(&opts.upstream, "upstream", "", "", "upstream dns server")
	rootCmd.Flags().StringSliceVarP(&opts.ipsV4, "ipv4", "", nil, "IPV4 address to return")
	rootCmd.Flags().StringSliceVarP(&opts.ipsV6, "ipv6", "", nil, "IPV6 address to return")
	rootCmd.Flags().Uint32VarP(&opts.ttl, "ttl", "", defaultTTL, "time to live")
	rootCmd.Flags().StringVarP(&opts.rebindV4, "rebind-v4", "", "", "IPV4 rebind address")
	rootCmd.Flags().StringVarP(&opts.rebindV6, "rebind-v6", "", "", "IPV6 rebind address")
	rootCmd.Flags().IntVarP(&opts.rebindThreshold, "rebind-threshold", "", 1, "rebind threshold")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
