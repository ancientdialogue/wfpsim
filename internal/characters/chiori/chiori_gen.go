// Code generated by "pipeline"; DO NOT EDIT.
package chiori

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
	1: {"hold"},
}

func init() {
	base = &model.AvatarData{}
	err := prototext.Unmarshal(pbData, base)
	if err != nil {
		panic(err)
	}
	validation.RegisterCharParamValidationFunc(keys.Chiori, ValidateParamKeys)
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
	attack = [][][]float64{
		{attack_1},
		{attack_2},
		attack_3,
		{attack_4},
	}
)

var (
	// attack: attack_1 = [0]
	attack_1 = []float64{
		0.4941039979457855,
		0.5343220233917236,
		0.5745400190353394,
		0.6319940090179443,
		0.6722120046615601,
		0.7181749939918518,
		0.7813739776611328,
		0.844573974609375,
		0.9077730178833008,
		0.976718008518219,
		1.0456629991531372,
		1.1146080493927002,
		1.1835520267486572,
		1.2524969577789307,
		1.3214420080184937,
	}
	// attack: attack_2 = [1]
	attack_2 = []float64{
		0.4683389961719513,
		0.5064589977264404,
		0.5445799827575684,
		0.5990380048751831,
		0.637158989906311,
		0.6807249784469604,
		0.7406290173530579,
		0.8005329966545105,
		0.8604360222816467,
		0.925786018371582,
		0.9911360144615173,
		1.0564850568771362,
		1.1218349933624268,
		1.1871839761734009,
		1.252534031867981,
	}
	// attack: attack_3 = [2 3]
	attack_3 = [][]float64{
		{
			0.3041650056838989,
			0.32892200350761414,
			0.3536800146102905,
			0.3890480101108551,
			0.4138059914112091,
			0.44209998846054077,
			0.4810050129890442,
			0.5199099779129028,
			0.5588139891624451,
			0.6012560129165649,
			0.64369797706604,
			0.6861389875411987,
			0.7285810112953186,
			0.7710220217704773,
			0.8134639859199524,
		},
		{
			0.3041650056838989,
			0.32892200350761414,
			0.3536800146102905,
			0.3890480101108551,
			0.4138059914112091,
			0.44209998846054077,
			0.4810050129890442,
			0.5199099779129028,
			0.5588139891624451,
			0.6012560129165649,
			0.64369797706604,
			0.6861389875411987,
			0.7285810112953186,
			0.7710220217704773,
			0.8134639859199524,
		},
	}
	// attack: attack_4 = [4]
	attack_4 = []float64{
		0.7512270212173462,
		0.8123739957809448,
		0.873520016670227,
		0.9608719944953918,
		1.0220179557800293,
		1.0918999910354614,
		1.187986969947815,
		1.2840739488601685,
		1.380162000656128,
		1.484984040260315,
		1.5898059606552124,
		1.6946289539337158,
		1.7994509935379028,
		1.9042739868164062,
		2.0090959072113037,
	}
	// attack: charge = [5 6]
	charge = [][]float64{
		{
			0.5430899858474731,
			0.5872949957847595,
			0.6315000057220459,
			0.6946499943733215,
			0.7388550043106079,
			0.7893750071525574,
			0.8588399887084961,
			0.9283049702644348,
			0.9977700114250183,
			1.0735499858856201,
			1.1493300199508667,
			1.2251100540161133,
			1.3008899688720703,
			1.376670002937317,
			1.4524500370025635,
		},
		{
			0.5430899858474731,
			0.5872949957847595,
			0.6315000057220459,
			0.6946499943733215,
			0.7388550043106079,
			0.7893750071525574,
			0.8588399887084961,
			0.9283049702644348,
			0.9977700114250183,
			1.0735499858856201,
			1.1493300199508667,
			1.2251100540161133,
			1.3008899688720703,
			1.376670002937317,
			1.4524500370025635,
		},
	}
	// attack: highPlunge = [10]
	highPlunge = []float64{
		1.59676194190979,
		1.7267309427261353,
		1.8566999435424805,
		2.042370080947876,
		2.1723389625549316,
		2.3208749294281006,
		2.5251119136810303,
		2.72934889793396,
		2.9335858821868896,
		3.1563899517059326,
		3.3791940212249756,
		3.6019980907440186,
		3.8248019218444824,
		4.047605991363525,
		4.270410060882568,
	}
	// attack: lowPlunge = [9]
	lowPlunge = []float64{
		1.2783770561218262,
		1.3824310302734375,
		1.4864850044250488,
		1.635133981704712,
		1.7391870021820068,
		1.858106017112732,
		2.021620035171509,
		2.1851329803466797,
		2.3486459255218506,
		2.527024984359741,
		2.7054030895233154,
		2.8837809562683105,
		3.0621590614318848,
		3.24053692817688,
		3.418915033340454,
	}
	// attack: plunge = [8]
	plunge = []float64{
		0.6393240094184875,
		0.6913620233535767,
		0.743399977684021,
		0.8177400231361389,
		0.8697779774665833,
		0.9292500019073486,
		1.011023998260498,
		1.0927979946136475,
		1.1745719909667969,
		1.2637799978256226,
		1.3529880046844482,
		1.442196011543274,
		1.5314040184020996,
		1.6206120252609253,
		1.709820032119751,
	}
	// skill: thrustAtkScaling = [4]
	thrustAtkScaling = []float64{
		1.492799997329712,
		1.6047600507736206,
		1.7167199850082397,
		1.8660000562667847,
		1.9779599905014038,
		2.0899200439453125,
		2.2392001152038574,
		2.3884799480438232,
		2.537760019302368,
		2.687040090560913,
		2.836319923400879,
		2.985599994659424,
		3.1721999645233154,
		3.358799934387207,
		3.5453999042510986,
	}
	// skill: thrustDefScaling = [5]
	thrustDefScaling = []float64{
		1.8660000562667847,
		2.0059499740600586,
		2.145900011062622,
		2.3324999809265137,
		2.472450017929077,
		2.6124000549316406,
		2.7990000247955322,
		2.985599994659424,
		3.1721999645233154,
		3.358799934387207,
		3.5453999042510986,
		3.7320001125335693,
		3.965250015258789,
		4.198500156402588,
		4.4317498207092285,
	}
	// skill: turretAtkScaling = [0]
	turretAtkScaling = []float64{
		0.8208000063896179,
		0.8823599815368652,
		0.9439200162887573,
		1.0260000228881836,
		1.0875600576400757,
		1.1491199731826782,
		1.2311999797821045,
		1.3132799863815308,
		1.395359992980957,
		1.4774399995803833,
		1.5595200061798096,
		1.6416000127792358,
		1.7441999912261963,
		1.8467999696731567,
		1.9493999481201172,
	}
	// skill: turretDefScaling = [1]
	turretDefScaling = []float64{
		1.0260000228881836,
		1.1029499769210815,
		1.179900050163269,
		1.2825000286102295,
		1.3594499826431274,
		1.436400055885315,
		1.5390000343322754,
		1.6416000127792358,
		1.7441999912261963,
		1.8467999696731567,
		1.9493999481201172,
		2.052000045776367,
		2.1802499294281006,
		2.308500051498413,
		2.4367499351501465,
	}
	// skill: turretDuration = [2]
	turretDuration = []float64{
		17,
		17,
		17,
		17,
		17,
		17,
		17,
		17,
		17,
		17,
		17,
		17,
		17,
		17,
		17,
	}
	// skill: turretInterval = [3]
	turretInterval = []float64{
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
		3.5999999046325684,
	}
	// burst: burstAtkScaling = [0]
	burstAtkScaling = []float64{
		2.563199996948242,
		2.7554399967193604,
		2.9476799964904785,
		3.2039999961853027,
		3.396239995956421,
		3.588479995727539,
		3.8447999954223633,
		4.1011199951171875,
		4.357439994812012,
		4.613759994506836,
		4.87007999420166,
		5.126399993896484,
		5.446800231933594,
		5.767199993133545,
		6.087600231170654,
	}
	// burst: burstCD = [2]
	burstCD = []float64{
		13.5,
		13.5,
		13.5,
		13.5,
		13.5,
		13.5,
		13.5,
		13.5,
		13.5,
		13.5,
		13.5,
		13.5,
		13.5,
		13.5,
		13.5,
	}
	// burst: burstDefScaling = [1]
	burstDefScaling = []float64{
		3.2039999961853027,
		3.4442999362945557,
		3.6846001148223877,
		4.005000114440918,
		4.245299816131592,
		4.485599994659424,
		4.806000232696533,
		5.126399993896484,
		5.446800231933594,
		5.767199993133545,
		6.087600231170654,
		6.4079999923706055,
		6.808499813079834,
		7.209000110626221,
		7.609499931335449,
	}
)