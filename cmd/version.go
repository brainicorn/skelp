package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	logo = `       @@               @@.
       @@               @@.
 .@@+  @@ .++  ,@@@@@   @@. ++;@,
'@':@@ @@ @@  @@@@@ +@# @@. @@@@@
 @@@;  @@@@.  @@@@@   @ @@. @@ :@#
 '#@@@ @@@@@ '@@@@@  ;@ @@. @@ ,@#
#@+ @@ @@ #@' @@@@@##@@ @@. @@@@@,
 @@@@  @@  @@ '@@@@ @@@ @@. @@@@#
                +@@' '.     @@
                   @#@@@    @@
                    @@'
`
)

var (
	GitCommit = "work in progress"
	GitBranch = "master"
	Version   = "1.0-beta1"
)

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Displays skelp version info",
		Long:  `Displays skelp version info`,
		RunE:  executeVersion,
	}

	return cmd
}

func executeVersion(cmd *cobra.Command, args []string) error {
	fmt.Println(logo)
	fmt.Println("Version: ", Version)

	return nil
}
