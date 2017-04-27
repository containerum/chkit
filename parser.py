import argparse
import argcomplete

ONE_REQUIRED_ARGUMENT_ERROR = "you should pass at least one required argument: KIND or FILE"
KIND_OR_FILE_BOTH_ERROR = "you should pass either KIND, or FILE, not both"
NAME_OR_FILE_BOTH_ERROR = "you should pass either NAME, or FILE, not both"
NAME_WITH_KIND_ERROR = "NAME is required with KIND argument"


def create_parser(kinds, output_formats, version):
    parser = argparse.ArgumentParser(prog='client', description='This is help for %(prog)s')
    parser._positionals.title = 'client commands'

    parser.add_argument('--version', action='version', version='%(prog)s {}'.format(version))
    parser.add_argument("-d", '--debug', action='store_true', default=False, help='print debug messages to stderr')

    subparsers = parser.add_subparsers(help='use «[COMMAND] --help» to get detailed help for the command',  dest='command')

    parser_run = subparsers.add_parser('run', help='client run {deployment,deploy,deployments} NAME '
                                                   '—image=imagename [--replicas=1] '
                                                   '[--env="key1=value1"] [--env="key2=value2"]'
                                                   ' [--port=3000] [--port=3001] '
                                                   '[--command="/bin/bash"][--command="/bin/bash2"]'
                                                   ' [--volume="name:pathTo"] [--volume="name2:pathTo2"]')
    parser_run._optionals.title = 'run arguments'
    parser_run.add_argument('kind', help='object kind', choices=kinds)
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

    parser_create = subparsers.add_parser('create', help='create object')
    parser_create.add_argument('-f', '--file', help='input file', required=True)

    delete_usg = 'client delete (KIND -n NAME | -f FILE) [--namespace NAMESPACE]'
    parser_delete = subparsers.add_parser('delete', help='delete object', usage=delete_usg)
    parser_delete._optionals.title = 'delete arguments'
    parser_delete.add_argument('kind', help='object kind', choices=kinds)
    parser_delete.add_argument('name', help='object name to delete')
    parser_delete.add_argument('-f', '--file', help='input file')
    parser_delete.add_argument('--namespace', help='namespace, optional', required=False)

    parser_replace = subparsers.add_parser('replace', help='replace object')
    parser_replace.add_argument('-f', '--file', help='input file', required=True)
    parser_replace.add_argument('--namespace', help='namespace, optional', required=False)

    get_usg = 'client get (KIND [-n NAME] | -f FILE) [-o OUTPUT] [--namespace NAMESPACE] [--debug -D ]'
    parser_get = subparsers.add_parser('get', help='get object info', usage=get_usg)
    parser_get._optionals.title = 'get arguments'
    parser_get.add_argument('kind', help='object kind', choices=kinds)
    parser_get.add_argument('name', help='object name to get info, optional', nargs='*')
    parser_get.add_argument('-f', '--file', help='input file')

    parser_get.add_argument('-o', '--output', help='output format, default: json', choices=output_formats)
    parser_get.add_argument('--namespace', help='namespace, optional', required=False)

    config_usg = 'client config (--set-token -t TOKEN  --set_default_namespace -ns NAMESPACE)'
    parser_config = subparsers.add_parser('config', help='modify config', usage=config_usg)
    parser_config.add_argument('--set-token', '-t', help='token', required=False)
    parser_config.add_argument('--set_default_namespace', '-ns', help='default namespace', required=False)

    logout_usg = 'client logout'
    parser_logout = subparsers.add_parser('logout', help='logout user', usage=logout_usg)

    config_usg = 'client expose [-p --ports PORTNAME:TARGETPORT:PROTOCOL]'
    parser_expose = subparsers.add_parser('expose', help='expose service', usage=config_usg)
    parser_expose.add_argument('kind', help='object kind', choices=kinds)
    parser_expose.add_argument('name', help='object name to get info, optional', nargs='*')
    parser_expose.add_argument('--ports', '-p', help='target port', nargs='*', required=True)

    argcomplete.autocomplete(parser)

    return parser
