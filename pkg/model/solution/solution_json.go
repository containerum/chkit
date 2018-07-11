package solution

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

var (
	_ model.JSONrenderer = Solution{}
	_ json.Marshaler     = Solution{}
)

func (solution Solution) RenderJSON() (string, error) {
	data, err := solution.MarshalJSON()
	return string(data), err
}

func (solution Solution) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(solution.ToKube(), "", model.Indent)
}

func (solution *Solution) UnmarshalJSON(data []byte) error {
	var kubeSol kubeModels.UserSolution
	if err := json.Unmarshal(data, &kubeSol); err != nil {
		return err
	}
	*solution = SolutionFromKube(kubeSol)
	return nil
}
