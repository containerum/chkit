package cmd

import (
	"fmt"
	"os"

	jww "github.com/spf13/jwalterweatherman"
)

var np *jww.Notepad

var bashCompletionFunc = fmt.Sprintf(`__chkit_get_outformat()
{
	COMPREPLY=( $(compgen -W "json yaml pretty" -- "${cur}") )
}

__chkit_get_object_list()
{
	local prog
	prog="${COMP_WORDS[0]}"
	ret=$(${prog} get $1 | cut -d '|' -f 2 | grep -vE '(^\+)|(NAME)')
	code=$?
	echo "${ret}"
	return ${code}
}

__chkit_containers_in_deploy()
{
	local prog
	prog="${COMP_WORDS[0]}"
	ret=$(${prog} get deploy $1 | sed -e '1,/Containers:/d' | sed '/Ports/,/ImagePullPolicy/d' | sed 's/\t//g')
	code=$?
	echo "${ret}"
	return ${code}
}

__chkit_namespaces_list()
{
	__chkit_get_object_list namespace
}

__chkit_get_sort_columns()
{
	local cur cmd cols
	cur="${COMP_WORDS[COMP_CWORD]}"
	cmd="${COMP_WORDS[2]}"
	case "${cmd}" in
		"deployments" | "deployment" | "deploy" )
			cols="%s"
		;;
		"po" | "pods" | "pod" )
			cols="%s"
		;;
		"service" | "services" | "svc" )
			cols="%s"
		;;
		"ns" | "namespaces" | "namespace" )
			cols="%s"
		;;
	esac
	COMPREPLY=( $(compgen -W "${cols}" -- "${cur}") )
}

__custom_func()
{
	local cur prev list
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"
	case "${prev}" in
		"po" | "pods" | "pod" | "deployments" | "deployment" | "deploy" | "service" | "services" | "svc" | "ns" | "namespaces" | "namespace" )
			list="$(__chkit_get_object_list ${prev})"
			if [[ $? == 0 ]]; then
				COMPREPLY=( $(compgen -W "${list}" -- "${cur}") )
			fi
		;;
		* )
			if [[ ${last_command} == "chkit_set" && $COMP_CWORD == 4 ]]; then
				list="$(__chkit_containers_in_deploy ${COMP_WORDS[3]})"
				if [[ $? == 0 ]]; then
					COMPREPLY=( $(compgen -W "${list}" -- "${cur}") )
				fi
			fi
		;;
	esac
}`)

func exitOnErr(err error) {
	if err != nil {
		np.ERROR.Println(err)
		os.Exit(1)
	}
}
