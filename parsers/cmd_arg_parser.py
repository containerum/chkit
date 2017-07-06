import argparse

import argcomplete

from variables.data import kinds, output_formats, run_kinds, delete_kinds, expose_kinds, fields

ONE_REQUIRED_ARGUMENT_ERROR = "you should pass at least one required argument: KIND or FILE"
KIND_OR_FILE_BOTH_ERROR = "you should pass either KIND, or FILE, not both"
NAME_OR_FILE_BOTH_ERROR = "you should pass either NAME, or FILE, not both"
NAME_WITH_KIND_ERROR = "NAME is required with KIND argument"

formatter_class=lambda prog: MyFormatter(prog, max_help_position=80, width=140)


def create_parser(version):
    parser = argparse.ArgumentParser(prog='chkit', description='This is help for %(prog)s')
    parser._positionals.title = 'chkit commands'

    parser.add_argument('--version', action='version', version='%(prog)s {}'.format(version))
    parser.add_argument("-d", '--debug', action='store_true', default=False, help='print debug messages to stdout')
    subparsers = parser.add_subparsers(help='use «[COMMAND] --help» to get detailed help for the command',  dest='command')

    config_description = "Show and changing user's config settings"
    config_usg = 'chkit [--debug -d ] config (--set-token -t TOKEN  | --set-default-namespace -n NAMESPACE )[-h | --help]'
    parser_config = subparsers.add_parser('config', help=config_usg, usage=config_usg, description=config_description, formatter_class=formatter_class)
    parser_config.add_argument('--set-token', '-t', help='token', required=False)
    parser_config.add_argument('--set-default-namespace', '-n', help='default namespace', required=False)

    run_description = "Running deployement genereting json file"
    run_usg = 'chkit [--debug -d ] run  NAME --configure | --image -i IMAGE '\
                                                   '[--env -e "KEY=VALUE"]'\
                                                   '[--port -p PORT]'\
                                                   '[--replicas -r REPLICAS_COUNT]'\
                                                   '[--memory -m MEMORY]'\
                                                   '[--cpu -c CPU]'\
                                                   '[--command -cmd COMMAND] '\
                                                   '[--labels -ls "KEY=VALUE"]'\
                                                   '[--namespace -n NAMESPACE]'\
                                                   '[-h  --help]'
    parser_run = subparsers.add_parser('run', help=run_usg, usage=run_usg, description=run_description, formatter_class=formatter_class)
    parser_run._optionals.title = 'run arguments'
    parser_run.add_argument('name', help='name', metavar="NAME")
    parser_run.add_argument('--image', '-i', help='image', required=False)
    parser_run.add_argument('--env', '-e', nargs='*', help='environment names', required=False)
    parser_run.add_argument('--ports', '-p', type=int, nargs='*', help='ports', required=False)
    parser_run.add_argument('--commands', '-cmd', nargs='*', help='commands', required=False)
    parser_run.add_argument('--labels', '-ls', nargs='*', help='labels', required=False)
    parser_run.add_argument('--replicas', '-r', type=int, help='replicas, default: 1', default=1,
                            required=False)
    parser_run.add_argument('--memory', '-m', help='memory, default: 128Mi', default="128Mi", required=False)
    parser_run.add_argument('--cpu', '-c', help='CPU share, default: 100m, ', default="100m", required=False)
    parser_run.add_argument('--namespace', '-n', help='namespace, default \"default\"', required=False)
    parser_run.add_argument('--configure', action='store_true', default=False, help='input params in console')

    create_usg = "chkit [--debug -d ] create (--file -f FILE)[-h --help]"
    create_description = "Creating deployment from json file"
    parser_create = subparsers.add_parser('create', help=create_usg, usage=create_usg, description=create_description,
                                          formatter_class=formatter_class)
    parser_create.add_argument('--file', '-f', help='input file', required=True)

    expose_usg = 'chkit [--debug -d ] expose KIND NAME (-p --ports PORTS) [--name][-h | --help]'
    expose_description = "Exposing service genereting json file"
    parser_expose = subparsers.add_parser('expose', help=expose_usg, usage=expose_usg, description=expose_description,
                                          formatter_class=formatter_class)
    parser_expose.add_argument('kind', help='{deployment} object kind', choices=expose_kinds, metavar="KIND")
    parser_expose.add_argument('name', help='object name to get info', nargs='?', metavar="NAME")
    parser_expose.add_argument('--ports', '-p', help='target port, for external services PORTS = PORTNAME:TARGETPORT[:PROTOCOL],'
                                                     ' for internal services '
                                                     'PORTS = PORTNAME:TARGETPORT:PORT[:PROTOCOL]'
                                                     ' default: PROTOCOL = TCP', nargs='*', required=True)
    parser_expose.add_argument('--namespace', '-n', help='namespace, default: \"default\"', required=False)

    set_usg = 'chkit [--debug -d] set FIELD TYPE NAME CONTAINER_NAME=CONTAINER_IMAGE|COUNT [-n --namespace NAMESPACE][--help | -h]'
    set_description = 'Change image in containers | set replicas count'
    parser_set = subparsers.add_parser('set', help=set_usg, usage=set_usg, description=set_description,
                                       formatter_class=formatter_class)
    parser_set._optionals.title = 'set arguments'
    parser_set.add_argument('field', help='{image} spec field', choices=fields, metavar="FIELD")
    parser_set.add_argument('kind', help='{deployment} object kind', choices=run_kinds, metavar="KIND")
    parser_set.add_argument('name', help='object name to get info', metavar="NAME", nargs='+')
    parser_set.add_argument('args', help='pair of container and image|count of replicas', metavar="ARGS", nargs='+')
    #parser_set.add_argument('--file', '-f', help='input file')
    parser_set.add_argument('--namespace', '-n', help='namespace, default: \"default\"', required=False)

    get_usg = 'chkit [--debug -d ] get (KIND [NAME] | --file -f FILE) ' \
              '[--output -o OUTPUT] [--namespace -n NAMESPACE][--deploy -d DEPLOY][-h | --help]'
    get_description = "Show info about pod(s), service(s), namespace(s), deployment(s)"
    parser_get = subparsers.add_parser('get', help=get_usg, usage=get_usg, description=get_description,
                                       formatter_class=formatter_class)
    parser_get._optionals.title = 'get arguments'
    parser_get.add_argument('kind', help='{namespace,deployment,service,pod} object kind', choices=kinds, default="", metavar="KIND", nargs='?')
    parser_get.add_argument('name', help='object name to get info', metavar="NAME", nargs='?')
    parser_get.add_argument('--file', '-f', help='input file')
    parser_get.add_argument('--output', '-o', help='{yaml,json} output format, default: json', choices=output_formats, metavar="OUTPUT")
    parser_get.add_argument('--namespace', '-n', help='namespace, default: \"default\"', required=False)
    parser_get.add_argument('--deploy', '-d', help='filtering by deploy(only for pods ans services!)', required=False)

    restart_usg = 'chkit [--debug -d ] restart NAME [--namespace NAMESPACE][-h | --help]'
    restart_description = "Restarting pods by deploy name"
    parser_restart = subparsers.add_parser('restart', help=restart_usg, usage=restart_usg,
                                           description=restart_description,
                                           formatter_class=formatter_class)
    parser_restart._optionals.title = 'restart arguments'
    parser_restart.add_argument('name', help='deploy name to restart', metavar="NAME")
    parser_restart.add_argument('--namespace', '-n', help='namespace, optional', required=False)

    delete_usg = 'chkit [--debug -d ] delete (KIND NAME | --file -f FILE) [--pods][--namespace NAMESPACE][-h | --help]'
    delete_description = "Deleting pods,service,deployments by name"
    parser_delete = subparsers.add_parser('delete', help=delete_usg, usage=delete_usg, description=delete_description,
                                          formatter_class=formatter_class)
    parser_delete._optionals.title = 'delete arguments'
    parser_delete.add_argument('kind', help='{deployment,service,pod} object kind', nargs="?", choices=delete_kinds, metavar="KIND")
    parser_delete.add_argument('name', help='object name to delete', metavar="NAME", nargs="?")
    parser_delete.add_argument('--file', '-f', help='input file')
    parser_delete.add_argument('--pods', action='store_true', default=False, help='delete all pods in deploy')
    parser_delete.add_argument('--namespace', '-n', help='namespace, optional', required=False)

    # parser_replace = subparsers.add_parser('replace', help='replace object')
    # parser_replace.add_argument('--file', '-f', help='input file', required=True)
    # parser_replace.add_argument('--namespace', help='namespace, optional', required=False)

    update_usg = 'chkit update'
    update_description = "update app through GitHub Releases"
    parser_update = subparsers.add_parser('update', help=update_usg, usage=update_usg, description=update_description,
                                         formatter_class=formatter_class)

    login_usg = 'chkit login'
    login_description = "Sign in. Sets user's token to config"
    parser_login = subparsers.add_parser('login', help=login_usg, usage=login_usg, description=login_description,
                                         formatter_class=formatter_class)

    logout_usg = 'chkit logout'
    logout_description = "Clearing user's token from config"
    parser_logout = subparsers.add_parser('logout', help=logout_usg, usage=logout_usg, description=logout_description,
                                          formatter_class=formatter_class)

    scale_usg = 'chkit [--debug -d] scale KIND NAME COUNT [-n --namespace NAMESPACE][--help | -h]'
    scale_description = "Change replicas count"
    parser_scale = subparsers.add_parser('scale', help=scale_usg, usage=scale_usg, description=scale_description,
                                         formatter_class=formatter_class)
    parser_scale._optionals.title = 'scale arguments'
    parser_scale.add_argument('kind', help='{deployment} object kind', choices=run_kinds, metavar="KIND")
    parser_scale.add_argument('name', help='object name to get info', metavar="NAME", type=str)
    parser_scale.add_argument('count', help='count of replicas', metavar="COUNT", type=int, choices=range(1, 10))
    parser_scale.add_argument('--namespace', '-n', help='namespace, default: \"default\"', required=False)

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