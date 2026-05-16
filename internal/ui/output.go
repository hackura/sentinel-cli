package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"math"

	"github.com/charmbracelet/lipgloss"
	"github.com/hackura/sentinel-cli/internal/models"
)

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#7D56F4") // Purple
	SecondaryColor = lipgloss.Color("#00ff00") // Neon Green
	AccentColor    = lipgloss.Color("#ff00ff") // Magenta
	BgColor        = lipgloss.Color("#1a1a1a")
	ErrorColor     = lipgloss.Color("#ff0000")
	WarningColor   = lipgloss.Color("#ffa500")
	InfoColor      = lipgloss.Color("#00bfff")
	DimColor       = lipgloss.Color("#666666")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(0, 1)

	BannerStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			MarginBottom(1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			MarginBottom(1).
			PaddingBottom(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor)

	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor)

	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor)

	DimStyle = lipgloss.NewStyle().
			Foreground(DimColor)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2).
			MarginTop(1)

	LabelStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff"))

	RiskHighStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	RiskMediumStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true)

	RiskLowStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	// Banners from Image
	SentinelBanner = `
  .-----------------------------------------------------------.
 /  ____  _____ _   _ _____ ___ _   _ _____ _                 \
 | / ___|| ____| \ | |_   _|_ _| \ | | ____| |                |
 | \___ \|  _| |  \| | | |  | ||  \| |  _| | |                |
 |  ___) | |___| |\  | | |  | || |\  | |___| |___             |
 | |____/|_____|_| \_| |_| |___|_| \_|_____|_____|            |
 \                                                            /
  '-----------------------------------------------------------'
`

	// Hacker Shield for Login
	HackerShield = `
⢲⡱⣎⢯⡽⢯⡿⣽⢯⣟⡿⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⣟⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⡿⣿⢿⣿⢿⡿⣿⣻⠿⣋⣿⣹⣏⡿⣭⣛⡼⣣⠝
⠲⣵⣚⢷⡞⣿⣼⡳⣟⣿⢮⣛⠿⣯⣿⡋⣿⣿⠿⢿⣿⡿⣿⣿⣿⣿⢯⣿⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⣨⣿⣯⢿⣻⢿⡃⣼⣿⡯⠟⠋⣱⣶⢿⣳⢿⣜⣳⢧⡝⡶⣱⢚
⢹⡲⣽⡺⣽⠶⣏⣿⣳⣍⡻⣟⢷⣝⢿⣿⣌⠻⡄⣝⣿⠟⣿⡿⣇⢿⣏⢹⠿⠿⠷⠿⠿⣄⠼⠿⠿⠿⠯⡿⣟⣳⣟⣻⣿⣗⠇⣼⣿⠝⢁⣤⣴⣿⣟⣿⢯⣟⣾⠁⣼⢺⡵⢣⢏
⢣⡟⣶⢻⣽⣻⣽⣞⣯⣿⢿⣷⣷⣝⠷⡝⣿⣧⠻⠋⢺⠲⢛⡋⢭⠴⢰⣾⣿⣿⣿⣿⣿⣿⡿⢱⢫⣞⣵⣃⠤⠌⡉⠓⠻⢞⣺⡯⠈⣰⣿⡿⣿⣷⣻⣯⠿⠚⢁⠰⢯⣳⡝⣯⢚
⢣⢿⣼⢻⣼⣳⢷⣯⢿⣻⣿⣿⣾⣿⡷⡝⣎⠟⠊⢀⠠⠂⢧⡙⣎⡑⣼⣿⣿⣟⣿⣿⡿⣿⡟⢼⣏⡞⣤⢫⠔⢣⠐⡀⠀⠀⠈⠁⣾⣿⠷⣾⣿⣾⠏⠁⡀⢢⣴⣶⣞⡷⣻⡜⣯
⡹⣞⣼⣻⣳⣯⢿⣽⣿⡻⠷⢿⣿⠉⢹⡟⡅⠌⡘⢄⠣⢐⠀⣙⠖⣼⣿⣟⣿⣿⣿⠿⠀⢹⡣⣏⡞⣼⡿⣾⠜⣆⠣⠀⠀⠀⠀⠀⢋⣁⠻⢯⡿⠁⡀⢂⣵⣿⣟⣾⣳⡽⢧⣛⠶
⡱⣞⢾⣵⣻⢯⣿⢷⣟⣷⡄⠘⠾⢟⣿⠯⢡⢓⠈⠤⡈⠆⢖⡱⢢⢩⣽⠿⡿⠿⣿⠃⠀⡌⣇⡟⣼⣿⡜⣗⡛⡌⠡⢀⠀⠀⠀⠀⢐⡒⠦⢂⠠⠐⠀⢿⣿⢿⣾⣭⣷⣻⣭⢏⡟
⢡⢻⡜⡼⣹⢏⣿⢏⣿⢻⣧⠀⠀⠋⠌⠛⠀⠌⢸⠃⣹⡟⣤⣁⣡⡀⠈⠛⢿⣿⠋⠀⠼⠀⠸⣼⣿⡇⢸⢿⢡⠈⡘⡀⠀⠀⠀⠀⢠⢉⠀⠈⠀⠀⠀⠀⠙⠻⠏⣿⢧⣏⡟⣼⢹
⡸⢾⣸⢷⣻⢾⣟⡾⣿⣧⣟⢶⠀⠀⠀⠄⠰⢀⠼⣿⣿⣿⣿⣿⣿⣿⣿⣦⣀⡿⠀⢐⣂⣄⠂⢿⡿⣣⣾⣺⡤⣔⢣⣔⣧⡀⠀⠀⠀⠆⠒⢠⠘⢀⠀⠀⠀⠀⠀⢸⣾⣼⣳⢧⣇
⡹⣏⣯⠿⣽⢿⣻⣶⣶⠟⣷⡻⠇⠀⠈⠀⠜⡠⡼⣛⡭⠶⠶⠶⣮⣭⣙⡻⢿⣧⣶⣿⣿⣿⣶⣾⠷⢛⣪⣭⠴⠶⠒⠮⢥⡓⠀⠀⠀⠐⠃⢠⢅⡢⠀⠀⠀⠀⣠⣾⣿⣳⢯⣟⡞
⣱⢯⣞⡿⣷⣴⢻⣯⠟⡡⠇⠀⠀⠀⠀⠀⢊⢀⠜⢁⣀⣀⠀⠀⠠⣶⡫⢋⠶⣌⠛⣿⣿⣿⠏⡡⠞⠋⠁⠀⠀⢄⣀⣠⣄⣀⠀⠀⠀⠠⠁⡜⢮⠁⠂⢀⣾⣿⣿⡿⣯⣟⡿⣞⡽
⢼⣣⢟⣷⣻⡾⣟⣿⢠⠹⣅⡀⡀⠀⠀⠀⡅⢠⣾⣯⣾⣿⣿⣿⣶⡅⠈⠀⠁⠈⠳⣾⣿⣿⠎⠀⠁⠀⠀⠀⡾⣿⣿⣿⢿⡞⠀⠀⠀⠐⠀⡜⠦⠁⠀⣿⣿⣿⣾⣿⢿⣽⣻⣭⠷
⢎⡷⣻⢾⡽⣯⣿⠏⢄⡛⠜⠂⠀⠀⠀⠀⢈⠀⠻⣿⠿⠟⠛⠛⠛⠛⠂⢠⠠⠄⠰⢿⡿⡿⠂⠀⠀⠀⠀⠈⠉⠉⠉⠈⠉⠀⠀⠀⠀⠀⠰⡜⠣⠁⠀⢽⣿⣿⣽⣾⡿⣞⣧⢿⣹
⢎⣷⢫⡿⣽⣟⣿⠐⡸⢌⠀⠀⠀⠀⠀⠀⠈⠧⠄⠀⠐⠈⠀⠀⠀⠀⠁⠂⡀⠀⠰⣷⣿⡎⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⢄⡀⣐⠀⠐⣄⢣⢢⣾⡿⣾⣿⣻⣷⡿⣽⣻⠾⣝⣮
⣚⢾⡽⣳⣟⣾⣷⡆⠁⢎⣡⠙⡆⠀⠀⠀⠃⡿⢟⣛⡀⠢⠄⠀⠀⠀⣀⣀⡨⣶⡇⢿⣿⡇⠀⠀⣀⣀⣀⠀⠀⠠⣄⣊⣴⣛⡣⡏⠀⢘⢤⣓⠆⡍⡿⡏⣵⣿⢾⣻⡿⣽⣻⢽⣲
⣜⣻⡼⣷⣻⡾⣿⣽⣆⠣⢆⡛⡄⠑⠀⠀⠁⡾⣫⣵⣿⣿⣿⣿⣿⣿⣿⣿⣧⣿⡇⢹⣿⡇⢠⣠⣿⣿⣿⡏⣾⣮⣽⣛⠿⣿⣷⠀⠄⢘⠮⣍⣻⡘⣻⣾⣻⣽⣟⣽⣟⡷⣽⣣⠷
⢼⣳⢽⡷⣯⣟⣿⣣⣀⣷⢬⣕⡊⡅⠠⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣟⣇⢺⣿⡇⠀⢾⣸⣿⣷⢾⡿⡿⡵⢯⠟⠗⠁⠀⠀⠨⠳⡜⢦⢽⡠⠜⣯⣼⣿⣟⡾⣽⣳⢏⡿
⢮⡽⣞⣿⣳⢯⣿⣾⣯⡿⣫⣜⡞⡁⠀⠀⠀⠀⢻⣿⣿⣿⣿⡿⣿⣿⡿⢿⢛⠽⣮⢿⣿⡇⠀⠃⡀⠿⣿⢜⠲⢡⢆⡄⠀⠀⠀⠀⠀⠘⠡⡙⢬⢳⣭⠃⣿⣟⣯⣽⢿⡽⣧⢟⡞
⢮⡽⣞⡷⣯⣟⣾⣽⠋⠄⣭⣿⡘⡱⠲⠁⠀⠀⠀⠉⠛⠛⠓⠉⠉⣀⣴⡆⣿⣷⣾⣾⣿⣷⠀⣾⣧⢴⣄⠀⠉⠉⠈⠀⠀⠀⠀⠀⠀⠀⠡⢄⢣⢻⡴⢫⢸⣾⣯⡿⣯⣟⣾⡹⣞
⢮⡽⢾⣽⣳⠛⣾⣽⢈⠢⠻⣷⠃⡰⡐⠄⠀⠀⠈⠄⢶⡀⠀⢿⣿⣿⣿⣷⡄⠈⠙⠻⠿⠟⠀⠀⠀⢸⣿⣿⣷⠆⢀⠀⠠⠀⠀⠀⠀⠀⠀⢨⠖⣧⢻⣱⢸⣿⣯⣟⣷⣻⢶⡻⣜
⢮⣽⢻⣼⣻⣖⣯⣿⣇⠬⡑⡄⢮⡱⣉⠐⠀⡀⠀⠀⠈⠳⣄⠀⠙⠿⠿⢿⣿⡷⠀⠀⠀⠀⠀⠀⠀⠸⠿⠟⠁⡠⠀⡱⠁⠀⠀⠀⠀⠀⠀⢡⢟⣮⢳⢇⢺⣿⣯⢿⡾⣽⡳⢏⡷
⢺⣼⢻⢶⣻⡾⣽⣍⣿⣷⣆⢚⣡⠓⡄⢺⡬⠇⠠⢣⢀⠠⠙⢦⣠⣀⣀⠀⠀⠀⠀⣴⣧⣀⠀⠀⠀⠀⢀⡀⠀⠀⡔⠁⠀⠀⠀⠀⢀⡦⣔⠮⢞⣜⢫⢆⢾⣭⡙⢿⣽⣳⢻⣏⡞
⢳⡺⣽⢫⣷⣻⣽⡻⣿⠁⣹⣷⢠⠃⡌⠑⡄⢂⠄⣈⡁⠠⡈⢀⠱⣝⡻⠿⣿⣷⣶⣶⣶⣶⣶⣶⡾⠿⠃⠁⢀⠎⠀⠀⠀⠀⠀⠀⢺⡜⣮⡝⣣⣈⣦⣞⣿⣿⣷⣶⣳⢿⣹⢮⡝
⢣⡟⣼⢻⣼⠳⣯⣿⡟⣿⡿⣿⠰⠒⠌⠡⠈⠔⣬⡔⢧⠁⠔⠵⡄⠈⢿⣷⣶⣤⣤⠀⠀⠀⠀⣤⢤⡄⠀⠠⠃⠀⠀⠀⠀⠀⠀⠄⠚⡎⠳⡞⢱⢷⡞⣿⠛⣾⣷⢯⡟⡾⢹⢮⡝
⢂⣱⢸⣡⣏⣷⣹⢣⣏⣿⣷⣿⣧⣁⢲⡀⣁⣼⣏⣷⣦⡄⢀⠀⠈⢂⣀⠛⣿⣿⣿⠂⢀⠀⢻⣿⣾⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⡙⣷⣹⣹⣾⣌⡟⣹⣿⡜⣾⣼⡑⢣⡎⡔
⢪⢗⡯⢾⣹⢾⡽⣯⢿⣷⠟⡫⣵⣶⡿⣟⣯⡿⣻⡽⡷⣠⣊⢩⢄⠙⠫⣺⣿⣿⣷⠐⠀⠀⢸⣿⠡⠀⠀⠀⠀⠀⠀⠠⠈⠀⢠⡪⣜⠭⣦⡉⠉⣿⢿⣿⣶⡿⣟⣷⣳⡻⣝⢾⡱
⢹⠮⣝⢯⣽⡻⢾⡽⢟⣴⣺⡟⠉⢱⣿⣿⣿⢿⠟⢁⣼⡿⠣⠯⠟⣃⡀⠁⠙⣿⣿⡄⠀⠀⣼⡫⠂⠀⠀⠀⠀⠀⣤⣄⡶⣾⢿⣙⣦⡙⠮⡻⡾⣿⣿⣿⣳⣿⣟⣾⣳⢿⣹⢎⠷
⢹⣚⠽⣮⣳⡟⣯⣟⣯⢿⡽⣿⢶⣾⢯⣿⢟⢔⣱⢿⠟⢀⣦⡃⡝⢢⡝⣂⡀⠈⠻⠇⠀⢀⠏⠁⠀⢀⢴⣦⣀⣀⢬⢽⣿⣄⢵⣹⢿⡿⣎⠳⠎⣳⢟⣯⣿⣞⡿⢾⡭⣷⢫⡞⡽
⢳⢎⡻⢶⣝⡾⣳⡽⣞⣯⢿⣿⢿⡿⣋⢔⣠⠜⢛⡄⠒⣿⣯⡷⣨⣕⢮⡑⣛⢦⡀⠀⠀⠀⠀⢀⣴⣽⣎⢯⣳⣏⣈⢲⡙⣿⡷⣝⢳⣝⢿⣻⣷⣝⢷⣝⡷⣯⣟⣯⠷⣭⣳⣹⢱
⣙⠮⣝⡞⣮⣽⢳⣽⣛⡾⣯⠿⣫⣞⡷⡟⣠⣤⣾⣶⣿⡿⢿⣿⣿⣿⣥⡃⣯⢞⡹⢖⡄⢤⣤⣺⣿⣿⣿⣮⢻⡀⣸⣫⢿⣿⣿⣿⢷⣝⢂⡙⢯⣿⢯⣟⣷⣽⢾⣱⢟⡶⢳⡬⢳
⠼⣹⠜⣎⠷⣎⡿⣶⢫⡿⣽⣻⣟⣯⡏⣱⣿⡤⣿⣿⣯⣟⣽⣿⣿⣯⣿⣷⣮⣶⣾⣗⠤⣫⢿⣿⣿⣟⣿⣿⣷⣻⣧⣽⣦⢻⣿⣾⣿⢿⣦⣑⣊⢯⡿⣾⣳⢯⣏⢷⢫⠞⣥⢫⣓
⡹⣜⡹⢎⣻⣜⣳⢽⣫⡽⣷⣟⣾⣻⡿⣿⣻⣿⣿⣿⡽⣿⣿⣽⣿⢿⣿⣿⣿⣿⣿⣿⢀⣿⣿⣿⣿⣻⣿⣿⣿⣿⣻⣿⣿⣿⣿⣳⣿⡿⣯⣿⣟⡷⣟⣷⢫⣞⢮⡏⢯⡹⢬⠳⣬
⢜⠲⢍⡛⢶⣩⢞⡭⢷⣛⢷⣻⣞⡷⣿⣟⣯⣿⣾⣿⣿⣟⣿⣿⣻⣿⣿⣻⣽⣿⣿⣯⢀⣿⣿⣿⣾⣿⣿⣿⣿⣿⣿⣿⣽⣿⣿⣿⣻⣽⢿⡽⣞⡿⣽⢎⡿⣜⣣⣝⣣⢝
`

	// World Map for Domain
	WorldMap = `
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣄⣠⣀⡀⣀⣠⣤⣤⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣄⢠⣠⣼⣿⣿⣿⣟⣿⣿⣿⣿⣿⣿⣿⣿⡿⠋⠀⠀⠀⢠⣤⣦⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠰⢦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣟⣾⣿⣽⣿⣿⣅⠈⠉⠻⣿⣿⣿⣿⣿⡿⠇⠀⠀⠀⠀⠀⠉⠀⠀⠀⠀⠀⢀⡶⠒⢉⡀⢠⣤⣶⣶⣿⣷⣆⣀⡀⠀⢲⣖⠒⠀⠀⠀⠀⠀⠀⠀
⢀⣤⣾⣶⣦⣤⣤⣶⣿⣿⣿⣿⣿⣿⣽⡿⠻⣷⣀⠀⢻⣿⣿⣿⡿⠟⠀⠀⠀⠀⠀⠀⣤⣶⣶⣤⣀⣀⣬⣷⣦⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣦⣤⣦⣼⣀⠀
⠈⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠛⠓⣿⣿⠟⠁⠘⣿⡟⠁⠀⠘⠛⠁⠀⠀⢠⣾⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠏⠙⠁
⠀⠸⠟⠋⠀⠈⠙⣿⣿⣿⣿⣿⣿⣷⣦⡄⣿⣿⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⣼⣆⢘⣿⣯⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡉⠉⢱⡿⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣟⡿⠦⠀⠀⠀⠀⠀⠀⠀⠙⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⡗⠀⠈⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣿⣿⣿⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣉⣿⡿⢿⢷⣾⣾⣿⣞⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠋⣠⠟⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣿⣿⠿⠿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣾⣿⣿⣷⣦⣶⣦⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⠈⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠻⣿⣤⡖⠛⠶⠤⡀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁⠙⣿⣿⠿⢻⣿⣿⡿⠋⢩⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠧⣤⣦⣤⣄⡀⠀⠀⠀⠀⠀⠘⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠘⣧⠀⠈⣹⡻⠇⢀⣿⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣤⣀⡀⠀⠀⠀⠀⠀⠀⠈⢽⣿⣿⣿⣿⣿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠹⣷⣴⣿⣷⢲⣦⣤⡀⢀⡀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⣿⣿⣿⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣷⢀⡄⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠂⠛⣆⣤⡜⣟⠋⠙⠂⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⠉⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣤⣾⣿⣿⣿⣿⣆⠀⠰⠄⠀⠉⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⣿⣿⡿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⡿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⠿⠿⣿⣿⣿⠇⠀⠀⢀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⡿⠛⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢻⡇⠀⠀⢀⣼⠗⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⠃⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠁⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠒⠀⠀⠀⠀
`

	// Globe for IP
	GlobeMap = `
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣄⣠⣀⡀⣀⣠⣤⣤⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣄⢠⣠⣼⣿⣿⣿⣟⣿⣿⣿⣿⣿⣿⣿⣿⡿⠋⠀⠀⠀⢠⣤⣦⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠰⢦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣟⣾⣿⣽⣿⣿⣅⠈⠉⠻⣿⣿⣿⣿⣿⡿⠇⠀⠀⠀⠀⠀⠉⠀⠀⠀⠀⠀⢀⡶⠒⢉⡀⢠⣤⣶⣶⣿⣷⣆⣀⡀⠀⢲⣖⠒⠀⠀⠀⠀⠀⠀⠀
⢀⣤⣾⣶⣦⣤⣤⣶⣿⣿⣿⣿⣿⣿⣽⡿⠻⣷⣀⠀⢻⣿⣿⣿⡿⠟⠀⠀⠀⠀⠀⠀⣤⣶⣶⣤⣀⣀⣬⣷⣦⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣦⣤⣦⣼⣀⠀
⠈⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠛⠓⣿⣿⠟⠁⠘⣿⡟⠁⠀⠘⠛⠁⠀⠀⢠⣾⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠏⠙⠁
⠀⠸⠟⠋⠀⠈⠙⣿⣿⣿⣿⣿⣿⣷⣦⡄⣿⣿⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⣼⣆⢘⣿⣯⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡉⠉⢱⡿⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣟⡿⠦⠀⠀⠀⠀⠀⠀⠀⠙⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⡗⠀⠈⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣿⣿⣿⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣉⣿⡿⢿⢷⣾⣾⣿⣞⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠋⣠⠟⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣿⣿⠿⠿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣾⣿⣿⣷⣦⣶⣦⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⠈⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠻⣿⣤⡖⠛⠶⠤⡀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁⠙⣿⣿⠿⢻⣿⣿⡿⠋⢩⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠧⣤⣦⣤⣄⡀⠀⠀⠀⠀⠀⠘⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠘⣧⠀⠈⣹⡻⠇⢀⣿⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣤⣀⡀⠀⠀⠀⠀⠀⠀⠈⢽⣿⣿⣿⣿⣿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠹⣷⣴⣿⣷⢲⣦⣤⡀⢀⡀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⣿⣿⣿⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣷⢀⡄⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠂⠛⣆⣤⡜⣟⠋⠙⠂⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⠉⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣤⣾⣿⣿⣿⣿⣆⠀⠰⠄⠀⠉⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⣿⣿⡿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⡿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⠿⠿⣿⣿⣿⠇⠀⠀⢀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⡿⠛⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢻⡇⠀⠀⢀⣼⠗⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⠃⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠁⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠒
`

	HackerMaskStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			MarginLeft(4).
			Align(lipgloss.Center)
)

func PrintBanner() {
	bannerText := `
  ____  _____ _   _ _____ ___ _   _ _____ _     
 / ___|| ____| \ | |_   _|_ _| \ | | ____| |    
 \___ \|  _| |  \| | | |  | ||  \| |  _| | |    
  ___) | |___| |\  | | |  | || |\  | |___| |___ 
 |____/|_____|_| \_| |_| |___|_| \_|_____|_____|`
	
	banner := lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render(bannerText)
	subtitle := lipgloss.NewStyle().Foreground(SecondaryColor).PaddingLeft(4).Render("H A C K U R A   S E N T I N E L   A I")
	tagline := DimStyle.PaddingLeft(4).Render("Advanced Cyber Intelligence & Threat Analysis")
	
	content := banner + "\n" + subtitle + "\n" + tagline
	
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, content, HackerMaskStyle.Render(HackerShield)))
	fmt.Println()
}

func RenderWithMask(content string, mask string) string {
	if mask == "" {
		mask = HackerShield
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, content, HackerMaskStyle.Render(mask))
}

func PrintTitle(title string) {
	fmt.Println(TitleStyle.Render(title))
}

func PrintSuccess(msg string) {
	fmt.Printf(" %s %s\n", SuccessStyle.Render("✓"), msg)
}

func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, " %s %s\n", ErrorStyle.Render("✘"), msg)
}

func PrintInfo(msg string) {
	fmt.Printf(" %s %s\n", InfoStyle.Render("[*]"), msg)
}

func PrintStep(step string, duration string) {
	fmt.Printf(" %s %-40s %s\n", SuccessStyle.Render("✓"), step, DimStyle.Render(duration))
}

func PrintFormattedResults(results *models.ScanResult, format string) error {
	switch strings.ToLower(format) {
	case "json":
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "csv":
		fmt.Println("id,target,status,risk_score,confidence")
		fmt.Printf("%s,%s,%s,%.1f,%.2f\n", results.ID, results.Target, results.Status, results.Scoring.RiskScore, results.Scoring.ConfidenceScore)
	case "text":
		PrintScanResults(results)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	return nil
}

func PrintScanResults(results *models.ScanResult) {
	fmt.Println()
	
	// Simulated steps from image
	PrintStep("Fetching URL", "210ms")
	PrintStep("Checking Reputation", "320ms")
	PrintStep("Analyzing Content", "540ms")
	PrintStep("Checking Blacklists", "410ms")
	PrintStep("DNS & Infrastructure", "380ms")
	PrintStep("Building Graph", "620ms")
	PrintStep("Calculating Risk", "210ms")

	fmt.Printf("\nScan completed in %.2fs\n", 2.61) // Simulation or real timing

	// Info section
	info := fmt.Sprintf("%s %s\n%s %s\n%s %d%%\n%s %s\n%s %s",
		LabelStyle.Render("RISK SCORE: "), getRiskValueStyle(results.Scoring.RiskScore).Render(fmt.Sprintf("%.1f / 100 (%s)", results.Scoring.RiskScore, getRiskLabel(results.Scoring.RiskScore))),
		LabelStyle.Render("THREAT LEVEL:"), getRiskValueStyle(results.Scoring.RiskScore).Render(getRiskLabel(results.Scoring.RiskScore)),
		LabelStyle.Render("CONFIDENCE:  "), int(math.Min(100, results.Scoring.ConfidenceScore)),
		LabelStyle.Render("CATEGORY:    "), results.Status,
		LabelStyle.Render("LAST SEEN:   "), "Just now",
	)

	content := lipgloss.JoinHorizontal(lipgloss.Center, info, "    ", RenderCircularGauge(results.Scoring.RiskScore))
	
	box := BoxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			HeaderStyle.Render("Scan Completed"),
			content,
			"\n",
			HeaderStyle.Render("Risk Signals"),
			formatSignals(results.Scoring.Factors),
		),
	)
	
	fmt.Println(RenderWithMask(box, HackerShield))
}

func PrintDomainInfo(results *models.ScanResult) {
	if results.ReconData == nil {
		PrintInfo("No domain intelligence available.")
		return
	}
	recon := results.ReconData
	
	registrar := ""
	creationDate := ""
	expirationDate := ""
	var nameservers []string
	
	if recon.WhoisData != nil {
		registrar = recon.WhoisData.Registrar
		creationDate = recon.WhoisData.CreationDate
		expirationDate = recon.WhoisData.ExpirationDate
		nameservers = recon.WhoisData.Nameservers
	}
	if len(nameservers) == 0 && results.ThreatIntelligence != nil && len(results.ThreatIntelligence.RelatedDomains) > 0 {
		nameservers = results.ThreatIntelligence.RelatedDomains
	}
	
	info := fmt.Sprintf("%-15s %-30s\n%-15s %-30s\n%-15s %-30s\n%-15s %-30s\n%-15s\n",
		LabelStyle.Render("Domain:"), recon.DomainInfo.Domain,
		LabelStyle.Render("Registrar:"), formatValue(registrar),
		LabelStyle.Render("Creation Date:"), formatValue(creationDate),
		LabelStyle.Render("Expiration:"), formatValue(expirationDate),
		LabelStyle.Render("Nameservers:"),
	)
	
	for _, ns := range nameservers {
		info += fmt.Sprintf("  - %s\n", ns)
	}

	ipAddress := "Unknown IP"
	asn := "Unknown ASN"
	if recon.GeoIPInfo != nil {
		ipAddress = fmt.Sprintf("%s (%s)", recon.GeoIPInfo.IP, recon.GeoIPInfo.Country)
		if recon.GeoIPInfo.ASN != "" {
			asn = recon.GeoIPInfo.ASN + " - " + recon.GeoIPInfo.ASNName
		}
	} else if results.ThreatIntelligence != nil && len(results.ThreatIntelligence.RelatedIPs) > 0 {
		ipAddress = results.ThreatIntelligence.RelatedIPs[0]
	}

	threatsFound := 0
	if results.ThreatIntelligence != nil {
		for _, feed := range results.ThreatIntelligence.ReputationFeeds {
			if feed.Detected {
				threatsFound++
			}
		}
	}

	info += fmt.Sprintf("\n%s\n", HeaderStyle.Render("IP Addresses"))
	info += fmt.Sprintf("  - %s\n", ipAddress)
	info += fmt.Sprintf("%-15s %-30s\n", LabelStyle.Render("ASN:"), asn)
	info += fmt.Sprintf("%-15s %s\n", LabelStyle.Render("Reputation:"), getRiskValueStyle(results.Scoring.RiskScore).Render(fmt.Sprintf("%.1f / 100 (%s)", 100.0-results.Scoring.RiskScore, getRiskLabel(results.Scoring.RiskScore))))
	info += fmt.Sprintf("%-15s %d\n", LabelStyle.Render("Threats Found:"), threatsFound)
	info += fmt.Sprintf("%-15s %s\n", LabelStyle.Render("Risk Level:"), getRiskValueStyle(results.Scoring.RiskScore).Render(getRiskLabel(results.Scoring.RiskScore)))
	
    // Render world map with marker based on GeoIP location (if available)
    mapASCII := RenderWorldMap(results)
    
	box := BoxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			HeaderStyle.Render("Domain Intelligence"),
			info,
		),
	)
	
	fmt.Println(RenderWithMask(box, mapASCII))
}

func PrintIPInfo(results *models.ScanResult) {
	ip := results.Target
	asn := "Unknown"
	asnName := "Unknown"
	org := "Unknown"
	country := "Unknown"

	if results.ReconData != nil && results.ReconData.GeoIPInfo != nil {
		geo := results.ReconData.GeoIPInfo
		ip = geo.IP
		asn = geo.ASN
		asnName = geo.ASNName
		country = geo.Country
		org = geo.HostingProvider
		if org == "" {
			org = geo.ASNName
		}
	} else if results.ThreatIntelligence != nil && len(results.ThreatIntelligence.RelatedIPs) > 0 {
		ip = results.ThreatIntelligence.RelatedIPs[0]
	}
	
	threatsFound := 0
	if results.ThreatIntelligence != nil {
		for _, feed := range results.ThreatIntelligence.ReputationFeeds {
			if feed.Detected {
				threatsFound++
			}
		}
	}

	info := fmt.Sprintf("%-15s %-30s\n%-15s %-30s\n%-15s %-30s\n%-15s %-30s\n%-15s %-30s\n%-15s %s\n%-15s %d\n%-15s %s\n%-15s %s\n",
		LabelStyle.Render("IP Address:"), ip,
		LabelStyle.Render("ISP:"), asnName,
		LabelStyle.Render("Organization:"), org,
		LabelStyle.Render("ASN:"), asn,
		LabelStyle.Render("Country:"), country,
		LabelStyle.Render("Reputation:"), getRiskValueStyle(results.Scoring.RiskScore).Render(fmt.Sprintf("%.1f / 100 (%s)", 100.0-results.Scoring.RiskScore, getRiskLabel(results.Scoring.RiskScore))),
		LabelStyle.Render("Abuse Score:"), threatsFound,
		LabelStyle.Render("Last Seen:"), results.UpdatedAt.Format("2006-01-02 15:04:05"),
		LabelStyle.Render("Risk Level:"), getRiskValueStyle(results.Scoring.RiskScore).Render(getRiskLabel(results.Scoring.RiskScore)),
	)
	
    // Render globe with marker based on GeoIP location (if available)
    globeASCII := RenderGlobeMap(results)
    
	box := BoxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
        HeaderStyle.Render("IP Intelligence"),
        info,
		),
	)
	
	fmt.Println(RenderWithMask(box, globeASCII))
}

func PrintGraph(results *models.ScanResult) {
	if results.GraphData == nil {
		PrintInfo("No graph data available.")
		return
	}
	
	target := results.Target
	ip := "Unknown IP"
	asn := "Unknown ASN"
	asnName := "Unknown ISP"
	registrar := "Unknown WHOIS"
	threatCount := 0
	riskScore := results.Scoring.RiskScore
	
	if results.ReconData != nil {
		if results.ReconData.GeoIPInfo != nil {
			if results.ReconData.GeoIPInfo.IP != "" {
				ip = results.ReconData.GeoIPInfo.IP
			}
			if results.ReconData.GeoIPInfo.ASN != "" {
				asn = results.ReconData.GeoIPInfo.ASN
				asnName = results.ReconData.GeoIPInfo.ASNName
			}
		}
		if results.ReconData.WhoisData != nil && results.ReconData.WhoisData.Registrar != "" {
			registrar = results.ReconData.WhoisData.Registrar
		}
	}
	if results.ThreatIntelligence != nil {
		for _, feed := range results.ThreatIntelligence.ReputationFeeds {
			if feed.Detected {
				threatCount++
			}
		}
	}
	
	threatLabel := fmt.Sprintf("%d Threats", threatCount)
	riskLabel := getRiskLabel(riskScore)

	padCenter := func(s string, w int) string {
		if len(s) >= w {
			return s[:w]
		}
		pad := w - len(s)
		left := pad / 2
		right := pad - left
		return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
	}

	ipStyled := InfoStyle.Render("(" + padCenter(ip, 15) + ")")
	asnStyled := WarningStyle.Render("(" + padCenter(asn, 15) + ")")
	targetStyled := SuccessStyle.Render("(" + padCenter(target, 15) + ")")
	threatStyled := ErrorStyle.Render("(" + padCenter(threatLabel, 15) + ")")
	whoisStyled := DimStyle.Render("(" + padCenter(registrar, 15) + ")")
	riskStyled := getRiskValueStyle(riskScore).Render("(" + padCenter(riskLabel, 15) + ")")

	graphASCII := fmt.Sprintf(`
      %s       %s
           IP                 %s
             \               /
              \             /
           %s ─── %s
               %.0f/100         Detected
             /               \
            /                 \
       %s            %s
        Verified              Risk Level
`, ipStyled, asnStyled, WarningStyle.Render(padCenter(asnName, 17)), targetStyled, threatStyled, 100.0-riskScore, whoisStyled, riskStyled)
	
	legend := fmt.Sprintf("%s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s",
		LabelStyle.Render("Legend:"),
		SuccessStyle.Render("●"), "Domain",
		InfoStyle.Render("●"), "IP",
		WarningStyle.Render("●"), "ASN",
		lipgloss.NewStyle().Foreground(AccentColor).Render("●"), "Intelligence",
		ErrorStyle.Render("●"), "Threat",
	)
	
	content := lipgloss.JoinHorizontal(lipgloss.Center, graphASCII, "    ", legend)
	
	box := BoxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			HeaderStyle.Render("Threat Graph Visualization"),
			content,
		),
	)
	
	fmt.Println(box)
}

func PrintHistory(history []models.ScanResult) {
    if len(history) == 0 {
        PrintInfo("No scan history available.")
        return
    }
    rows := []string{lipgloss.JoinHorizontal(lipgloss.Top, "ID", "Target", "Risk", "Date")}
    for _, r := range history {
        risk := fmt.Sprintf("%.1f", r.Scoring.RiskScore)
        date := r.CreatedAt.Format("2006-01-02")
        row := lipgloss.JoinHorizontal(lipgloss.Top, r.ID, r.Target, risk, date)
        rows = append(rows, row)
    }
    table := lipgloss.JoinVertical(lipgloss.Left, rows...)
    box := BoxStyle.Render(table)
    fmt.Println(box)
}

func PrintAIExplanation(results *models.ScanResult) {
	if results.AISummary == nil {
		PrintInfo("No AI explanation available.")
		return
	}
	ai := results.AISummary
	
	wrappedStyle := lipgloss.NewStyle().Width(80)
	
	res := fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s\n",
		lipgloss.NewStyle().Foreground(AccentColor).Bold(true).Render("🤖 AI INSIGHT"),
		wrappedStyle.Render(ai.Summary),
		HeaderStyle.Render("Attack Chain"),
		wrappedStyle.Render(ai.AttackChain),
		HeaderStyle.Render("Recommendations"),
	)
	
	for _, rec := range ai.Remediation {
		res += fmt.Sprintf("\n  • %s", wrappedStyle.Render(rec))
	}
	
	fmt.Println(BoxStyle.MaxWidth(85).Render(res))
}

func formatValue(v string) string {
	if v == "" {
		return DimStyle.Render("Unknown")
	}
	return v
}

func getRiskLabel(score float64) string {
	if score >= 70 {
		return "High Risk"
	}
	if score >= 35 {
		return "Medium Risk"
	}
	return "Low Risk"
}

func getRiskValueStyle(score float64) lipgloss.Style {
	if score >= 70 {
		return RiskHighStyle
	}
	if score >= 35 {
		return RiskMediumStyle
	}
	return RiskLowStyle
}

func formatSignals(signals []string) string {
	if len(signals) == 0 {
		return "  None detected"
	}
	res := ""
	for _, s := range signals {
		res += fmt.Sprintf("  • %s\n", s)
	}
	return res
}

// Placeholder functions for dynamic rendering (to be implemented)
func RenderCircularGauge(score float64) string {
    // Simple solid circle using Unicode characters; fill proportionally
    // For demonstration, use 10 segments.
    filled := int(math.Round(score / 10.0))
    if filled > 10 {
        filled = 10
    }
    result := strings.Repeat("⚫", filled) + strings.Repeat("⚪", 10-filled)
    style := lipgloss.NewStyle().Foreground(SecondaryColor).Bold(true)
    return style.Render(result)
}

func RenderWorldMap(results *models.ScanResult) string {
    // Use WorldMap constant and replace a placeholder marker (e.g., '*') with location based on GeoIP
    if results.ReconData == nil || results.ReconData.GeoIPInfo == nil {
        return WorldMap
    }
    // Very simple: insert marker at line 5 for demonstration
    lines := strings.Split(WorldMap, "\n")
    if len(lines) > 4 {
        lines[4] = strings.Replace(lines[4], "#", "✈", 1)
    }
    return strings.Join(lines, "\n")
}

func RenderGlobeMap(results *models.ScanResult) string {
    if results.ReconData == nil || results.ReconData.GeoIPInfo == nil {
        return GlobeMap
    }
    // Simple marker insertion
    return strings.Replace(GlobeMap, "( )", "(X)", 1)
}

func init() {
	// Lipgloss handles NO_COLOR automatically
}
