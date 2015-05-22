package main

import (
	"github.com/buildkite/agent/command"
	"github.com/codegangsta/cli"
)

var Commands []cli.Command

var AgentDescription = `Usage:

   buildkite-agent start [arguments...]

Description:

   When a job is ready to run it will call the "bootstrap-script"
   and pass it all the environment variables required for the job to run.
   This script is responsible for checking out the code, and running the
   actual build script defined in the project.

   The agent will run any jobs within a PTY (pseudo terminal) if available.

Example:

   $ buildkite-agent start --token xxx`

var ShasumHelpDescription = `Usage:

   buildkite-agent artifact shasum [arguments...]

Description:

   Prints to STDOUT the SHA-1 for the artifact provided. If your search query
   for artifacts matches multiple agents, and error will be raised.

   Note: You need to ensure that your search query is surrounded by quotes if
   using a wild card as the built-in shell path globbing will provide files,
   which will break the download.

Example:

   $ buildkite-agent artifact shasum "pkg/release.tar.gz" --build xxx

   This will search for all the files in the build with the path "pkg/release.tar.gz" and will
   print to STDOUT it's SHA-1 checksum.

   If you would like to target artifacts from a specific build step, you can do
   so by using the --job argument.

   $ buildkite-agent artifact shasum "pkg/release.tar.gz" --job "release" --build xxx

   You can also use the job's id (provided by the environment variable $BUILDKITE_JOB_ID)`

var DownloadHelpDescription = `Usage:

   buildkite-agent artifact download [arguments...]

Description:

   Downloads artifacts from Buildkite to the local machine.

   Note: You need to ensure that your search query is surrounded by quotes if
   using a wild card as the built-in shell path globbing will provide files,
   which will break the download.

Example:

   $ buildkite-agent artifact download "pkg/*.tar.gz" . --build xxx

   This will search across all the artifacts for the build with files that match that part.
   The first argument is the search query, and the second argument is the download destination.

   If you're trying to download a specific file, and there are multiple artifacts from different
   jobs, you can target the particular job you want to download the artifact from:

   $ buildkite-agent artifact download "pkg/*.tar.gz" . --job "tests" --build xxx

   You can also use the job's id (provided by the environment variable $BUILDKITE_JOB_ID)`

var UploadHelpDescription = `Usage:

   buildkite-agent artifact upload <pattern> <destination> [arguments...]

Description:

   Uploads files to a job as artifacts.

   You need to ensure that the paths are surrounded by quotes otherwise the
   built-in shell path globbing will provide the files, which is currently not
   supported.

Example:

   $ buildkite-agent artifact upload "log/**/*.log"

   You can also upload directly to Amazon S3 if you'd like to host your own artifacts:

   $ export BUILDKITE_S3_ACCESS_KEY_ID=xxx
   $ export BUILDKITE_S3_SECRET_ACCESS_KEY=yyy
   $ export BUILDKITE_S3_DEFAULT_REGION=eu-central-1 # default is us-east-1
   $ export BUILDKITE_S3_ACL=private # default is public-read
   $ buildkite-agent artifact upload "log/**/*.log" s3://name-of-your-s3-bucket/$BUILDKITE_JOB_ID`

var SetHelpDescription = `Usage:

   buildkite-agent meta-data set <key> <value> [arguments...]

Description:

   Set arbitrary data on a build using a basic key/value store.

Example:

   $ buildkite-agent meta-data set "foo" "bar"`

var GetHelpDescription = `Usage:

   buildkite-agent meta-data get <key> [arguments...]

Description:

   Get data from a builds key/value store.

Example:

   $ buildkite-agent meta-data get "foo"`

var DefaultEndpoint = "https://agent.buildkite.com/v2"

func init() {
	Commands = []cli.Command{
		{
			Name:        "start",
			Usage:       "Starts a Buildkite agent",
			Description: AgentDescription,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "config",
					Value:  "",
					Usage:  "Path to a configration file",
					EnvVar: "BUILDKITE_AGENT_CONFIG",
				},
				cli.StringFlag{
					Name:   "token",
					Value:  "",
					Usage:  "Your account agent token",
					EnvVar: "BUILDKITE_AGENT_TOKEN",
				},
				cli.StringFlag{
					Name:   "name",
					Value:  "",
					Usage:  "The name of the agent",
					EnvVar: "BUILDKITE_AGENT_NAME",
				},
				cli.StringFlag{
					Name:   "priority",
					Value:  "",
					Usage:  "The priority of the agent (higher priorities are assigned work first)",
					EnvVar: "BUILDKITE_AGENT_PRIORITY",
				},
				cli.StringSliceFlag{
					Name:   "meta-data",
					Value:  &cli.StringSlice{},
					Usage:  "Meta data for the agent (default is \"queue=default\")",
					EnvVar: "BUILDKITE_AGENT_META_DATA",
				},
				cli.BoolFlag{
					Name:  "meta-data-ec2-tags",
					Usage: "Populate the meta data from the current instances EC2 Tags",
				},
				cli.StringFlag{
					Name:   "bootstrap-script",
					Value:  "",
					Usage:  "Path to the bootstrap script",
					EnvVar: "BUILDKITE_BOOTSTRAP_SCRIPT_PATH",
				},
				cli.StringFlag{
					Name:   "build-path",
					Value:  "",
					Usage:  "Path to where the builds will run from",
					EnvVar: "BUILDKITE_BUILD_PATH",
				},
				cli.StringFlag{
					Name:   "hooks-path",
					Value:  "",
					Usage:  "Directory where the hook scripts are found",
					EnvVar: "BUILDKITE_HOOKS_PATH",
				},
				cli.BoolFlag{
					Name:   "no-pty",
					Usage:  "Do not run jobs within a pseudo terminal",
					EnvVar: "BUILDKITE_NO_PTY",
				},
				cli.BoolFlag{
					Name:   "no-automatic-ssh-fingerprint-verification",
					Usage:  "Don't automatically verify SSH fingerprints",
					EnvVar: "BUILDKITE_NO_AUTOMATIC_SSH_FINGERPRINT_VERIFICATION",
				},
				cli.BoolFlag{
					Name:   "no-command-eval",
					Usage:  "Don't allow this agent to run arbitrary console commands",
					EnvVar: "BUILDKITE_NO_COMMAND_EVAL",
				},
				cli.StringFlag{
					Name:   "endpoint",
					Value:  DefaultEndpoint,
					Usage:  "The Agent API endpoint",
					EnvVar: "BUILDKITE_AGENT_ENDPOINT",
				},
				cli.BoolFlag{
					Name:   "debug",
					Usage:  "Enable debug mode",
					EnvVar: "BUILDKITE_AGENT_DEBUG",
				},
				cli.BoolFlag{
					Name:   "no-color",
					Usage:  "Don't show colors in logging",
					EnvVar: "BUILDKITE_AGENT_NO_COLOR",
				},
			},
			Action: command.AgentStartCommandAction,
		},
		{
			Name:  "artifact",
			Usage: "Upload/download artifacts from Buildkite jobs",
			Subcommands: []cli.Command{
				{
					Name:        "download",
					Usage:       "Downloads artifacts from Buildkite to the local machine",
					Description: DownloadHelpDescription,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "step",
							Value: "",
							Usage: "Scope the search to a paticular step by using either it's name of job ID",
						},
						cli.StringFlag{
							Name:  "build",
							Value: "BUILDKITE_JOB_ID",
							Usage: "The build that the artifacts were uploaded to",
						},
						cli.StringFlag{
							Name:   "agent-access-token",
							Value:  "",
							Usage:  "The access token used to identify the agent",
							EnvVar: "BUILDKITE_AGENT_ACCESS_TOKEN",
						},
						cli.StringFlag{
							Name:   "endpoint",
							Value:  DefaultEndpoint,
							Usage:  "The Agent API endpoint",
							EnvVar: "BUILDKITE_AGENT_ENDPOINT",
						},
						cli.BoolFlag{
							Name:   "debug",
							Usage:  "Enable debug mode",
							EnvVar: "BUILDKITE_AGENT_DEBUG",
						},
						cli.BoolFlag{
							Name:   "no-color",
							Usage:  "Don't show colors in logging",
							EnvVar: "BUILDKITE_AGENT_NO_COLOR",
						},
					},
					Action: command.ArtifactDownloadCommandAction,
				},
				{
					Name:        "upload",
					Usage:       "Uploads files to a job as artifacts",
					Description: UploadHelpDescription,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "job",
							Value:  "",
							Usage:  "Which job should the artifacts be uploaded to",
							EnvVar: "BUILDKITE_JOB_ID",
						},
						cli.StringFlag{
							Name:   "agent-access-token",
							Value:  "",
							Usage:  "The access token used to identify the agent",
							EnvVar: "BUILDKITE_AGENT_ACCESS_TOKEN",
						},
						cli.StringFlag{
							Name:   "endpoint",
							Value:  DefaultEndpoint,
							Usage:  "The Agent API endpoint",
							EnvVar: "BUILDKITE_AGENT_ENDPOINT",
						},
						cli.BoolFlag{
							Name:   "debug",
							Usage:  "Enable debug mode",
							EnvVar: "BUILDKITE_AGENT_DEBUG",
						},
						cli.BoolFlag{
							Name:   "no-color",
							Usage:  "Don't show colors in logging",
							EnvVar: "BUILDKITE_AGENT_NO_COLOR",
						},
					},
					Action: command.ArtifactUploadCommandAction,
				},
				{
					Name:        "shasum",
					Usage:       "Prints the SHA-1 checksum for the artifact provided to STDOUT",
					Description: ShasumHelpDescription,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "step",
							Value: "",
							Usage: "Scope the search to a paticular step by using either it's name of job ID",
						},
						cli.StringFlag{
							Name:  "build",
							Value: "BUILDKITE_JOB_ID",
							Usage: "The build that the artifacts were uploaded to",
						},
						cli.StringFlag{
							Name:   "agent-access-token",
							Value:  "",
							Usage:  "The access token used to identify the agent",
							EnvVar: "BUILDKITE_AGENT_ACCESS_TOKEN",
						},
						cli.StringFlag{
							Name:   "endpoint",
							Value:  DefaultEndpoint,
							Usage:  "The Agent API endpoint",
							EnvVar: "BUILDKITE_AGENT_ENDPOINT",
						},
						cli.BoolFlag{
							Name:   "debug",
							Usage:  "Enable debug mode",
							EnvVar: "BUILDKITE_AGENT_DEBUG",
						},
						cli.BoolFlag{
							Name:   "no-color",
							Usage:  "Don't show colors in logging",
							EnvVar: "BUILDKITE_AGENT_NO_COLOR",
						},
					},
					Action: command.ArtifactShasumCommandAction,
				},
			},
		},
		{
			Name:  "meta-data",
			Usage: "Get/set data from Buildkite jobs",
			Subcommands: []cli.Command{
				{
					Name:        "set",
					Usage:       "Set data on a build",
					Description: SetHelpDescription,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "job",
							Value:  "",
							Usage:  "Which job should the meta-data be set on",
							EnvVar: "BUILDKITE_JOB_ID",
						},
						cli.StringFlag{
							Name:   "agent-access-token",
							Value:  "",
							Usage:  "The access token used to identify the agent",
							EnvVar: "BUILDKITE_AGENT_ACCESS_TOKEN",
						},
						cli.StringFlag{
							Name:   "endpoint",
							Value:  DefaultEndpoint,
							Usage:  "The Agent API endpoint",
							EnvVar: "BUILDKITE_AGENT_ENDPOINT",
						},
						cli.BoolFlag{
							Name:   "debug",
							Usage:  "Enable debug mode",
							EnvVar: "BUILDKITE_AGENT_DEBUG",
						},
						cli.BoolFlag{
							Name:   "no-color",
							Usage:  "Don't show colors in logging",
							EnvVar: "BUILDKITE_AGENT_NO_COLOR",
						},
					},
					Action: command.DataSetCommandAction,
				},
				{
					Name:        "get",
					Usage:       "Get data from a build",
					Description: GetHelpDescription,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "job",
							Value:  "",
							Usage:  "Which job should the meta-data be retrieved from",
							EnvVar: "BUILDKITE_JOB_ID",
						},
						cli.StringFlag{
							Name:   "agent-access-token",
							Value:  "",
							Usage:  "The access token used to identify the agent",
							EnvVar: "BUILDKITE_AGENT_ACCESS_TOKEN",
						},
						cli.StringFlag{
							Name:   "endpoint",
							Value:  DefaultEndpoint,
							Usage:  "The Agent API endpoint",
							EnvVar: "BUILDKITE_AGENT_ENDPOINT",
						},
						cli.BoolFlag{
							Name:   "debug",
							Usage:  "Enable debug mode",
							EnvVar: "BUILDKITE_AGENT_DEBUG",
						},
						cli.BoolFlag{
							Name:   "no-color",
							Usage:  "Don't show colors in logging",
							EnvVar: "BUILDKITE_AGENT_NO_COLOR",
						},
					},
					Action: command.DataGetCommandAction,
				},
			},
		},
	}
}
