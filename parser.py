import argparse
import argcomplete
from data import kinds, output_formats, run_kinds, delete_kinds, expose_kinds, fields


ONE_REQUIRED_ARGUMENT_ERROR = "you should pass at least one required argument: KIND or FILE"
KIND_OR_FILE_BOTH_ERROR = "you should pass either KIND, or FILE, not both"
NAME_OR_FILE_BOTH_ERROR = "you should pass either NAME, or FILE, not both"
NAME_WITH_KIND_ERROR = "NAME is required with KIND argument"

formatter_class=lambda prog: MyFormatter(prog, max_help_position=80, width=140)


def create_parser(version):
    parser = argparse.ArgumentParser(prog='client', description='This is help for %(prog)s')
    parser._positionals.title = 'client commands'

    parser.add_argument('--version', action='version', version='%(prog)s {}'.format(version))
    parser.add_argument("-D", '--debug', action='store_true', default=False, help='print debug messages to stdout')

    subparsers = parser.add_subparsers(help='use «[COMMAND] --help» to get detailed help for the command',  dest='command')
    run_description = "Running deployement genereting json file"
    run_usg = 'client [--debug -d ] run  NAME --configure | —image -i IMAGE_NAME '\
                                                   '[--env -e "KEY=VALUE"]'\
                                                   '[--port -p PORT]'\
                                                   '[--replicas -r REPLICAS_COUNT]'\
                                                   '[--memory -m MEMORY]'\
                                                   '[--cpu -c CPU]'\
                                                   '[--command -cmd COMMAND] '\
                                                   '[--labels -ls "KEY=VALUE"]'\
                                                   '[--namespace -n NAMESPACE]'\
                                                   '[-h | --help]'
    parser_run = subparsers.add_parser('run', help=run_usg, usage=run_usg, description=run_description, formatter_class=formatter_class)
    parser_run._optionals.title = 'run arguments'
    parser_run.add_argument('name', help='name', metavar="NAME")
    parser_run.add_argument('--image', '-i', help='image', required=False)
    parser_run.add_argument('--env', '-e', nargs='*', help='environment names', required=False)
    parser_run.add_argument('--ports', '-p', type=int, nargs='*', help='ports', required=False)
    parser_run.add_argument('--commands', '-cmd', nargs='*', help='commands', required=False)
    parser_run.add_argument('--labels', '-ls', nargs='*', help='labels', required=False)
    parser_run.add_argument('--replicas', '-r', type=int, help='replicas, default: 1', default=1, required=False)
    parser_run.add_argument('--memory', '-m', help='memory, default: 128Mi', default="128Mi", required=False)
    parser_run.add_argument('--cpu', '-c', help='CPU share, default: 100m, ', default="100m", required=False)
    parser_run.add_argument('--namespace', '-n', help='namespace, default \"default\"', required=False)
    parser_run.add_argument('--configure', action='store_true', default=False, help='input params in console')

    create_usg = "client [--debug -d ] create (-f FILE | --file FILE)[-h | --help]"
    create_description = "Creating deployment from json file"
    parser_create = subparsers.add_parser('create', help=create_usg, usage=create_usg, description=create_description,
                                          formatter_class=formatter_class)
    parser_create.add_argument('-f', '--file', help='input file', required=True)


    delete_usg = 'client [--debug -d ] delete (KIND NAME | -f FILE) [--namespace NAMESPACE][-h | --help]'
    delete_description = "Deleting pods,service,deployments by name"
    parser_delete = subparsers.add_parser('delete', help=delete_usg, usage=delete_usg, description=delete_description,
                                          formatter_class=formatter_class)
    parser_delete._optionals.title = 'delete arguments'
    parser_delete.add_argument('kind', help='{deployment,service,pod} object kind', choices=delete_kinds, metavar="KIND")
    parser_delete.add_argument('name', help='object name to delete', metavar="NAME")
    parser_delete.add_argument('--file', '-f', help='input file')
    parser_delete.add_argument('--namespace', '-n', help='namespace, optional', required=False)

    # parser_replace = subparsers.add_parser('replace', help='replace object')
    # parser_replace.add_argument('--file', '-f', help='input file', required=True)
    # parser_replace.add_argument('--namespace', help='namespace, optional', required=False)

    get_usg = 'client [--debug -d ] get (KIND [NAME] | -f FILE) [-o OUTPUT] [--namespace NAMESPACE][-h | --help]'
    get_description = "Show info about pod(s), service(s), namespace(s), deployment(s)"
    parser_get = subparsers.add_parser('get', help=get_usg, usage=get_usg, description=get_description,
                                       formatter_class=formatter_class)
    parser_get._optionals.title = 'get arguments'
    parser_get.add_argument('kind', help='{namespace,deployment,service,pod} object kind', choices=kinds, metavar="KIND")
    parser_get.add_argument('name', help='object name to get info', metavar="NAME", nargs='*')
    parser_get.add_argument('--file', '-f', help='input file')
    parser_get.add_argument('--output', '-o', help='{yaml,json} output format, default: json', choices=output_formats, metavar="OUTPUT")
    parser_get.add_argument('--namespace', '-n', help='namespace, default: \"default\"', required=False)

    config_description = "Show and changing user's config settings"
    config_usg = 'client config (--set-token TOKEN  | --set-default-namespace NAMESPACE )[-h | --help]'
    parser_config = subparsers.add_parser('config', help=config_usg, usage=config_usg, description=config_description,
                                          formatter_class=formatter_class)
    parser_config.add_argument('--set-token', help='token', required=False)
    parser_config.add_argument('--set-default-namespace', help='default namespace', required=False)


    logout_usg = 'client logout'
    logout_description = "Clearing user's token from config"
    parser_logout = subparsers.add_parser('logout', help=logout_usg, usage=logout_usg, description=logout_description,
                                          formatter_class=formatter_class)

    expose_usg = 'client [--debug -d ] expose KIND NAME (-p --ports PORTS) [--name][-h | --help]'
    expose_description = "Exposing service genereting json file"
    parser_expose = subparsers.add_parser('expose', help=expose_usg, usage=expose_usg, description=expose_description,
                                          formatter_class=formatter_class)
    parser_expose.add_argument('kind', help='{deployment} object kind', choices=expose_kinds, metavar="KIND")
    parser_expose.add_argument('name', help='object name to get info', nargs='*', metavar="NAME")
    parser_expose.add_argument('--ports', '-p', help='target port, PORTS = PORTNAME:TARGETPORT:PROTOCOL, default: PROTOCOL = TCP', nargs='*', required=True)
    parser_expose.add_argument('--namespace', '-n', help='namespace, default: \"default\"', required=False)

    set_usg = 'client set FIELD (-f FILENAME | TYPE NAME) CONTAINER_NAME_1=CONTAINER_IMAGE_1 ... CONTAINER_NAME_N=CONTAINER_IMAGE_N'
    set_description = 'Change image in containers'
    parser_set = subparsers.add_parser('set', help=set_usg, usage=set_usg, description=set_description,
                                       formatter_class=formatter_class)
    parser_set._optionals.title = 'set arguments'
    parser_set.add_argument('field', help='{image} spec field', choices=fields, metavar="FIELD")
    parser_set.add_argument('kind', help='{deployment} object kind', choices=run_kinds, metavar="KIND")
    parser_set.add_argument('name', help='object name to get info', metavar="NAME", nargs='?')
    parser_set.add_argument('container', help='pair of container and image', metavar="CONTAINER", nargs='?')
    parser_set.add_argument('--file', '-f', help='input file')
    parser_set.add_argument('--namespace', '-n', help='namespace, default: \"default\"', required=False)
    argcomplete.autocomplete(parser)

    return parser


class MyFormatter(argparse.HelpFormatter):
    """
    Corrected _max_action_length for the indenting of subactions
    """
    def add_argument(self, action):
        if action.help is not argparse.SUPPRESS:

            # find all invocations
            get_invocation = self._format_action_invocation
            invocations = [get_invocation(action)]
            current_indent = self._current_indent
            for subaction in self._iter_indented_subactions(action):
                # compensate for the indent that will be added
                indent_chg = self._current_indent - current_indent
                added_indent = 'x'*indent_chg
                invocations.append(added_indent+get_invocation(subaction))
            # print('inv', invocations)

            # update the maximum item length
            invocation_length = max([len(s) for s in invocations])
            action_length = invocation_length + self._current_indent
            self._action_max_length = max(self._action_max_length,
                                          action_length)

            # add the item to the list
            self._add_item(self._format_action, [action])