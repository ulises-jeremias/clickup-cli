package doc

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/triptechtravel/clickup-cli/internal/apiv3"
	"github.com/triptechtravel/clickup-cli/internal/tableprinter"
	"github.com/triptechtravel/clickup-cli/pkg/cmdutil"
)

type listOptions struct {
	creator    int
	deleted    bool
	archived   bool
	parentID   string
	parentType string
	limit      int
	cursor     string
	jsonFlags  cmdutil.JSONFlags
}

// NewCmdList returns a command to list ClickUp Docs in the workspace.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List ClickUp Docs in the workspace",
		Long: `List Docs in the configured ClickUp workspace.

Supports filtering by creator, status, parent location, and pagination.`,
		Example: `  # List all Docs
  clickup doc list

  # List non-deleted, non-archived Docs in JSON
  clickup doc list --json

  # List Docs in a specific space
  clickup doc list --parent-id 123456 --parent-type SPACE

  # Paginate
  clickup doc list --limit 10 --cursor <cursor>`,
		PersistentPreRunE: cmdutil.NeedsAuth(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.parentType != "" {
				if _, err := parseParentType(opts.parentType); err != nil {
					return err
				}
			}
			if err := validateParentIDAndType(opts.parentID, opts.parentType); err != nil {
				return err
			}
			return runList(f, opts)
		},
	}

	cmd.Flags().IntVar(&opts.creator, "creator", 0, "Filter by creator user ID")
	cmd.Flags().BoolVar(&opts.deleted, "deleted", false, "Include deleted Docs")
	cmd.Flags().BoolVar(&opts.archived, "archived", false, "Include archived Docs")
	cmd.Flags().StringVar(&opts.parentID, "parent-id", "", "Filter by parent ID")
	cmd.Flags().StringVar(&opts.parentType, "parent-type", "", "Parent type (SPACE|FOLDER|LIST|WORKSPACE|EVERYTHING)")
	cmd.Flags().IntVar(&opts.limit, "limit", 0, "Maximum number of Docs to return")
	cmd.Flags().StringVar(&opts.cursor, "cursor", "", "Pagination cursor from a previous response")

	cmdutil.AddJSONFlags(cmd, &opts.jsonFlags)

	return cmd
}

func runList(f *cmdutil.Factory, opts *listOptions) error {
	ios := f.IOStreams
	cs := ios.ColorScheme()

	workspaceID, err := resolveWorkspaceID(f)
	if err != nil {
		return err
	}

	client, err := f.ApiClient()
	if err != nil {
		return err
	}

	params := apiv3.SearchDocsPublicParams{
		Deleted:  opts.deleted,
		Archived: opts.archived,
		Creator:  opts.creator,
		ParentID: opts.parentID,
		Limit:    opts.limit,
		Cursor:   opts.cursor,
	}
	if opts.parentID != "" && opts.parentType != "" {
		pt, err := parseParentType(opts.parentType)
		if err != nil {
			return err
		}
		params.ParentType = pt
	}

	ctx := context.Background()
	result, err := apiv3.SearchDocsPublic(ctx, client, workspaceID, params)
	if err != nil {
		return fmt.Errorf("failed to list docs: %w", err)
	}

	if opts.jsonFlags.WantsJSON() {
		return opts.jsonFlags.OutputJSON(ios.Out, result)
	}

	if len(result.Docs) == 0 {
		fmt.Fprintln(ios.Out, cs.Gray("No Docs found."))
		return nil
	}

	tp := tableprinter.New(ios)
	tp.AddField(cs.Bold("ID"))
	tp.AddField(cs.Bold("NAME"))
	tp.AddField(cs.Bold("VISIBILITY"))
	tp.AddField(cs.Bold("STATUS"))
	tp.EndRow()
	tp.SetTruncateColumn(1)

	for _, d := range result.Docs {
		status := "active"
		if d.Deleted {
			status = cs.Red("deleted")
		} else if d.Archived {
			status = cs.Gray("archived")
		}
		tp.AddField(d.ID)
		tp.AddField(d.Name)
		tp.AddField(strings.ToLower(d.Visibility))
		tp.AddField(status)
		tp.EndRow()
	}
	if err := tp.Render(); err != nil {
		return err
	}

	if result.NextCursor != "" {
		fmt.Fprintf(ios.Out, "\n%s  clickup doc list --cursor %s\n", cs.Gray("Next page:"), result.NextCursor)
	}

	return nil
}
