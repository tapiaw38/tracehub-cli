package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tapiaw38/tracehub-cli/internal/client"
	"github.com/tapiaw38/tracehub-cli/internal/config"
)

var (
	cfg       *config.Config
	apiClient *client.Client
)

var rootCmd = &cobra.Command{
	Use:   "tracehub",
	Short: "TraceHub CLI - Monitor and analyze traces from your terminal",
	Long: `TraceHub CLI is a command-line interface for interacting with TraceHub server.
Connect to your projects, stream traces in real-time, and analyze errors with AI.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load configuration
		var err error
		cfg, err = config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load config: %v\n", err)
			cfg = &config.Config{}
		}

		// Initialize client if configured
		if cfg.ServerURL != "" && cfg.APIKey != "" {
			apiClient = client.New(cfg.ServerURL, cfg.APIKey)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(streamCmd)
	rootCmd.AddCommand(errorsCmd)
	rootCmd.AddCommand(versionCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure TraceHub server connection",
	Long:  "Set up the TraceHub server URL and API key for authentication.",
	Run: func(cmd *cobra.Command, args []string) {
		serverURL, _ := cmd.Flags().GetString("server")
		apiKey, _ := cmd.Flags().GetString("api-key")

		if serverURL == "" || apiKey == "" {
			fmt.Println("Error: Both --server and --api-key are required")
			os.Exit(1)
		}

		cfg.ServerURL = serverURL
		cfg.APIKey = apiKey

		if err := config.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save config: %v\n", err)
			os.Exit(1)
		}

		// Test connection
		testClient := client.New(serverURL, apiKey)
		if healthy, err := testClient.Health(); err != nil || !healthy {
			fmt.Printf("⚠️  Warning: Could not connect to server: %v\n", err)
		} else {
			fmt.Println("✓ Successfully connected to TraceHub server")
		}

		fmt.Printf("✓ Configuration saved to %s\n", config.GetConfigPath())
	},
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List available projects",
	Long:  "Display all projects available on the TraceHub server.",
	Run: func(cmd *cobra.Command, args []string) {
		if apiClient == nil {
			fmt.Println("Error: Not configured. Run 'tracehub configure' first.")
			os.Exit(1)
		}

		response, err := apiClient.ListProjects(50, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to list projects: %v\n", err)
			os.Exit(1)
		}

		if len(response.Projects) == 0 {
			fmt.Println("No projects found.")
			return
		}

		fmt.Printf("Found %d projects:\n\n", response.Total)
		for _, proj := range response.Projects {
			fmt.Printf("  • %s\n", proj.Name)
			fmt.Printf("    ID: %s\n", proj.ID)
			fmt.Printf("    Language: %s\n", proj.Language)
			if proj.Description != "" {
				fmt.Printf("    Description: %s\n", proj.Description)
			}
			fmt.Println()
		}
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect [project-id]",
	Short: "Connect to a project",
	Long:  "Set the current project for streaming traces and viewing errors.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if apiClient == nil {
			fmt.Println("Error: Not configured. Run 'tracehub configure' first.")
			os.Exit(1)
		}

		projectID := args[0]

		// Verify project exists
		project, err := apiClient.GetProject(projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get project: %v\n", err)
			os.Exit(1)
		}

		// Save current project
		cfg.CurrentProject = &config.ProjectConfig{
			ID:   project.ID,
			Name: project.Name,
		}

		if err := config.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Connected to project: %s\n", project.Name)
		fmt.Println("\nAvailable commands:")
		fmt.Println("  tracehub stream    - Stream traces in real-time")
		fmt.Println("  tracehub errors    - View detected errors")
	},
}

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream traces in real-time",
	Long:  "Display traces from the current project as they arrive.",
	Run: func(cmd *cobra.Command, args []string) {
		if apiClient == nil {
			fmt.Println("Error: Not configured. Run 'tracehub configure' first.")
			os.Exit(1)
		}

		if cfg.CurrentProject == nil {
			fmt.Println("Error: No project connected. Run 'tracehub connect <project-id>' first.")
			os.Exit(1)
		}

		level, _ := cmd.Flags().GetString("level")
		limit, _ := cmd.Flags().GetInt("limit")

		fmt.Printf("Streaming traces from project: %s\n", cfg.CurrentProject.Name)
		fmt.Println("Press Ctrl+C to stop\n")

		filters := make(map[string]string)
		if level != "" {
			filters["level"] = level
		}
		if limit > 0 {
			filters["limit"] = fmt.Sprintf("%d", limit)
		}

		// Query traces
		response, err := apiClient.QueryTraces(cfg.CurrentProject.ID, filters)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to query traces: %v\n", err)
			os.Exit(1)
		}

		if len(response.Traces) == 0 {
			fmt.Println("No traces found.")
			return
		}

		// Display traces
		for _, trace := range response.Traces {
			levelColor := getLevelColor(trace.Level)
			fmt.Printf("[%s] %s%-7s%s %s\n",
				trace.Timestamp.Format("15:04:05"),
				levelColor,
				trace.Level,
				"\033[0m",
				trace.Message,
			)
			if trace.Source != "" {
				fmt.Printf("          Source: %s\n", trace.Source)
			}
			if trace.StackTrace != "" {
				fmt.Printf("          Stack: %s\n", truncate(trace.StackTrace, 100))
			}
		}

		fmt.Printf("\nShowing %d traces (use --limit to adjust)\n", len(response.Traces))
	},
}

var errorsCmd = &cobra.Command{
	Use:   "errors",
	Short: "List detected errors",
	Long:  "Display errors detected in the current project.",
	Run: func(cmd *cobra.Command, args []string) {
		if apiClient == nil {
			fmt.Println("Error: Not configured. Run 'tracehub configure' first.")
			os.Exit(1)
		}

		if cfg.CurrentProject == nil {
			fmt.Println("Error: No project connected. Run 'tracehub connect <project-id>' first.")
			os.Exit(1)
		}

		// Query error-level traces
		filters := map[string]string{
			"level": "error",
			"limit": "20",
		}

		response, err := apiClient.QueryTraces(cfg.CurrentProject.ID, filters)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to query errors: %v\n", err)
			os.Exit(1)
		}

		if len(response.Traces) == 0 {
			fmt.Println("No errors found. ✓")
			return
		}

		fmt.Printf("Found %d recent errors:\n\n", len(response.Traces))
		for _, trace := range response.Traces {
			fmt.Printf("  ❌ %s\n", trace.Message)
			fmt.Printf("     Time: %s\n", trace.Timestamp.Format("2006-01-02 15:04:05"))
			if trace.Source != "" {
				fmt.Printf("     Source: %s\n", trace.Source)
			}
			fmt.Println()
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TraceHub CLI v1.0.0")
	},
}

func init() {
	configureCmd.Flags().String("server", "", "TraceHub server URL (e.g., http://localhost:8080)")
	configureCmd.Flags().String("api-key", "", "API key for authentication")
	configureCmd.MarkFlagRequired("server")
	configureCmd.MarkFlagRequired("api-key")

	streamCmd.Flags().String("level", "", "Filter by level (info, warn, error, fatal)")
	streamCmd.Flags().Int("limit", 100, "Number of traces to display")
}

func getLevelColor(level string) string {
	switch level {
	case "error", "fatal":
		return "\033[31m" // Red
	case "warn":
		return "\033[33m" // Yellow
	case "info":
		return "\033[32m" // Green
	case "debug":
		return "\033[36m" // Cyan
	default:
		return ""
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
