import argparse
import argcomplete


def create_parser(kinds, output_formats, version):
    parser = argparse.ArgumentParser(prog='client', description='This is help for %(prog)s')

    parser._positionals.title = 'client commands'

    parser.add_argument('--version', action='version', version='%(prog)s {}'.format(version))

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
    parser_delete.add_argument('-k', '--kind', help='oqbject kind', choices=kinds)
    parser_delete.add_argument('-n', '--name', help='object name to delete')
    parser_delete.add_argument('-f', '--file', help='input file')
    parser_delete.add_argument('--namespace', help='namespace, optional', required=False)

    parser_replace = subparsers.add_parser('replace', help='replace object')
    parser_replace.add_argument('-f', '--file', help='input file', required=True)
    parser_replace.add_argument('--namespace', help='namespace, optional', required=False)

    get_usg = 'client get (--kind KIND [-n NAME] | -f FILE) [-o OUTPUT] [--namespace NAMESPACE]'
    parser_get = subparsers.add_parser('get', help='get object info', usage=get_usg)
    parser_get._optionals.title = 'get arguments'
    parser_get.add_argument('-k', '--kind', help='object kind', choices=kinds)
    parser_get.add_argument('-f', '--file', help='input file')
    parser_get.add_argument('-n', '--name', help='object name to get info, optional', required=False)
    parser_get.add_argument('-o', '--output', help='output format, default: json', choices=output_formats)
    parser_get.add_argument('--namespace', help='namespace, optional', required=False)

    argcomplete.autocomplete(parser)

    return parser
