package character

import (
	"fmt"
	"log"
	"sort"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/avatar"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Config struct {
	PackageName    string   `yaml:"package_name,omitempty"`
	CharStructName string   `yaml:"char_struct_name,omitempty"`
	GenshinID      int32    `yaml:"genshin_id,omitempty"`
	SubID          int32    `yaml:"sub_id,omitempty"`
	Key            string   `yaml:"key,omitempty"`
	KeyVarName     string   `yaml:"key_var_name,omitempty"`
	Shortcuts      []string `yaml:"shortcuts,omitempty"`

	// skill data generation
	SkillDataMapping map[string]map[string][]int `yaml:"skill_data_mapping"`

	// param validation
	ActionParamKeys map[string][]paramData `yaml:"action_param_keys"`

	// icd groups
	ICDGroups []ICDGroupData `yaml:"icd_groups"`
	ICDTags   []string       `yaml:"icd_tags"`

	// doc related
	Documentation DocConfig `yaml:"documentation"`

	// extra fields to be populate but not read from yaml
	RelativePath string `yaml:"-"`
}

type DocConfig struct {
	FieldsData []FieldDocData `json:"fields_data" yaml:"fields_data"`
}

type FieldDocData struct {
	Fields []string `json:"fields" yaml:"fields"`
	Desc   string   `json:"desc"   yaml:"desc"`
}

type ICDGroupData struct {
	GroupName      string    `yaml:"group_name"`
	ResetTimer     int       `yaml:"reset_timer"`
	EleAppSequence []float64 `yaml:"ele_app_sequence"`
	DamageSequence []float64 `yaml:"damage_sequence"`
}

type paramData struct {
	Param string `yaml:"param,omitempty"`
	Desc  string `yaml:"desc,omitempty"`
}

type Generator struct {
	GeneratorConfig
	src   *avatar.DataSource
	chars []Config
	data  map[string]*model.AvatarData
}

type GeneratorConfig struct {
	Root   string
	Excels string
}

func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	g := &Generator{
		GeneratorConfig: cfg,
		data:            make(map[string]*model.AvatarData),
	}

	src, err := avatar.NewDataSource(g.Excels)
	if err != nil {
		return nil, err
	}
	g.src = src

	chars, err := ParseCharConfig(g.Root)
	if err != nil {
		return nil, err
	}
	g.chars = chars

	for i := range chars {
		v := chars[i]
		// validate char config
		if _, ok := g.data[v.Key]; ok {
			return nil, fmt.Errorf("duplicated key %v found; second instance at %v", v.Key, v.RelativePath)
		}
		char, err := src.GetAvatarData(v.GenshinID, v.SubID)
		if err != nil {
			log.Printf("Error loading %v data: %v; skipping\n", v.Key, err)
			continue
		}
		char.Key = v.Key
		g.data[v.Key] = char
	}

	return g, nil
}

func (g *Generator) Data() []*model.AvatarData {
	keys := make([]string, 0, len(g.data))
	for k := range g.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var res []*model.AvatarData
	for _, k := range keys {
		res = append(res, g.data[k])
	}
	return res
}
