// Copyright © 2017 blacktop <https://github.com/blacktop>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/cli"
	"github.com/maliceio/malice/cli/command"
	"github.com/maliceio/malice/cli/command/commands"
	cliconfig "github.com/moby/moby/cli/config"
	cliflags "github.com/moby/moby/cli/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	version string = "dev"
	commit  string = "dev"
	date    string = "dev"
)

func newMaliceCommand(maliceCli *command.MaliceCli) *cobra.Command {
	opts := cliflags.NewClientOptions()
	var flags *pflag.FlagSet

	cmd := &cobra.Command{
		Use:           "malice [OPTIONS] COMMAND [ARG...]",
		Short:         "Open Source Malware Analysis Framework",
		SilenceUsage:  true,
		SilenceErrors: true,
		// TraverseChildren: true,
		// Args:             noArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Version {
				showVersion()
				return nil
			}
			return maliceCli.ShowHelp(cmd, args)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// flags must be the top-level command flags, not cmd.Flags()
			// opts.Common.SetDefaultOptions(flags)
			malicePreRun(opts)
			if err := maliceCli.Initialize(opts); err != nil {
				return err
			}
			return nil
			// return isSupported(cmd, maliceCli.Client().ClientVersion(), maliceCli.OSType(), maliceCli.HasExperimental())
		},
	}
	cli.SetupRootCommand(cmd)

	flags = cmd.Flags()
	flags.BoolVarP(&opts.Version, "version", "v", false, "Print version information and quit")
	// flags.StringVar(&opts.ConfigDir, "config", cliconfig.Dir(), "Location of client config files")
	// opts.Common.InstallFlags(flags)

	// setFlagErrorFunc(maliceCli, cmd, flags, opts)

	// setHelpFunc(maliceCli, cmd, flags, opts)

	cmd.SetOutput(maliceCli.Out())
	// cmd.AddCommand(newDaemonCommand())
	commands.AddCommands(cmd, maliceCli)

	setValidateArgs(maliceCli, cmd, flags, opts)

	return cmd
}

func initializeMaliceCli(maliceCli *command.MaliceCli, flags *pflag.FlagSet, opts *cliflags.ClientOptions) {
	if maliceCli.Client() == nil { // when using --help, PersistentPreRun is not called, so initialization is needed.
		// flags must be the top-level command flags, not cmd.Flags()
		// opts.Common.SetDefaultOptions(flags)
		malicePreRun(opts)
		maliceCli.Initialize(opts)
	}
}

func main() {
	stdin, stdout, stderr := StdStreams()
	logrus.SetOutput(os.Stderr)

	maliceCli := command.NewMaliceCli(stdin, stdout, stderr)
	cmd := newMaliceCommand(maliceCli)

	if err := cmd.Execute(); err != nil {
		if sterr, ok := err.(cli.StatusError); ok {
			if sterr.Status != "" {
				fmt.Fprintln(os.Stderr, 1)
			}
			// StatusError should only be used for errors, and all errors should
			// have a non-zero exit status, so never exit with 0
			if sterr.StatusCode == 0 {
				os.Exit(1)
			}
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func showVersion() {
	fmt.Printf("Malice version %s, build %s, date %s\n", version, commit, date)
}

func malicePreRun(opts *cliflags.ClientOptions) {
	cliflags.SetLogLevel(opts.Common.LogLevel)

	if opts.ConfigDir != "" {
		cliconfig.SetDir(opts.ConfigDir)
	}

	if opts.Common.Debug {
		// debug.Enable()
	}
}