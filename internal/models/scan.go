package models

import "time"

type ScanRequest struct {
	Input string `json:"input"`
	Type  string `json:"type"`
}

type ScanResponse struct {
	ScanID    string    `json:"scanId"`
	Status    string    `json:"status"`
	Input     string    `json:"input"`
	CreatedAt time.Time `json:"created_at"`
}

type ScanResultWrapper struct {
	Success bool       `json:"success"`
	Data    ScanResult `json:"data"`
}

type ScanResult struct {
	ID               string                 `json:"id"`
	ScanID           string                 `json:"scanId"`
	Target           string                 `json:"input"`
	Status           string                 `json:"status"`
	Progress         int                    `json:"progress"`
	Scoring          Scoring                `json:"scoring"`
	ReconData        *ReconData             `json:"recon_data"`
	ThreatIntelligence *ThreatIntelligence  `json:"threat_intelligence"`
	GraphData        *GraphData             `json:"graph_data"`
	AISummary        *AISummary             `json:"ai_summary"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Errors           []string               `json:"errors"`
}

type Scoring struct {
	RiskScore       float64  `json:"risk_score"`
	ConfidenceScore float64  `json:"confidence_score"`
	RiskLevel       string   `json:"risk_level"`
	Factors         []string `json:"factors"`
}

type ReconData struct {
	DomainInfo DomainInfo     `json:"domain_info"`
	GeoIPInfo  *GeoIPInfo     `json:"geoip_info"`
	DNSRecords map[string]any `json:"dns_records"`
	WhoisData  *WhoisData     `json:"whois_data"`
}

type DomainInfo struct {
	Domain     string `json:"domain"`
	Registered bool   `json:"registered"`
	Target     string `json:"target"`
}

type WhoisData struct {
	Registrar      string   `json:"registrar"`
	CreationDate   string   `json:"creationDate"`
	ExpirationDate string   `json:"expirationDate"`
	Nameservers    []string `json:"nameservers"`
}

type GeoIPInfo struct {
	IP              string   `json:"ip"`
	ASN             string   `json:"asn"`
	ASNName         string   `json:"asnName"`
	Country         string   `json:"country"`
	HostingProvider string   `json:"hostingProvider"`
	RiskTags        []string `json:"riskTags"`
}

type ThreatIntelligence struct {
	Summary        string           `json:"summary"`
	RelatedDomains []string         `json:"related_domains"`
	RelatedIPs     []string         `json:"related_ips"`
	ReputationFeeds map[string]Feed  `json:"reputation_feeds"`
}

type Feed struct {
	Detected bool   `json:"detected"`
	Verdict  string `json:"verdict"`
	Score    int    `json:"score"`
}

type GraphData struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

type GraphNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type GraphEdge struct {
	Source       string `json:"source"`
	Target       string `json:"target"`
	Relationship string `json:"relationship"`
}

type AISummary struct {
	Summary     string   `json:"summary"`
	Remediation []string `json:"remediation"`
	AttackChain string   `json:"attack_chain"`
}

type ScanStatusResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID              string  `json:"id"`
		Status          string  `json:"status"`
		Progress        int     `json:"progress"`
		RiskScore       float64 `json:"risk_score"`
		RiskLevel       string  `json:"risk_level"`
		ConfidenceScore float64 `json:"confidence_score"`
		ErrorMessage    *string `json:"error_message"`
	} `json:"data"`
}

type ThreatSignal struct {
	Name     string `json:"name"`
	Severity string `json:"severity"`
}

type ThreatActor struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}
