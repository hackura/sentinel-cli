package api

import (
	"fmt"
	"strings"
	"net"
	"github.com/hackura/sentinel-cli/internal/models"
)

func (c *Client) InitiateScan(target string, async bool) (*models.ScanResponse, error) {
	scanType := detectType(target)
	req := models.ScanRequest{
		Input: target,
		Type:  scanType,
	}
	
	// The backend returns { "success": true, "data": { "scanId": "...", "status": "..." } } 
	// for POST /scan or different for /scan/initiate
	// Based on routes/scan.ts: POST /initiate returns { success: true, scanId: ... }
	
	var res struct {
		Success bool                `json:"success"`
		ScanID  string              `json:"scanId"`
		Data    models.ScanResponse `json:"data"`
	}
	
	err := c.Post("/scan", req, &res)
	if err != nil {
		return nil, err
	}
	
	finalRes := &res.Data
	if res.ScanID != "" {
		finalRes.ScanID = res.ScanID
	}
	
	return finalRes, nil
}


func (c *Client) GetScanStatus(id string) (*models.ScanStatusResponse, error) {
	var res models.ScanStatusResponse
	err := c.Get(fmt.Sprintf("/scan/%s/status", id), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetScanResults(id string) (*models.ScanResult, error) {
	var res models.ScanResultWrapper
	err := c.Get(fmt.Sprintf("/scan/%s", id), &res)
	if err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) ListRecentResults() ([]models.ScanResult, error) {
	// The backend might use /scans or /results. routes/scans.ts likely handles listing.
	var res struct {
		Success bool                `json:"success"`
		Data    []models.ScanResult `json:"data"`
	}
	err := c.Get("/scans", &res)
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}

func detectType(input string) string {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return "url"
	}
	if net.ParseIP(input) != nil {
		return "ip"
	}
	// Simple check for domain vs text
	if strings.Contains(input, ".") && !strings.Contains(input, " ") {
		return "domain"
	}
	return "text"
}
