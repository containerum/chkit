package ingress

import kubeModels "git.containerum.net/ch/kube-client/pkg/model"

type Path kubeModels.Path

func PathFromKube(kubePath kubeModels.Path) Path {
	return Path(kubePath)
}

func (path Path) ToKube() kubeModels.Path {
	return kubeModels.Path(path)
}

type PathList []Path

func PathListFromKube(kubeList []kubeModels.Path) PathList {
	var list PathList = make([]Path, 0, len(kubeList))
	for _, p := range kubeList {
		list = append(list, PathFromKube(p))
	}
	return list
}

func (list PathList) Copy() PathList {
	return append(make([]Path, 0, len(list)), list...)
}

func (list *PathList) Delete(i int) PathList {
	cp := list.Copy()
	return append(cp[:i], cp[i+1:]...)
}

func (list PathList) Append(paths ...Path) PathList {
	return append(list.Copy(), paths...)
}

func (list PathList) ToKube() []kubeModels.Path {
	kubeList := make([]kubeModels.Path, 0, len(list))
	for _, path := range list {
		kubeList = append(kubeList, path.ToKube())
	}
	return kubeList
}
