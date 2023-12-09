// See LICENSE file for copyright and license details

package help

import "djmo.ch/go-shrt/cmd/shrt/internal/base"

var EnvCmd = &base.Command{
	Name:      "environment",
	ShortHelp: "environment variables",
	LongHelp: `
The shrt command consults environment variables for configuration. If
an environment variable is unset, the shrt command uses a sensible
default setting. To see the effective setting of the variable <NAME>,
run 'shrt env <NAME>'. To change the default setting, run 'shrt env -w
<NAME>=<VALUE>'. Defaults changed using 'shrt env -w' are recorded in
a Shrt environment configuration file stored in /etc/shrt/config on
Unix systems and C:\ProgramData\shrt\config on Windows. The location
of the configuration file can be changed by setting the environment
variable DGITENV, and 'shrt env DGITENV' prints the effective
location, but 'shrt env -w' cannot change the default location. See
'shrt help env' for details.

Environment variables:

	SHRTENV
		The location of the Shrt environment configuration
		file. Cannot be set using 'shrt env -w'.
	SHRT_SRVNAME
		The server name of the Shrt host.
	SHRT_SCMTYPE
		The SCM (or VCS) type.
	SHRT_SUFFIX
		The SCM repository suffix, if required by repository
		host.
	SHRT_RDRNAME
		The server name of the repository host.
	SHRT_BARERDR
		Where requests with an empty path should redirect.
	SHRT_DBPATH
		The absolute path to the database file.
	SHRT_GOSOURCEDIR
		The string to append to the URL for go-get redirects
		to form the directory entry in the go-source meta tag.
		This key is experimental and may be removed in a
		future release.
	SHRT_GOSOURCEFILE
		The string to append to the URL for go-get redirects
		to form the file entry in the go-source meta tag. This
		key is experimental and may be removed in a future
		release.
`,
}
