package gojira

type ProjectsMeta struct {
	Projects map[string]ProjectMeta
}

func NewProjectsMeta() ProjectsMeta {
	pm := ProjectsMeta{}
	pm.init()
	return pm
}

func (pm *ProjectsMeta) AddMap(info map[string]float32) {
	pm.init()
	for key, teamSize := range info {
		proj, ok := pm.Projects[key]
		if !ok {
			proj = ProjectMeta{Key: key}
		}
		proj.TeamSize = teamSize
		pm.Projects[key] = proj
	}
}

func (pm *ProjectsMeta) init() {
	if pm.Projects == nil {
		pm.Projects = map[string]ProjectMeta{}
	}
}

type ProjectMeta struct {
	Key      string
	TeamSize float32
}
