package plan

import "github.com/lindb/lindb/models"

type PlanFragment struct {
	ID FragmentID `json:"id"`

	Receivers []models.InternalNode `json:"receivers"`

	Root               PlanNode            `json:"root,omitempty"`
	RemoteParentNodeID *PlanNodeID         `json:"parentNode,omitempty"`
	RemoteSources      []*RemoteSourceNode `json:"remoteSources,omitempty"`

	Partitions map[models.InternalNode][]int `json:"-"`
}

func NewPlanFragment(id FragmentID, root PlanNode) *PlanFragment {
	fragment := &PlanFragment{
		ID:   id,
		Root: root,
	}

	fragment.findRemoteSources(root)
	return fragment
}

func (pf *PlanFragment) findRemoteSources(node PlanNode) {
	sources := node.GetSources()
	for i := range sources {
		pf.findRemoteSources(sources[i])
	}

	if remoteSource, ok := node.(*RemoteSourceNode); ok {
		pf.RemoteSources = append(pf.RemoteSources, remoteSource)
	}
}
