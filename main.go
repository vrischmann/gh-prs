package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// pullRequest represents a simplified GitHub Pull Request structure for JSON unmarshaling.
type pullRequest struct {
	Number     int        `json:"number"`
	Repository repository `json:"repository"`
	Title      string     `json:"title"`
	URL        string     `json:"url"`
}

// repository represents the repository information within a PR.
type repository struct {
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
}

type fetchType int

const (
	reviewedFetchType fetchType = iota + 1
	toReviewFetchType
)

// fetchPRs fetches Pull Requests from GitHub using the 'gh search prs' command.
// It constructs the search query and parses the JSON output.
func fetchPRs(typ fetchType, org string) ([]pullRequest, error) {
	args := []string{
		"--json", "url,repository,number,title",
		"search", "prs", fmt.Sprintf("org:%s", org), "is:open",
	}
	switch typ {
	case reviewedFetchType:
		args = append(args, "reviewed-by:@me")
	case toReviewFetchType:
		args = append(args, "review-requested:@me")
	}

	cmd := exec.Command("gh", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		stderrStr := stderr.String()
		if strings.Contains(stderrStr, "No pull requests match your search") {
			return nil, nil
		}
		return nil, fmt.Errorf("gh search failed: %w\nStderr: %s", err, stderrStr)
	}

	if stdout.Len() == 0 {
		return nil, nil
	}

	var prs []pullRequest
	if err := json.Unmarshal(stdout.Bytes(), &prs); err != nil {
		return nil, fmt.Errorf("failed to parse gh output JSON: %w\nOutput: %s", err, stdout.String())
	}

	return prs, nil
}

// openBrowser uses the 'gh browse' command to open a URL in the default browser.
func openBrowser(url string) error {
	cmd := exec.Command("open", url)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open URL in browser: %w\nStderr: %s", err, stderr.String())
	}
	return nil
}

//
// UI
//

// model represents the state of our interactive UI.
type model struct {
	prs         []pullRequest
	cursor      int
	selectedPR  string // Stores the URL of the selected PR
	description string // "to review" or "reviewed"
	quitting    bool
	err         error
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model's state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.prs)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor >= 0 && m.cursor < len(m.prs) {
				m.selectedPR = m.prs[m.cursor].URL
			}
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the UI.
func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if m.quitting {
		return "" // Don't render anything if quitting
	}

	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Select a PR %s:\n\n", m.description))

	for i, pr := range m.prs {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render a single line for each PR
		prLine := fmt.Sprintf("%s [%s] #%d - %s", cursor, pr.Repository.Name, pr.Number, pr.Title)

		if m.cursor == i {
			// Highlight the selected row
			s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(prLine))
		} else {
			s.WriteString(prLine)
		}
		s.WriteString("\n")
	}

	s.WriteString("\n(Press 'q' or 'esc' to quit, 'enter' to select)")
	return s.String()
}

// startInteractiveUI runs the Bubble Tea program and returns the selected PR URL.
func startInteractiveUI(prs []pullRequest, description string) (string, error) {
	p := tea.NewProgram(model{prs: prs, description: description}, tea.WithOutput(os.Stderr)) // Use stderr for UI to keep stdout clean for potential piping if needed

	m, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run interactive UI: %w", err)
	}

	if m, ok := m.(model); ok && m.selectedPR != "" {
		return m.selectedPR, nil
	}

	return "", nil // No PR selected or program quit
}

var (
	version = "dev"
	orgFlag string
)

var rootCmd = &cobra.Command{
	Use:   "gh prs",
	Short: "Interactively list and open GitHub Pull Requests",
	Long: `A GitHub CLI extension that fetches pull requests based on review status 
and provides an interactive selection interface to open them in your browser.

Requires GitHub CLI (gh) to be installed and authenticated.`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var toReviewCmd = &cobra.Command{
	Use:   "to-review",
	Short: "List open PRs with review requests for you",
	Long: `Fetches and interactively lists open Pull Requests in the specified
organization where you are a requested reviewer and have not yet submitted a review.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkGHInstalled(); err != nil {
			return err
		}

		org, err := getOrganization()
		if err != nil {
			return err
		}

		prs, err := fetchPRs(toReviewFetchType, org)
		if err != nil {
			return fmt.Errorf("failed to fetch 'to review' PRs: %w", err)
		}

		if len(prs) == 0 {
			fmt.Printf("No open Pull Requests found in '%s' where you are a requested reviewer.\n", org)
			return nil
		}

		selectedURL, err := startInteractiveUI(prs, "to review") // Call StartInteractiveUI directly
		if err != nil {
			return fmt.Errorf("interactive selection failed: %w", err)
		}

		if selectedURL != "" {
			fmt.Printf("Opening %s in browser...\n", selectedURL)
			if err := openBrowser(selectedURL); err != nil { // Call OpenBrowser directly
				return fmt.Errorf("failed to open browser: %w", err)
			}
		} else {
			fmt.Println("No PR selected.")
		}

		return nil
	},
}

// --- Cmd/reviewed.go content ---

var reviewedCmd = &cobra.Command{
	Use:   "reviewed",
	Short: "List open PRs you have already reviewed",
	Long: `Fetches and interactively lists open Pull Requests in the specified
organization that you have already submitted a review for.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkGHInstalled(); err != nil {
			return err
		}

		org, err := getOrganization()
		if err != nil {
			return err
		}

		prs, err := fetchPRs(reviewedFetchType, org)
		if err != nil {
			return fmt.Errorf("failed to fetch 'reviewed' PRs: %w", err)
		}

		if len(prs) == 0 {
			fmt.Printf("No open Pull Requests found in '%s' that you have reviewed.\n", org)
			return nil
		}

		selectedURL, err := startInteractiveUI(prs, "reviewed") // Call StartInteractiveUI directly
		if err != nil {
			return fmt.Errorf("interactive selection failed: %w", err)
		}

		if selectedURL != "" {
			fmt.Printf("Opening %s in browser...\n", selectedURL)
			if err := openBrowser(selectedURL); err != nil { // Call OpenBrowser directly
				return fmt.Errorf("failed to open browser: %w", err)
			}
		} else {
			fmt.Println("No PR selected.")
		}

		return nil
	},
}

// --- main.go content (modified init to add commands) ---

// checkGHInstalled verifies that GitHub CLI is installed and authenticated.
func checkGHInstalled() error {
	cmd := exec.Command("gh", "auth", "status")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("GitHub CLI is not installed or not authenticated. Please install gh and run 'gh auth login': %w\nStderr: %s", err, stderr.String())
	}
	return nil
}

// getOrganization determines the organization to use for PR searches.
// It uses the --org flag if provided, otherwise attempts to detect from current repo.
func getOrganization() (string, error) {
	if orgFlag != "" {
		return orgFlag, nil
	}

	// Try to get organization from current repository
	cmd := exec.Command("gh", "repo", "view", "--json", "owner")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		stderrStr := stderr.String()
		if strings.Contains(stderrStr, "not a git repository") || strings.Contains(stderrStr, "not found") {
			return "", fmt.Errorf("not in a GitHub repository. Please use --org flag to specify organization")
		}
		return "", fmt.Errorf("unable to determine organization: %w\nStderr: %s", err, stderrStr)
	}

	var repo struct {
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	}

	if err := json.Unmarshal(stdout.Bytes(), &repo); err != nil {
		return "", fmt.Errorf("failed to parse repository information: %w", err)
	}

	if repo.Owner.Login == "" {
		return "", fmt.Errorf("unable to determine organization from current repository")
	}

	return repo.Owner.Login, nil
}

func init() {
	// Add persistent flags
	rootCmd.PersistentFlags().StringVarP(&orgFlag, "org", "o", "", "GitHub organization to search (defaults to current repo's organization)")
	
	// Add commands to the root command
	rootCmd.AddCommand(toReviewCmd)
	rootCmd.AddCommand(reviewedCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
