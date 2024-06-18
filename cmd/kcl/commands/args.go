package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/constants"
	"kcl-lang.io/kpm/pkg/opt"
	"oras.land/oras-go/v2/registry"
)

func argsGet(a []string, n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
}

func ParseUrlFromArgs(cli *client.KpmClient, args []string, localPath string) (*url.URL, error) {
	var sourceUrl url.URL

	if len(args) == 0 {
		if len(git) != 0 {
			gitUrl, err := url.Parse(git)
			if err != nil {
				return nil, err
			}

			gitUrl.Scheme = constants.GitScheme
			query := gitUrl.Query()
			query.Add(constants.Tag, tag)
			query.Add(constants.GitCommit, commit)
			query.Add(constants.GitBranch, branch)
			gitUrl.RawQuery = query.Encode()
			sourceUrl = *gitUrl
		} else if len(oci) != 0 {
			ociUrl, err := url.Parse(oci)
			if err != nil {
				return nil, err
			}

			ociUrl.Scheme = constants.OciScheme
			query := ociUrl.Query()
			query.Add(constants.Tag, tag)
			ociUrl.RawQuery = query.Encode()
			sourceUrl = *ociUrl
		}
	} else {
		url, err := url.Parse(args[0])
		if err != nil {
			return nil, err
		}
		query := url.Query()

		// It is a oci ref
		if len(url.Scheme) == 0 {
			url.Scheme = constants.OciScheme
			var repo string
			var reg string
			ref, err := registry.ParseReference(args[0])
			if err != nil {
				var pkgName string
				pkgName, tag, err = opt.ParseOciPkgNameAndVersion(args[0])
				if err != nil {
					return nil, err
				}
				reg = cli.GetSettings().DefaultOciRegistry()
				repo = cli.GetSettings().DefaultOciRepo()
				if !strings.HasPrefix(pkgName, "/") {
					repo = fmt.Sprintf("%s/%s", repo, pkgName)
				} else {
					repo = fmt.Sprintf("%s%s", repo, pkgName)
				}

			} else {
				reg = ref.Registry
				repo = ref.Repository
				tag = ref.ReferenceOrDefault()
			}
			url.Host = reg
			url.Path = repo
			query.Add(constants.Tag, tag)
		} else {
			query.Add(constants.Tag, tag)
			query.Add(constants.GitCommit, commit)
			query.Add(constants.GitBranch, branch)
		}
		url.RawQuery = query.Encode()
		sourceUrl = *url
	}
	return &sourceUrl, nil
}
