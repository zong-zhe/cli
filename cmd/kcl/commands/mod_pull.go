package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/downloader"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/utils"
)

const (
	modPullDesc = `This command pulls kcl modules from the registry.
`
	modPullExample = `  # Pull the the module named "k8s" to the local path from the registry
  kcl mod pull k8s

  # Pull the module dependency named "k8s" with the version "1.28"
  kcl mod add k8s:1.28

  # Pull the module from the GitHub by git url
  kcl mod pull git://github.com/kcl-lang/konfig --tag v0.4.0

  # Pull the module from the OCI Registry by oci url
  kcl mod pull oci://github.com/kcl-lang/konfig --tag v0.4.0

  # Pull the module from the GitHub by flag
  kcl mod pull --git https://github.com/kcl-lang/konfig --tag v0.4.0

  # Pull the module from the OCI Registry by flag
  kcl mod pull --oci https://ghcr.io/kcl-lang/helloworld --tag 0.1.0`
)

// NewModPullCmd returns the mod pull command.
func NewModPullCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pull",
		Short:   "pull kcl package from the registry",
		Long:    modPullDesc,
		Example: modPullExample,
		RunE: func(_ *cobra.Command, args []string) error {
			// source := argsGet(args, 0)
			localPath := argsGet(args, 1)
			// return cli.PullFromOci(localPath, source, tag)
			return pull(cli, args, localPath)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&git, "git", "", "git repository url")
	cmd.Flags().StringVar(&oci, "oci", "", "oci repository url")
	cmd.Flags().StringVar(&tag, "tag", "", "git or oci repository tag")
	cmd.Flags().StringVar(&commit, "commit", "", "git repository commit")
	cmd.Flags().StringVar(&branch, "branch", "", "git repository branch")

	return cmd
}

func pull(cli *client.KpmClient, args []string, localPath string) error {
	sourceUrl, err := ParseUrlFromArgs(cli, args, localPath)
	if err != nil {
		return err
	}
	var source pkg.Source
	source.ParseFromUrl(*sourceUrl)

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	// clean the temp dir.
	defer os.RemoveAll(tmpDir)

	err = cli.DepDownloader.Download(*downloader.NewDownloadOptions(
		downloader.WithLocalPath(tmpDir),
		downloader.WithSource(source),
		downloader.WithLogWriter(cli.GetLogWriter()),
		downloader.WithSettings(*cli.GetSettings()),
	))

	if err != nil {
		return err
	}

	localPath = filepath.Join(
		localPath,
		sourceUrl.Host,
		sourceUrl.Path,
		sourceUrl.Query().Get("tag"),
		sourceUrl.Query().Get("commit"),
		sourceUrl.Query().Get("branch"),
	)

	if utils.DirExists(localPath) {
		err := os.RemoveAll(localPath)
		if err != nil {
			return err
		}
	}

	if runtime.GOOS != "windows" {
		err = os.Rename(tmpDir, localPath)
		if err != nil {
			return err
		}
	} else {
		err = copy.Copy(tmpDir, localPath)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}
