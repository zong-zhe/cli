package cmd

import (
	"net/url"

	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/constants"
	"kcl-lang.io/kpm/pkg/downloader"
	"kcl-lang.io/kpm/pkg/opt"
	"kcl-lang.io/kpm/pkg/utils"
)

func argsGet(a []string, n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
}

func ParseUrlFromArgs(cli *client.KpmClient, args []string) (*url.URL, error) {
	var sourceUrl url.URL

	// Parse the source url from the args
	if len(git) != 0 {
		gitUrl, err := url.Parse(git)
		if err != nil {
			return nil, err
		}

		gitUrl.Scheme = constants.GitScheme
		query := gitUrl.Query()
		if tag != "" {
			query.Add(constants.Tag, tag)
		}
		if commit != "" {
			query.Add(constants.GitCommit, commit)
		}
		if branch != "" {
			query.Add(constants.GitBranch, branch)
		}
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
	} else if len(path) != 0 {
		pathUrl, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		sourceUrl = *pathUrl
	}

	// Parse the mod spec from the args
	if len(args) != 0 {
		source, err := downloader.NewSourceFromStr(sourceUrl.String())
		if err != nil {
			return nil, err
		}
		source.ModSpec = &downloader.ModSpec{}
		err = source.ModSpec.FromString(args[0])
		if err != nil {
			url, err := url.Parse(args[0])
			if err != nil {
				return nil, err
			}
			query := url.Query()
			url.Opaque = ""
			regOpts, err := opt.NewRegistryOptionsFrom(args[0], cli.GetSettings())
			if err != nil {
				return nil, err
			}

			if regOpts.Git != nil {
				if url.Scheme != constants.GitScheme && url.Scheme != constants.SshScheme {
					url.Scheme = constants.GitScheme
				}
				if tag != "" {
					query.Add(constants.Tag, tag)
				}
				if commit != "" {
					query.Add(constants.GitCommit, commit)
				}
				if branch != "" {
					query.Add(constants.GitBranch, branch)
				}
			} else if regOpts.Oci != nil {
				url.Scheme = constants.OciScheme
				url.Host = regOpts.Oci.Reg
				url.Path = regOpts.Oci.Repo
				if regOpts.Oci.Tag != "" {
					query.Add(constants.Tag, regOpts.Oci.Tag)
				}
				if tag != "" {
					query.Add(constants.Tag, tag)
				}
			} else if regOpts.Registry != nil {
				url.Scheme = constants.DefaultOciScheme
				url.Host = regOpts.Registry.Reg
				url.Path = regOpts.Registry.Repo
				if regOpts.Registry.Tag != "" {
					query.Add(constants.Tag, regOpts.Registry.Tag)
				}
				if tag != "" {
					query.Add(constants.Tag, tag)
				}
			}

			url.RawQuery = query.Encode()
			sourceUrl = *url
		} else {
			urlStr, err := source.ToString()
			if err != nil {
				return nil, err
			}

			urlWithSpec, err := url.Parse(urlStr)
			if err != nil {
				return nil, err
			}
			urlWithSpec.Scheme = sourceUrl.Scheme
			sourceUrl = *urlWithSpec
		}
	}

	source, err := downloader.NewSourceFromStr(sourceUrl.String())
	if err != nil {
		return nil, err
	}

	if source.SpecOnly() {
		source.Oci = &downloader.Oci{
			Reg:  cli.GetSettings().DefaultOciRegistry(),
			Repo: utils.JoinPath(cli.GetSettings().DefaultOciRepo(), source.ModSpec.Name),
			Tag:  source.ModSpec.Version,
		}
	}

	sourceUrlStr, err := source.ToString()
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(sourceUrlStr)
	if err != nil {
		return nil, err
	}

	return u, nil
}
