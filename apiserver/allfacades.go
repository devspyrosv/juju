// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver

// This file imports all of the facades so they get registered at runtime.
// When adding a new facade implementation, import it here so that its init()
// function will get called to register it.
//
// TODO(fwereade): this is silly. We should be declaring our full API in *one*
// place, not scattering it across packages and depending on magic import lists.
import (
	_ "github.com/juju/juju/apiserver/action"
	_ "github.com/juju/juju/apiserver/agent"
	_ "github.com/juju/juju/apiserver/agenttools"
	_ "github.com/juju/juju/apiserver/annotations"
	_ "github.com/juju/juju/apiserver/application"
	_ "github.com/juju/juju/apiserver/applicationscaler"
	_ "github.com/juju/juju/apiserver/backups"
	_ "github.com/juju/juju/apiserver/block"
	_ "github.com/juju/juju/apiserver/charmrevisionupdater"
	_ "github.com/juju/juju/apiserver/charms"
	_ "github.com/juju/juju/apiserver/cleaner"
	_ "github.com/juju/juju/apiserver/client"
	_ "github.com/juju/juju/apiserver/cloud"
	_ "github.com/juju/juju/apiserver/controller"
	_ "github.com/juju/juju/apiserver/deployer"
	_ "github.com/juju/juju/apiserver/discoverspaces"
	_ "github.com/juju/juju/apiserver/diskmanager"
	_ "github.com/juju/juju/apiserver/firewaller"
	_ "github.com/juju/juju/apiserver/highavailability"
	_ "github.com/juju/juju/apiserver/hostkeyreporter"
	_ "github.com/juju/juju/apiserver/imagemanager"
	_ "github.com/juju/juju/apiserver/imagemetadata"
	_ "github.com/juju/juju/apiserver/instancepoller"
	_ "github.com/juju/juju/apiserver/keymanager"
	_ "github.com/juju/juju/apiserver/keyupdater"
	_ "github.com/juju/juju/apiserver/leadership"
	_ "github.com/juju/juju/apiserver/lifeflag"
	_ "github.com/juju/juju/apiserver/logfwd"
	_ "github.com/juju/juju/apiserver/logger"
	_ "github.com/juju/juju/apiserver/machine"
	_ "github.com/juju/juju/apiserver/machineactions"
	_ "github.com/juju/juju/apiserver/machinemanager"
	_ "github.com/juju/juju/apiserver/meterstatus"
	_ "github.com/juju/juju/apiserver/metricsadder"
	_ "github.com/juju/juju/apiserver/metricsdebug"
	_ "github.com/juju/juju/apiserver/metricsmanager"
	_ "github.com/juju/juju/apiserver/migrationflag"
	_ "github.com/juju/juju/apiserver/migrationmaster"
	_ "github.com/juju/juju/apiserver/migrationminion"
	_ "github.com/juju/juju/apiserver/migrationtarget"
	_ "github.com/juju/juju/apiserver/modelmanager"
	_ "github.com/juju/juju/apiserver/provisioner"
	_ "github.com/juju/juju/apiserver/proxyupdater"
	_ "github.com/juju/juju/apiserver/reboot"
	_ "github.com/juju/juju/apiserver/resumer"
	_ "github.com/juju/juju/apiserver/retrystrategy"
	_ "github.com/juju/juju/apiserver/singular"
	_ "github.com/juju/juju/apiserver/spaces"
	_ "github.com/juju/juju/apiserver/sshclient"
	_ "github.com/juju/juju/apiserver/statushistory"
	_ "github.com/juju/juju/apiserver/storage"
	_ "github.com/juju/juju/apiserver/storageprovisioner"
	_ "github.com/juju/juju/apiserver/subnets"
	_ "github.com/juju/juju/apiserver/undertaker"
	_ "github.com/juju/juju/apiserver/unitassigner"
	_ "github.com/juju/juju/apiserver/uniter"
	_ "github.com/juju/juju/apiserver/upgrader"
	_ "github.com/juju/juju/apiserver/usermanager"
)
