#coding=utf8
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

    parser_run = subparsers.add_parser('run', help='run deployment')
    parser_run._optionals.title = 'run arguments'
    parser_run.add_argument('-n', '--name', help='name, required', required=True)
    parser_run.add_argument('-i', '--image', help='image, required', required=True)
    parser_run.add_argument('-e', '--env', nargs='*', help='environment names, optional', required=False)
    parser_run.add_argument('-p', '--ports', type=int, nargs='*', help='ports, optional', required=False)
    parser_run.add_argument('-r', '--replicas', type=int, help='replicas, optional', required=False)
    parser_run.add_argument('--namespace', help='namespace, optional', required=False)

    parser_create = subparsers.add_parser('create', help='create object')
    parser_create.add_argument('-f', '--file', help='input file', required=True)

    delete_usg = 'client delete (--kind KIND -n NAME | -f FILE) [--namespace NAMESPACE]'
    parser_delete = subparsers.add_parser('delete', help='delete object', usage=delete_usg)
    parser_delete._optionals.title = 'delete arguments'
    parser_delete.add_argument('-k', '--kind', help='object kind', choices=kinds)
    parser_delete.add_argument('-n', '--name', help='object name to delete')
    parser_delete.add_argument('-f', '--file', help='input file')
    parser_delete.add_argument('--namespace', help='namespace, optional', required=False)

    parser_replace = subparsers.add_parser('replace', help='replace object')
    parser_replace.add_argument('-f', '--file', help='input file', required=True)
    parser_replace.add_argument('--namespace', help='namespace, optional', required=False)

    get_usg = 'client get (--kind KIND [-n NAME] | -f FILE) [-o OUTPUT] [--namespace NAMESPACE] [--debug -D {True|False}]'
    parser_get = subparsers.add_parser('get', help='get object info', usage=get_usg)
    parser_get._optionals.title = 'get arguments'
    parser_get.add_argument('-k', '--kind', help='object kind', choices=kinds)
    parser_get.add_argument('-f', '--file', help='input file')
    parser_get.add_argument('-n', '--name', help='object name to get info, optional', required=False)
    parser_get.add_argument('-o', '--output', help='output format, default: json', choices=output_formats)
    parser_get.add_argument('--namespace', help='namespace, optional', required=False)

    config_usg = 'client config (--set-token TOKEN)'
    parser_create = subparsers.add_parser('config', help='modify config', usage=config_usg)
    parser_create.add_argument('--set-token', '-t', help='token', required=True)

    logout_usg = 'client logout'
    parser_create = subparsers.add_parser('logout', help='logout user', usage=logout_usg)


    argcomplete.autocomplete(parser)

    return parser
