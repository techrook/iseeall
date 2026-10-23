package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Node represents a file or directory in the tree
type Node struct {
	Name     string
	IsDir    bool
	Children []*Node
	Err      error
}

// Global variables to hold our flag values
var (
	showAll      bool
	maxDepth     int
	ignoreErrors bool
)


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iseeall [directory]",
	Short: "A fast, robust CLI tool to visualize directory structures as a tree.",
	Long: `iseeall is a production-grade directory tree visualizer built for speed and reliability. 

Unlike standard tree commands, iseeall is designed with the real-world filesystem in mind:
  • Gracefully handles permission errors without crashing (warnings go to stderr).
  • Hidden files are hidden by default to reduce clutter (use -a to show them).
  • Prevents runaway recursion in massive directories with a configurable max depth (-L).
  • Follows Unix philosophy: clean tree output to stdout, errors to stderr.

Examples:
  iseeall                  # Visualize the current directory
  iseeall /var/log -L 2    # Visualize /var/log up to 2 levels deep
  iseeall -a --ignore-errors ./my-project # Show all files, silently skip unreadable dirs`,

  	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Determine target directory (default to current directory ".")
		targetDir := "."
		if len(args) > 0 {
			targetDir = args[0]
		}


		// 2. Build the tree in memory
		// We start at depth 0. (Note: we will handle maxDepth == 0 as "infinite" inside buildTree)
		rootNode, err := buildTree(targetDir, 0, maxDepth)
		if err != nil {
			return err
		}

		// 3. Print the tree to stdout
		printTree(rootNode, "", true)

		return nil
	},

}

func buildTree(dirPath string, currentDepth int, maxDepth int) (*Node, error) {
	entries, err := os.ReadDir(dirPath)
	
	// Handle permission errors gracefully
	if err != nil {
		return &Node{
			Name:  filepath.Base(dirPath),
			IsDir: true,
			Err:   err,
		}, nil 
	}

	currentNode := &Node{
		Name:     filepath.Base(dirPath),
		IsDir:    true,
		Children: []*Node{},
	}

	// Check max depth (0 means infinite)
	if maxDepth > 0 && currentDepth >= maxDepth {
		return currentNode, nil
	}

	for _, entry := range entries {
		// Skip hidden files unless showAll is true
		if !showAll && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		childPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			childNode, _ := buildTree(childPath, currentDepth+1, maxDepth)
			currentNode.Children = append(currentNode.Children, childNode)
		} else {
			currentNode.Children = append(currentNode.Children, &Node{
				Name:  entry.Name(),
				IsDir: false,
			})
		}
	}

	return currentNode, nil
}

func printTree(node *Node, prefix string, isLast bool) {
	// 1. Handle errors on the node (print to stderr, not stdout!)
	if node.Err != nil {
		if !ignoreErrors {
			fmt.Fprintf(os.Stderr, "Warning: %s: %v\n", node.Name, node.Err)
		}
		return
	}

	// 2. Determine the connector
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	// 3. Print the current node (Root node has no prefix)
	if prefix == "" {
		fmt.Println(node.Name)
	} else {
		fmt.Println(prefix + connector + node.Name)
	}

	// 4. Build the prefix for the children
	childPrefix := prefix
	if prefix != "" {
		if isLast {
			childPrefix += "    " // 4 spaces (no vertical line needed below)
		} else {
			childPrefix += "│   " // Vertical line continues down
		}
	}

	// 5. Recursively print all children
	for i, child := range node.Children {
		isLastChild := (i == len(node.Children)-1)
		printTree(child, childPrefix, isLastChild)
	}
}
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&showAll, "all", "a", false, "include hidden files and directories")
	rootCmd.Flags().IntVarP(&maxDepth, "level", "L", 0, "max display depth of the directory tree (0 for infinite)")
	rootCmd.Flags().BoolVar(&ignoreErrors, "ignore-errors", false, "do not print permission errors to stderr")
}


