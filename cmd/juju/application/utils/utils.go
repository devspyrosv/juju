// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package utils

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/mattn/go-isatty"

	"github.com/juju/charm/v10"
	charmresource "github.com/juju/charm/v10/resource"
	"github.com/juju/cmd/v3"
	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"github.com/juju/loggo"

	app "github.com/juju/juju/apiserver/facades/client/application"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/resources"
)

var logger = loggo.GetLogger("juju.cmd.juju.application.utils")

// GetMetaResources retrieves metadata resources for the given
// charm.URL.
func GetMetaResources(charmURL *charm.URL, client CharmClient) (map[string]charmresource.Meta, error) {
	charmInfo, err := client.CharmInfo(charmURL.String())
	if err != nil {
		return nil, errors.Trace(err)
	}
	return charmInfo.Meta.Resources, nil
}

// ParsePlacement validates provided placement of a unit and
// returns instance.Placement.
func ParsePlacement(spec string) (*instance.Placement, error) {
	if spec == "" {
		return nil, nil
	}
	placement, err := instance.ParsePlacement(spec)
	if err == instance.ErrPlacementScopeMissing {
		spec = fmt.Sprintf("model-uuid:%s", spec)
		placement, err = instance.ParsePlacement(spec)
	}
	if err != nil {
		return nil, errors.Errorf("invalid --to parameter %q", spec)
	}
	return placement, nil
}

// GetFlags returns the flags with the given names. Only flags that are set and
// whose name is included in flagNames are included.
func GetFlags(flagSet *gnuflag.FlagSet, flagNames []string) []string {
	flags := make([]string, 0, flagSet.NFlag())
	flagSet.Visit(func(flag *gnuflag.Flag) {
		for _, name := range flagNames {
			if flag.Name == name {
				flags = append(flags, flagWithMinus(name))
			}
		}
	})
	return flags
}

func flagWithMinus(name string) string {
	if len(name) > 1 {
		return "--" + name
	}
	return "-" + name
}

// GetUpgradeResources
func GetUpgradeResources(
	newCharmURL *charm.URL,
	resourceLister ResourceLister,
	applicationID string,
	cliResources map[string]string,
	meta map[string]charmresource.Meta,
) (map[string]charmresource.Meta, error) {
	if len(meta) == 0 {
		return nil, nil
	}
	current, err := getResources(applicationID, resourceLister)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return filterResources(newCharmURL, meta, current, cliResources)
}

func getResources(
	applicationID string,
	resourceLister ResourceLister,
) (map[string]resources.Resource, error) {
	svcs, err := resourceLister.ListResources([]string{applicationID})
	if err != nil {
		return nil, errors.Trace(err)
	}
	return resources.AsMap(svcs[0].Resources), nil
}

func filterResources(
	newCharmURL *charm.URL,
	meta map[string]charmresource.Meta,
	current map[string]resources.Resource,
	uploads map[string]string,
) (map[string]charmresource.Meta, error) {
	filtered := make(map[string]charmresource.Meta)
	for name, res := range meta {
		doUpgrade, err := shouldUpgradeResource(newCharmURL, res, uploads, current)
		if err != nil {
			return nil, errors.Trace(err)
		}
		if doUpgrade {
			filtered[name] = res
		}
	}
	return filtered, nil
}

// shouldUpgradeResource reports whether we should upload the metadata for the given
// resource.  This is always true for resources we're adding with the --resource
// flag. For resources we're not adding with --resource, we only upload metadata
// for charmstore resources.  Previously uploaded resources stay pinned to the
// data the user uploaded.
func shouldUpgradeResource(newCharmURL *charm.URL, res charmresource.Meta, uploads map[string]string, current map[string]resources.Resource) (bool, error) {
	// Always upload metadata for resources the user is uploading during
	// upgrade-charm.
	if _, ok := uploads[res.Name]; ok {
		logger.Tracef("%q provided to upgrade existing resource", res.Name)
		return true, nil
	}

	cur, ok := current[res.Name]
	if !ok {
		// If there's no information on the server, there might be a new resource added to the charm.
		if newCharmURL.Schema == charm.Local.String() {
			return false, errors.NewNotValid(nil, fmt.Sprintf("new resource %q was missing, please provide it via --resource", res.Name))
		}

		logger.Tracef("resource %q does not exist in controller, so it will be uploaded", res.Name)
		return true, nil
	}
	if newCharmURL.Schema == charm.Local.String() {
		// We are switching to a local charm, and this resource was not provided
		// by --resource, so no need to override existing resource.
		logger.Tracef("switching to a local charm, resource %q will not be upgraded because it was not provided by --resource", res.Name)
		return false, nil
	}
	// Never override existing resources a user has already uploaded.
	return cur.Origin != charmresource.OriginUpload, nil
}

const maxValueSize = 5242880 // Max size for a config file.

// ReadValue reads the value of an option out of the named file.
// An empty content is valid, like in parsing the options. The upper
// size is 5M.
func ReadValue(ctx *cmd.Context, filesystem modelcmd.Filesystem, filename string) (string, error) {
	absFilename := ctx.AbsPath(filename)
	fi, err := filesystem.Stat(absFilename)
	if err != nil {
		return "", errors.Errorf("cannot read option from file %q: %v", filename, err)
	}
	if fi.Size() > maxValueSize {
		return "", errors.Errorf("size of option file is larger than 5M")
	}
	content, err := os.ReadFile(ctx.AbsPath(filename))
	if err != nil {
		return "", errors.Errorf("cannot read option from file %q: %v", filename, err)
	}
	return string(content), nil
}

// isTerminal checks if the file descriptor is a terminal.
func IsTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}

	return isatty.IsTerminal(f.Fd())
}

type configFlag interface {
	// AbsoluteFileNames returns the absolute path of any file names specified.
	AbsoluteFileNames(ctx *cmd.Context) ([]string, error)

	// ReadConfigPairs returns just the k=v attributes
	ReadConfigPairs(ctx *cmd.Context) (map[string]interface{}, error)
}

// ProcessConfig processes the config defined by the config flag and returns
// the map of config values and any YAML file content.
// We may have a single file arg specified, in which case
// it points to a YAML file keyed on the charm name and
// containing values for any charm settings.
// We may also have key/value pairs representing
// charm settings which overrides anything in the YAML file.
// If more than one file is specified, that is an error.
func ProcessConfig(ctx *cmd.Context, filesystem modelcmd.Filesystem, configOptions configFlag, trust bool) (map[string]string, string, error) {
	var configYAML []byte
	files, err := configOptions.AbsoluteFileNames(ctx)
	if err != nil {
		return nil, "", errors.Trace(err)
	}
	if len(files) > 1 {
		return nil, "", errors.Errorf("only a single config YAML file can be specified, got %d", len(files))
	}
	if len(files) == 1 {
		configYAML, err = os.ReadFile(files[0])
		if err != nil {
			return nil, "", errors.Trace(err)
		}
	}
	attr, err := configOptions.ReadConfigPairs(ctx)
	if err != nil {
		return nil, "", errors.Trace(err)
	}
	appConfig := make(map[string]string)
	for k, v := range attr {
		appConfig[k] = v.(string)

		// Handle @ syntax for including file contents as values so we
		// are consistent to how 'juju config' works
		if len(appConfig[k]) < 1 || appConfig[k][0] != '@' {
			continue
		}

		if appConfig[k], err = ReadValue(ctx, filesystem, appConfig[k][1:]); err != nil {
			return nil, "", errors.Trace(err)
		}
	}

	// Expand the trust flag into the appConfig
	appConfig[app.TrustConfigOptionName] = strconv.FormatBool(trust)
	return appConfig, string(configYAML), nil
}
