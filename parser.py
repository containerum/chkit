import argparse
import argcomplete
from data import kinds, output_formats, run_kinds, delete_kinds, expose_kinds

ONE_REQUIRED_ARGUMENT_ERROR = "you should pass at least one required argument: KIND or FILE"
KIND_OR_FILE_BOTH_ERROR = "you should pass either KIND, or FILE, not both"
NAME_OR_FILE_BOTH_ERROR = "you should pass either NAME, or FILE, not both"
NAME_WITH_KIND_ERROR = "NAME is required with KIND argument"


def create_parser(version):
    parser = argparse.ArgumentParser(prog='client', description='This is help for %(prog)s')
    parser._positionals.title = 'client commands'

    parser.add_argument('--version', action='version', version='%(prog)s {}'.format(version))
    parser.add_argument("-D", '--debug', action='store_true', default=False, help='print debug messages to stdout')

    subparsers = parser.add_subparsers(help='use «[COMMAND] --help» to get detailed help for the command',  dest='command')
    run_description = "Running deployement genereting json file"
    run_usg = 'client run {deployment,deploy,deployments} NAME\\ '\
                                                   '—image IMAGE_NAME [--replicas=1]\\ '\
                                                   '[--env="key1=value1"][--env="key2=value2"]\\'\
                                                   '[--port=3000] [--port=3001]\\ '\
                                                   '[--command="/bin/bash"][--command="/bin/bash2"]\\ '\
                                                   '[--volume="name:pathTo"][--volume="name2:pathTo2"]\\' \
                                                   '[-h | --help]'
    parser_run = subparsers.add_parser('run', help=run_usg, usage=run_usg, description=run_description)
    parser_run._optionals.title = 'run arguments'
    parser_run.add_argument('kind', help='object kind', choices=run_kinds)
    parser_run.add_argument('name', help='name, required')
    parser_run.add_argument('-i', '--image', help='image, required', required=False)
    parser_run.add_argument('-e', '--env', nargs='*', help='environment names, optional', required=False)
    parser_run.add_argument('-p', '--ports', type=int, nargs='*', help='ports , optional', required=False)
    parser_run.add_argument('--commands', nargs='*', help='commands , optional', required=False)
    parser_run.add_argument('--labels', nargs='*', help='labels , optional', required=False)
    parser_run.add_argument('-r', '--replicas', type=int, help='replicas, optional', default=1, required=False)
    parser_run.add_argument('-m', '--memory', help='memory, optional', default="128Mi", required=False)
    parser_run.add_argument('-c', '--cpu', help='cpu, optional', default="100m", required=False)
    parser_run.add_argument('--namespace', help='namespace, optional', required=False)
    parser_run.add_argument('--configure', action='store_true', default=False, help='input params in console')

    create_usg = "client create (-f FILE | --file FILE)[-h | --help]"
    create_description = "Creating deployment from json file"
    parser_create = subparsers.add_parser('create', help=create_usg, usage=create_usg, description=create_description)
    parser_create.add_argument('-f', '--file', help='input file', required=True)

    delete_usg = 'client delete (KIND NAME | -f FILE) [--namespace NAMESPACE][-h | --help]'
    delete_description = "Deleting pods,service,deployments by name"
    parser_delete = subparsers.add_parser('delete', help=delete_usg, usage=delete_usg, description=delete_description)
    parser_delete._optionals.title = 'delete arguments'
    parser_delete.add_argument('kind', help='object kind', choices=delete_kinds)
    parser_delete.add_argument('name', help='object name to delete')
    parser_delete.add_argument('-f', '--file', help='input file')
    parser_delete.add_argument('--namespace', help='namespace, optional', required=False)

    parser_replace = subparsers.add_parser('replace', help='replace object')
    parser_replace.add_argument('-f', '--file', help='input file', required=True)
    parser_replace.add_argument('--namespace', help='namespace, optional', required=False)

    get_usg = 'client get (KIND [NAME] | -f FILE) [-o OUTPUT] [--namespace NAMESPACE] [--debug -D ][-h | --help]'
    get_description = "Show info about pod(s), service(s), namespace(s), deployment(s)"
    parser_get = subparsers.add_parser('get', help=get_usg, usage=get_usg, description=get_description)
    parser_get._optionals.title = 'get arguments'
    parser_get.add_argument('kind', help='object kind', choices=kinds)
    parser_get.add_argument('name', help='object name to get info, optional', nargs='*')
    parser_get.add_argument('-f', '--file', help='input file')

    parser_get.add_argument('-o', '--output', help='output format, default: json', choices=output_formats)
    parser_get.add_argument('--namespace', help='namespace, optional', required=False)

    config_description = "Show and changing user's config settings"
    config_usg = 'client config (--set-token -t TOKEN  | --set-default-namespace -ns NAMESPACE | -v)[-h | --help]'
    parser_config = subparsers.add_parser('config', help=config_usg, usage=config_usg, description=config_description)
    parser_config.add_argument('--set-token', '-t', help='token', required=False)
    parser_config.add_argument('--set-default-namespace', '-ns', help='default namespace', required=False)
    parser_config.add_argument("-v", action='store_true', default=False, help='print current config settings')

    logout_usg = 'client logout'
    logout_description = "Clearing user's token from config"
    parser_logout = subparsers.add_parser('logout', help=logout_usg, usage=logout_usg, description=logout_description)

    expose_usg = 'client expose KIND NAME [-p --ports PORTNAME:TARGETPORT:PROTOCOL][-h | --help]'
    expose_description = "Exposing service genereting json file"
    parser_expose = subparsers.add_parser('expose', help=expose_usg, usage=expose_usg, description=expose_description)
    parser_expose.add_argument('kind', help='object kind', nargs='*', default="deploy")
    parser_expose.add_argument('name', help='object name to get info, optional', nargs='*')
    parser_expose.add_argument('--ports', '-p', help='target port', nargs='*', required=True)

    argcomplete.autocomplete(parser)

    return parser
