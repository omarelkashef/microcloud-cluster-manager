package main

import (
	"errors"
	"github.com/canonical/lxd-site-manager/client"
	cli "github.com/canonical/lxd/shared/cmd"
	microClient "github.com/canonical/microcluster/client"
	"github.com/canonical/microcluster/microcluster"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"sort"
	"strings"
)

type cmdSite struct {
	common *CmdControl
}

func (c *cmdSite) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "Manage sites",
		RunE:  c.Run,
	}

	siteListCmd := &cmdSiteList{
		common: c.common,
		site:   c,
	}

	cmd.AddCommand(siteListCmd.Command())

	siteShowCmd := &cmdSiteShow{
		common: c.common,
		site:   c,
	}

	cmd.AddCommand(siteShowCmd.Command())

	return cmd
}

func (c *cmdSite) Run(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

type cmdSiteList struct {
	common *CmdControl
	site   *cmdSite
}

func (c *cmdSiteList) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sites",
		RunE:  c.Run,
	}

	return cmd
}

func (c *cmdSiteList) Run(cmd *cobra.Command, args []string) error {
	m, err := microcluster.App(microcluster.Args{StateDir: c.common.FlagStateDir, Verbose: c.common.FlagLogVerbose, Debug: c.common.FlagLogDebug})
	if err != nil {
		return err
	}

	var microclusterClient *microClient.Client
	if len(args) == 1 {
		microclusterClient, err = m.RemoteClient(args[0])
	} else if len(args) == 0 {
		microclusterClient, err = m.LocalClient()
	} else {
		return errors.New("Invalid number of arguments")
	}

	if err != nil {
		return err
	}

	sites, err := client.GetSites(cmd.Context(), microclusterClient)
	if err != nil {
		return err
	}

	headers := []string{
		"NAME",
		"ADDRESSES",
		"STATUS",
	}

	var rows [][]string
	for _, site := range sites {
		rows = append(rows, []string{
			site.Name,
			strings.Join(site.Addresses, "\n"),
			site.Status,
		})
	}

	sort.Sort(cli.SortColumnsNaturally(rows))
	return cli.RenderTable(cli.TableFormatTable, headers, rows, sites)
}

type cmdSiteShow struct {
	common *CmdControl
	site   *cmdSite
}

func (c *cmdSiteShow) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Show a site",
		RunE:  c.Run,
	}

	return cmd
}

func (c *cmdSiteShow) Run(cmd *cobra.Command, args []string) error {
	m, err := microcluster.App(microcluster.Args{StateDir: c.common.FlagStateDir, Verbose: c.common.FlagLogVerbose, Debug: c.common.FlagLogDebug})
	if err != nil {
		return err
	}

	var siteName string
	var microclusterClient *microClient.Client
	if len(args) == 1 {
		siteName = args[0]
		microclusterClient, err = m.LocalClient()
	} else if len(args) == 2 {
		microclusterClient, err = m.RemoteClient(args[0])
		siteName = args[1]
	} else {
		return errors.New("Invalid number of arguments")
	}

	if err != nil {
		return err
	}

	site, err := client.GetSite(cmd.Context(), microclusterClient, siteName)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(os.Stdout).Encode(site)
}
