package pokemon

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/config"
)

func init() {
	enums := reflect.ValueOf(&TypeEnum)
	for i := range enums.Elem().NumField() {
		v := int32(1) << i
		field := enums.Elem().Field(i)
		if !field.CanSet() {
			continue
		}
		field.SetUint(uint64(v))
	}

	file, err := os.Open(filepath.Join(config.DataPath, "type_restraint_relationship.csv"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	header := data[0][1:]
	for _, line := range data[1:] {
		ht := parseChineseType(line[0])
		typeRestraintRelationship[ht] = make(map[Type]SkillEffect)
		for i, v := range line[1:] {
			vt := parseChineseType(header[i])
			switch strings.TrimSpace(v) {
			case "0":
				typeRestraintRelationship[ht][vt] = SkillEffectEnum.Null
			case "0.5":
				typeRestraintRelationship[ht][vt] = SkillEffectEnum.Bad
			case "1":
				typeRestraintRelationship[ht][vt] = SkillEffectEnum.Normal
			case "2":
				typeRestraintRelationship[ht][vt] = SkillEffectEnum.Excellent
			}
		}
	}
}

// 属性克制关系
var typeRestraintRelationship = make(map[Type]map[Type]SkillEffect)

// Type 属性
type Type uint32

var TypeEnum = enum.New[struct {
	Unknown  Type // ???
	None     Type // 无
	Normal   Type // 一般
	Flying   Type // 飞行
	Fire     Type // 火
	Psychic  Type // 超能力
	Water    Type // 水
	Bug      Type // 虫
	Electric Type // 电
	Rock     Type // 岩石
	Grass    Type // 草
	Ghost    Type // 幽灵
	Ice      Type // 冰
	Dragon   Type // 龙
	Fighting Type // 格斗
	Dark     Type // 恶
	Poison   Type // 毒
	Steel    Type // 钢
	Ground   Type // 地面
	Fairy    Type // 妖精
}]()

func parseChineseType(s string) Type {
	switch s {
	case "一般":
		return TypeEnum.Normal
	case "飞行":
		return TypeEnum.Flying
	case "火":
		return TypeEnum.Fire
	case "超能力":
		return TypeEnum.Psychic
	case "水":
		return TypeEnum.Water
	case "虫":
		return TypeEnum.Bug
	case "电":
		return TypeEnum.Electric
	case "岩石":
		return TypeEnum.Rock
	case "草":
		return TypeEnum.Grass
	case "幽灵":
		return TypeEnum.Ghost
	case "冰":
		return TypeEnum.Ice
	case "龙":
		return TypeEnum.Dragon
	case "格斗":
		return TypeEnum.Fighting
	case "恶":
		return TypeEnum.Dark
	case "毒":
		return TypeEnum.Poison
	case "钢":
		return TypeEnum.Steel
	case "地面":
		return TypeEnum.Ground
	case "妖精":
		return TypeEnum.Fairy
	default:
		return TypeEnum.Unknown
	}
}

// Contain 是否包含某属性
func (t Type) Contain(dst Type) bool {
	return t&dst == dst
}

// Flatten 取消聚合，扁平化
func (t Type) Flatten() []Type {
	allTypes := enum.Values[Type](TypeEnum)
	res := make([]Type, 0, len(allTypes))
	for _, dst := range allTypes {
		if !t.Contain(dst) {
			continue
		}
		res = append(res, dst)
	}
	return res
}

// GetEffectTo 当目标是指定属性时，获取效果
func (t Type) GetEffectTo(dst Type) float64 {
	fromType, toTypes := t.Flatten()[0], dst.Flatten()
	if len(toTypes) == 1 {
		// 如果目标是单属性，直接用属性相克表的值
		return typeRestraintRelationship[fromType][dst].Multiples()
	} else {
		// 如果目标不是单属性，将所有属性的值相乘
		v := float64(1)
		for _, dst = range toTypes {
			v *= typeRestraintRelationship[fromType][dst].Multiples()
		}
		return v
	}
}

// SkillEffect 技能效果
type SkillEffect uint8

var SkillEffectEnum = enum.New[struct {
	Null      SkillEffect // 没有效果
	Bad       SkillEffect // 效果不好
	Normal    SkillEffect // 效果一般
	Excellent SkillEffect // 效果绝佳
}]()

// Multiples 倍数
func (e SkillEffect) Multiples() float64 {
	switch e {
	case SkillEffectEnum.Bad:
		return 0.5
	case SkillEffectEnum.Normal:
		return 1
	case SkillEffectEnum.Excellent:
		return 2
	default:
		return 0
	}
}
