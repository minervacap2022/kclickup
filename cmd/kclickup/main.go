package main

import (
	"fmt"
	"os"

	"github.com/minervacap2022/klik-clickup-cli/internal/api"
	"github.com/minervacap2022/klik-clickup-cli/internal/auth"
	"github.com/minervacap2022/klik-clickup-cli/internal/output"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kclickup",
	Short: "ClickUp CLI for KLIK platform",
	Long:  "Command-line interface for ClickUp REST API v2. Auth via CLICKUP_API_TOKEN env var.",
}

func getClient() *api.Client {
	token, err := auth.GetToken()
	if err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
	return api.NewClient(token)
}

// --- task commands ---

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks in a list",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		listID, _ := cmd.Flags().GetString("list")

		result, err := client.Get(fmt.Sprintf("/list/%s/task", listID))
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

var taskCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		listID, _ := cmd.Flags().GetString("list")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		payload := map[string]interface{}{
			"name": name,
		}
		if description != "" {
			payload["description"] = description
		}

		result, err := client.Post(fmt.Sprintf("/list/%s/task", listID), payload)
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

var taskUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		id, _ := cmd.Flags().GetString("id")
		name, _ := cmd.Flags().GetString("name")
		status, _ := cmd.Flags().GetString("status")

		payload := map[string]interface{}{}
		if name != "" {
			payload["name"] = name
		}
		if status != "" {
			payload["status"] = status
		}

		result, err := client.Put(fmt.Sprintf("/task/%s", id), payload)
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

var taskViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		id, _ := cmd.Flags().GetString("id")

		result, err := client.Get(fmt.Sprintf("/task/%s", id))
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

// --- space commands ---

var spaceCmd = &cobra.Command{
	Use:   "space",
	Short: "Manage spaces",
}

var spaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List spaces in a team/workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		team, _ := cmd.Flags().GetString("team")
		if team == "" {
			team = auth.GetTeamID()
		}
		if team == "" {
			return fmt.Errorf("--team flag or CLICKUP_TEAM_ID env var required")
		}

		result, err := client.Get(fmt.Sprintf("/team/%s/space", team))
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

// --- list commands ---

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Manage lists",
}

var listListCmd = &cobra.Command{
	Use:   "list",
	Short: "List lists in a folder",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		folder, _ := cmd.Flags().GetString("folder")

		result, err := client.Get(fmt.Sprintf("/folder/%s/list", folder))
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

// --- comment commands ---

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage comments",
}

var commentAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a comment to a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		taskID, _ := cmd.Flags().GetString("task")
		text, _ := cmd.Flags().GetString("text")

		payload := map[string]interface{}{
			"comment_text": text,
		}

		result, err := client.Post(fmt.Sprintf("/task/%s/comment", taskID), payload)
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

func init() {
	// task
	taskListCmd.Flags().String("list", "", "List ID")
	taskListCmd.MarkFlagRequired("list")

	taskCreateCmd.Flags().String("list", "", "List ID")
	taskCreateCmd.Flags().String("name", "", "Task name")
	taskCreateCmd.Flags().String("description", "", "Task description")
	taskCreateCmd.MarkFlagRequired("list")
	taskCreateCmd.MarkFlagRequired("name")

	taskUpdateCmd.Flags().String("id", "", "Task ID")
	taskUpdateCmd.Flags().String("name", "", "New name")
	taskUpdateCmd.Flags().String("status", "", "New status")
	taskUpdateCmd.MarkFlagRequired("id")

	taskViewCmd.Flags().String("id", "", "Task ID")
	taskViewCmd.MarkFlagRequired("id")

	taskCmd.AddCommand(taskListCmd, taskCreateCmd, taskUpdateCmd, taskViewCmd)

	// space
	spaceListCmd.Flags().String("team", "", "Team/Workspace ID (or CLICKUP_TEAM_ID env)")
	spaceCmd.AddCommand(spaceListCmd)

	// list
	listListCmd.Flags().String("folder", "", "Folder ID")
	listListCmd.MarkFlagRequired("folder")
	listCmd.AddCommand(listListCmd)

	// comment
	commentAddCmd.Flags().String("task", "", "Task ID")
	commentAddCmd.Flags().String("text", "", "Comment text")
	commentAddCmd.MarkFlagRequired("task")
	commentAddCmd.MarkFlagRequired("text")
	commentCmd.AddCommand(commentAddCmd)

	rootCmd.AddCommand(taskCmd, spaceCmd, listCmd, commentCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
