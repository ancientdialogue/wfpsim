// Code generated by "pipeline"; DO NOT EDIT.
package mualani

import (
	_ "embed"

	"fmt"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/validation"
	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/encoding/prototext"
	"slices"
)

//go:embed data_gen.textproto
var pbData []byte
var base *model.AvatarData
var paramKeysValidation = map[action.Action][]string{
	2: {"travel"},
}

func init() {
	base = &model.AvatarData{}
	err := prototext.Unmarshal(pbData, base)
	if err != nil {
		panic(err)
	}
	validation.RegisterCharParamValidationFunc(keys.Mualani, ValidateParamKeys)
}

func ValidateParamKeys(a action.Action, keys []string) error {
	valid, ok := paramKeysValidation[a]
	if !ok {
		return nil
	}
	for _, v := range keys {
		if !slices.Contains(valid, v) {
			return fmt.Errorf("key %v is invalid for action %v", v, a.String())
		}
	}
	return nil
}

func (x *char) Data() *model.AvatarData {
	return base
}

var (
	attack = [][]float64{
		attack_1,
		attack_2,
		attack_3,
	}
)

var (
	// attack: attack_1 = [0]
	attack_1 = []float64{
		0.5139600038528,
		0.5525069832802,
		0.5910540223122,
		0.6424499750137,
		0.6809970140457,
		0.7195439934731,
		0.7709400057793,
		0.8223360180855,
		0.873731970787,
		0.9251279830933,
		0.9765239953995,
		1.0279200077057,
		1.0921649932861,
		1.1564099788666,
		1.220654964447,
	}
	// attack: attack_2 = [1]
	attack_2 = []float64{
		0.4462560117245,
		0.4797250032425,
		0.5131940245628,
		0.5578200221062,
		0.5912889838219,
		0.6247580051422,
		0.6693840026855,
		0.7140100002289,
		0.7586349844933,
		0.8032609820366,
		0.8478860259056,
		0.8925120234489,
		0.9482939839363,
		1.0040760040283,
		1.0598579645157,
	}
	// attack: attack_3 = [2]
	attack_3 = []float64{
		0.7003440260887,
		0.7528700232506,
		0.8053960204124,
		0.8754299879074,
		0.9279559850693,
		0.9804819822311,
		1.0505160093307,
		1.1205500364304,
		1.1905850172043,
		1.2606190443039,
		1.3306540250778,
		1.4006880521774,
		1.4882309436798,
		1.5757739543915,
		1.6633169651031,
	}
	// attack: charge = [3]
	charge = []float64{
		1.4287999868393,
		1.5359599590302,
		1.6431200504303,
		1.7860000133514,
		1.8931599855423,
		2.0003199577332,
		2.1431999206543,
		2.2860798835754,
		2.4289600849152,
		2.5718400478363,
		2.7147200107574,
		2.8575999736786,
		3.0362000465393,
		3.2147998809814,
		3.3933999538422,
	}
	// bite: bite = [0]
	bite = []float64{
		0.0868000015616,
		0.0933099985123,
		0.0998200029135,
		0.1084999963641,
		0.1150100007653,
		0.121519997716,
		0.1301999986172,
		0.1388799995184,
		0.1475600004196,
		0.1562400013208,
		0.1649200022221,
		0.1736000031233,
		0.1844500005245,
		0.1952999979258,
		0.206149995327,
	}
	// skill: momentumBonus = [1]
	momentumBonus = []float64{
		0.0434000007808,
		0.0466549992561,
		0.0499100014567,
		0.0542499981821,
		0.0575050003827,
		0.060759998858,
		0.0650999993086,
		0.0694399997592,
		0.0737800002098,
		0.0781200006604,
		0.082460001111,
		0.0868000015616,
		0.0922250002623,
		0.0976499989629,
		0.1030749976635,
	}
	// skill: surgingBite = [2]
	surgingBite = []float64{
		0.2169999927282,
		0.2332749962807,
		0.2495499998331,
		0.2712500095367,
		0.287524998188,
		0.3037999868393,
		0.3255000114441,
		0.3472000062466,
		0.368900001049,
		0.3905999958515,
		0.412299990654,
		0.4339999854565,
		0.4611249864101,
		0.4882499873638,
		0.5153750181198,
	}
	// burst: burst = [0]
	burst = []float64{
		0.5843920111656,
		0.628220975399,
		0.672051012516,
		0.7304900288582,
		0.7743189930916,
		0.8181489706039,
		0.8765879869461,
		0.9350270032883,
		0.9934660196304,
		1.0519059896469,
		1.1103450059891,
		1.1687840223312,
		1.2418329715729,
		1.3148820400238,
		1.3879309892654,
	}

)
