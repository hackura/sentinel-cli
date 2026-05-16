package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/hackura/sentinel-cli/internal/auth"
	"github.com/hackura/sentinel-cli/internal/config"
	"github.com/hackura/sentinel-cli/internal/ui"
	"github.com/hackura/sentinel-cli/internal/api"
)

func main() {
	app := &cli.App{
		Name:    "sentinel",
		Usage:   "Hackura Sentinel AI Security CLI",
		Version: "1.0.0",
		Action: func(c *cli.Context) error {
			ui.PrintBanner()
			return cli.ShowAppHelp(c)
		},
		Commands: []*cli.Command{
			{
				Name:  "login",
				Usage: "Authenticate with your Sentinel account",
				Action: loginAction,
			},
			{
				Name:  "logout",
				Usage: "Revoke token and logout",
				Action: logoutAction,
			},
			{
				Name:  "status",
				Usage: "Check authentication status",
				Action: statusAction,
			},
			{
				Name:  "scan",
				Usage: "Scan a URL, domain, IP or file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "file",
						Usage: "Batch scan from file",
					},
					&cli.BoolFlag{
						Name:  "async",
						Usage: "Background scan",
					},
					&cli.StringFlag{
						Name:  "format",
						Usage: "Output format (text, json, csv)",
						Value: "text",
					},
				},
				Action: scanAction,
			},
			{
				Name:  "domain",
				Usage: "Domain intelligence and analysis",
				Action: domainAction,
			},
			{
				Name:  "ip",
				Usage: "IP address intelligence",
				Action: ipAction,
			},
			{
				Name:  "graph",
				Usage: "Visualize threat graph for a target",
				Action: graphAction,
			},
			{
				Name:  "explain",
				Usage: "AI explanation and risk summary",
				Action: explainAction,
			},
			{
				Name:  "history",
				Usage: "Show scan history",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "format",
						Usage: "Output format (text, json, csv)",
						Value: "text",
					},
				},
				Action: resultsListAction,
			},
			{
				Name:  "config",
				Usage: "Manage CLI configuration",
				Subcommands: []*cli.Command{
					{
						Name:  "set",
						Usage: "Set config key value",
						Action: configSetAction,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}
}

func loginAction(c *cli.Context) error {
	ui.PrintBanner()
	
	deviceID, err := auth.GenerateDeviceID()
	if err != nil {
		return err
	}

	cfg, _ := config.LoadConfig()
	apiURL := cfg.Profiles[cfg.ActiveProfile].APIURL

	ui.PrintInfo("Initiating authentication sequence...")

	// Step 1: Start auth session via backend
	startURL := fmt.Sprintf("%s/cli/auth/start", apiURL)
	payload := map[string]string{"device_id": deviceID}
	payloadBytes, _ := json.Marshal(payload)
	resp, err := http.Post(startURL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to reach backend: %v", err))
		return err
	}
	defer resp.Body.Close()

	var startResult struct {
		Success bool   `json:"success"`
		AuthURL string `json:"auth_url"`
		Error   string `json:"error"`
	}
	json.NewDecoder(resp.Body).Decode(&startResult)

	if !startResult.Success {
		return fmt.Errorf("failed to start auth session: %s", startResult.Error)
	}

	maskedID := auth.MaskDeviceID(deviceID)
	fmt.Printf(" [>] Authentication Portal: %s\n", ui.InfoStyle.Render(startResult.AuthURL))
	fmt.Printf(" [>] Device Identity: %s\n\n", ui.DimStyle.Render(maskedID))

	if err := auth.OpenBrowser(startResult.AuthURL); err != nil {
		ui.PrintInfo("Please open the URL manually in your browser.")
	}

	ui.PrintInfo("Waiting for authorization... ⚡")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	session, err := auth.PollForToken(ctx, apiURL, deviceID)
	if err != nil {
		return err
	}

	if err := config.SaveToken(session.Token); err != nil {
		return err
	}

	fmt.Println()
	ui.PrintSuccess("Authentication successful!")
	
	// Fetch real status for the success message
	cfg, _ = config.LoadConfig()
	client := api.NewClient(cfg.Profiles[cfg.ActiveProfile].APIURL)
	status, err := client.GetUserStatus()
	
	displayName := "hacker"
	if err == nil && status.Name != "" {
		displayName = status.Name
	}

	welcomeMsg := fmt.Sprintf("Welcome back, %s. You are now logged in.", displayName)
	
	info := fmt.Sprintf("%s\n\n%-15s %s\n%-15s %s\n%-15s %s\n",
		ui.HeaderStyle.Render(welcomeMsg),
		ui.LabelStyle.Render("Account:"), session.Email,
		ui.LabelStyle.Render("Plan:"), status.Plan,
		ui.LabelStyle.Render("Daily Scans:"), fmt.Sprintf("%d / %d used", status.ScansUsed, status.ScansLimit),
	)
	
	fmt.Println(ui.BoxStyle.Render(info))

	return nil
}

func logoutAction(c *cli.Context) error {
	if err := config.RemoveToken(); err != nil {
		return err
	}
	ui.PrintSuccess("Session revoked. Logged out.")
	return nil
}

func statusAction(c *cli.Context) error {
	cfg, _ := config.LoadConfig()
	client := api.NewClient(cfg.Profiles[cfg.ActiveProfile].APIURL)

	status, err := client.GetUserStatus()
	if err != nil {
		ui.PrintInfo("Not logged in or session expired. Run 'sentinel login' to begin.")
		return nil
	}

	ui.PrintSuccess(fmt.Sprintf("Authenticated as %s", status.Email))
	
	info := fmt.Sprintf("%-15s %s\n%-15s %s\n%-15s %s\n",
		ui.LabelStyle.Render("User:"), status.Name,
		ui.LabelStyle.Render("Plan:"), status.Plan,
		ui.LabelStyle.Render("Daily Scans:"), fmt.Sprintf("%d / %d used", status.ScansUsed, status.ScansLimit),
	)

	fmt.Println(ui.BoxStyle.Render(info))
	return nil
}

func scanAction(c *cli.Context) error {
	target := c.Args().First()
	file := c.String("file")
	isAsync := c.Bool("async")
	format := c.String("format")

	if target == "" && file == "" {
		return fmt.Errorf("target URL/domain or --file is required")
	}

	cfg, _ := config.LoadConfig()
	client := api.NewClient(cfg.Profiles[cfg.ActiveProfile].APIURL)

	if file != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		targets := bytes.Split(data, []byte("\n"))
		for _, t := range targets {
			target := string(bytes.TrimSpace(t))
			if target == "" {
				continue
			}
			ui.PrintInfo(fmt.Sprintf("Scanning %s...", target))
			res, err := client.InitiateScan(target, isAsync)
			if err != nil {
				ui.PrintError(fmt.Sprintf("Failed to initiate scan for %s: %v", target, err))
				continue
			}
			fmt.Printf(" [>] Scan ID: %s\n", res.ScanID)
		}
		return nil
	}

	if !isValidTarget(target) {
		return fmt.Errorf("invalid target format")
	}

	ui.PrintInfo(fmt.Sprintf("Scanning %s...", target))
	
	res, err := client.InitiateScan(target, isAsync)
	if err != nil {
		return err
	}

	if isAsync {
		ui.PrintSuccess("Scan initiated in background.")
		fmt.Printf(" [>] Scan ID: %s\n", res.ScanID)
		return nil
	}

	// Poll for results
	ui.PrintInfo("Waiting for results... ⚡")
	for {
		status, err := client.GetScanStatus(res.ScanID)
		if err != nil {
			return err
		}

		if status.Data.Status == "completed" {
			ui.PrintSuccess("Scan complete!")
			results, err := client.GetScanResults(res.ScanID)
			if err != nil {
				return err
			}
			return ui.PrintFormattedResults(results, format)
		}

		if status.Data.Status == "failed" {
			errMsg := "unknown error"
			if status.Data.ErrorMessage != nil {
				errMsg = *status.Data.ErrorMessage
			}
			return fmt.Errorf("scan failed: %s", errMsg)
		}

		time.Sleep(2 * time.Second)
	}
}

func domainAction(c *cli.Context) error {
	target := c.Args().First()
	if target == "" {
		return fmt.Errorf("domain is required")
	}
	ui.PrintInfo(fmt.Sprintf("Gathering intelligence for %s...", target))
	
	cfg, _ := config.LoadConfig()
	client := api.NewClient(cfg.Profiles[cfg.ActiveProfile].APIURL)
	
	res, err := client.InitiateScan(target, false)
	if err != nil {
		return err
	}
	
	return pollAndShow(client, res.ScanID, "domain")
}

func ipAction(c *cli.Context) error {
	target := c.Args().First()
	if target == "" {
		return fmt.Errorf("IP address is required")
	}
	ui.PrintInfo(fmt.Sprintf("Analyzing IP address %s...", target))
	
	cfg, _ := config.LoadConfig()
	client := api.NewClient(cfg.Profiles[cfg.ActiveProfile].APIURL)
	
	res, err := client.InitiateScan(target, false)
	if err != nil {
		return err
	}
	
	return pollAndShow(client, res.ScanID, "ip")
}

func graphAction(c *cli.Context) error {
	target := c.Args().First()
	if target == "" {
		return fmt.Errorf("target is required")
	}
	ui.PrintInfo(fmt.Sprintf("Building threat graph for %s...", target))
	
	cfg, _ := config.LoadConfig()
	client := api.NewClient(cfg.Profiles[cfg.ActiveProfile].APIURL)
	
	res, err := client.InitiateScan(target, false)
	if err != nil {
		return err
	}
	
	return pollAndShow(client, res.ScanID, "graph")
}

func explainAction(c *cli.Context) error {
	target := c.Args().First()
	if target == "" {
		return fmt.Errorf("target is required")
	}
	ui.PrintInfo(fmt.Sprintf("Generating AI analysis for %s...", target))
	
	cfg, _ := config.LoadConfig()
	client := api.NewClient(cfg.Profiles[cfg.ActiveProfile].APIURL)
	
	res, err := client.InitiateScan(target, false)
	if err != nil {
		return err
	}
	
	return pollAndShow(client, res.ScanID, "explain")
}

func pollAndShow(client *api.Client, scanID, mode string) error {
	ui.PrintInfo("Waiting for results... ⚡")
	for {
		status, err := client.GetScanStatus(scanID)
		if err != nil {
			return err
		}

		if status.Data.Status == "completed" {
			results, err := client.GetScanResults(scanID)
			if err != nil {
				return err
			}
			
			switch mode {
			case "domain":
				ui.PrintDomainInfo(results)
			case "ip":
				ui.PrintIPInfo(results)
			case "graph":
				ui.PrintGraph(results)
			case "explain":
				ui.PrintAIExplanation(results)
			default:
				ui.PrintScanResults(results)
			}
			return nil
		}

		if status.Data.Status == "failed" {
			return fmt.Errorf("analysis failed")
		}

		time.Sleep(2 * time.Second)
	}
}

func resultsListAction(c *cli.Context) error {
	format := c.String("format")
	cfg, _ := config.LoadConfig()
	client := api.NewClient(cfg.Profiles[cfg.ActiveProfile].APIURL)

	results, err := client.ListRecentResults()
	if err != nil {
		return err
	}

	if len(results) == 0 {
		ui.PrintInfo("No scan history found.")
		return nil
	}

	if format == "text" {
		ui.PrintTitle("Recent Scans")
		header := fmt.Sprintf(" %-5s %-30s %-10s %-10s %-20s\n", "ID", "Target", "Type", "Risk", "Date")
		table := header + strings.Repeat("-", 100) + "\n"
		for i, r := range results {
			riskLabel := "Low"
			riskStyle := ui.SuccessStyle
			if r.Scoring.RiskScore >= 70 {
				riskLabel = "High"
				riskStyle = ui.ErrorStyle
			} else if r.Scoring.RiskScore >= 35 {
				riskLabel = "Medium"
				riskStyle = ui.WarningStyle
			}
			
			table += fmt.Sprintf(" %-5d %-30s %-10s %-10s %-20s\n", i+1, r.Target, "domain", riskStyle.Render(riskLabel), r.CreatedAt.Format("2006-01-02 15:04"))
		}
		fmt.Println(ui.RenderWithMask(table, ui.HackerShield))
	} else {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
	}
	return nil
}

func configSetAction(c *cli.Context) error {
	key := c.Args().Get(0)
	val := c.Args().Get(1)

	if key == "" || val == "" {
		return fmt.Errorf("usage: sentinel config set <key> <val>")
	}

	cfg, _ := config.LoadConfig()
	profile := cfg.Profiles[cfg.ActiveProfile]

	switch key {
	case "api_url":
		profile.APIURL = val
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	cfg.Profiles[cfg.ActiveProfile] = profile
	if err := config.SaveConfig(cfg); err != nil {
		return err
	}

	ui.PrintSuccess(fmt.Sprintf("Set %s to %s", key, val))
	return nil
}

func isValidTarget(target string) bool {
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		_, err := url.ParseRequestURI(target)
		return err == nil
	}
	return strings.Contains(target, ".") && !strings.Contains(target, " ")
}
