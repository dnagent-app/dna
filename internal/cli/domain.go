package cli

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/datatug/cliformat"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo2namecheap"
	"github.com/spf13/cobra"
)

func domainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Manage domains",
	}
	cmd.AddCommand(domainListCmd())
	return cmd
}

func domainListCmd() *cobra.Command {
	var page int
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List registered domains",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDomainList(cmd, page, limit)
		},
	}
	cmd.Flags().IntVar(&page, "page", 1, "page number (1-based)")
	cmd.Flags().IntVar(&limit, "limit", 20, "domains per page (max 100)")
	cliformat.AddFlag(cmd, "name")
	return cmd
}

func runDomainList(cmd *cobra.Command, page, limit int) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	q := domainQuery{offset: (page - 1) * limit, limit: limit}
	reader, err := client.DomainsCollection().ExecuteQueryToRecordsReader(ctx, q)
	if err != nil {
		return fmt.Errorf("listing domains: %w", err)
	}
	defer reader.Close()

	var items []namecheap.DomainInfo
	for {
		rec, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading domain: %w", err)
		}
		if info, ok := rec.Data().(*namecheap.DomainInfo); ok {
			items = append(items, *info)
		}
	}

	format := cliformat.GetFormat(cmd)
	return cliformat.WriteList(cmd.OutOrStdout(), format, items, func(d namecheap.DomainInfo) string {
		return d.DomainName
	})
}

// newClient builds a namecheap.Client from env / ~/.namecheap-api credentials.
func newClient() (*namecheap.Client, error) {
	opts, err := namecheap.ConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("loading credentials: %w\n\nSet NAMECHEAP_API_USER / NAMECHEAP_API_KEY, or add them to ~/.namecheap-api", err)
	}
	opts = append(opts, namecheap.WithClientIPAutodetection())
	return namecheap.New(opts...)
}

// domainQuery implements dal.Query for the domain list command.
type domainQuery struct {
	offset int
	limit  int
}

func (q domainQuery) String() string { return "domain list" }
func (q domainQuery) Offset() int    { return q.offset }
func (q domainQuery) Limit() int     { return q.limit }
func (q domainQuery) GetRecordsReader(ctx context.Context, qe dal.QueryExecutor) (dal.RecordsReader, error) {
	return qe.ExecuteQueryToRecordsReader(ctx, q)
}
func (q domainQuery) GetRecordsetReader(ctx context.Context, qe dal.QueryExecutor) (dal.RecordsetReader, error) {
	return qe.ExecuteQueryToRecordsetReader(ctx, q)
}
