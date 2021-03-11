package tehai

import (
	"mahjong/model/hai"
	"testing"

	"github.com/stretchr/testify/assert"
)

func copy(hais []*hai.Hai) []*hai.Hai {
	return append([]*hai.Hai{}, hais...)
}

func TestAdd(t *testing.T) {
	cases := []struct {
		name       string
		beforeHais []*hai.Hai
		inHai      *hai.Hai
		afterHais  []*hai.Hai
		outError   error
	}{
		{
			name:       "success",
			beforeHais: []*hai.Hai{},
			inHai:      hai.Haku,
			afterHais:  []*hai.Hai{hai.Haku},
			outError:   nil,
		},
		{
			name:       "failure",
			beforeHais: append(copy(hai.Manzu), hai.Souzu[:4]...),
			inHai:      nil,
			outError:   TehaiHaiIsNilErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tehai := tehaiImpl{hais: c.beforeHais}
			err := tehai.Add(c.inHai)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterHais, tehai.hais)
		})
	}
}

func TestAdds(t *testing.T) {
	cases := []struct {
		name       string
		beforeHais []*hai.Hai
		inHais     []*hai.Hai
		afterHais  []*hai.Hai
		outError   error
	}{
		{
			name:       "success",
			beforeHais: copy(hai.Manzu),
			inHais:     hai.Pinzu[:3],
			afterHais:  append(copy(hai.Manzu), hai.Pinzu[:3]...),
			outError:   nil,
		},
		{
			name:       "failure",
			beforeHais: copy(hai.Manzu),
			inHais:     []*hai.Hai{nil},
			outError:   TehaiHaiIsNilErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tehai := tehaiImpl{hais: c.beforeHais}
			err := tehai.Adds(c.inHais)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterHais, tehai.hais)
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		name       string
		beforeHais []*hai.Hai
		inHai      *hai.Hai
		afterHais  []*hai.Hai
		outHai     *hai.Hai
		outError   error
	}{
		{
			name:       "success",
			beforeHais: copy(hai.Manzu),
			inHai:      hai.Manzu9,
			afterHais:  hai.Manzu[:8],
			outHai:     hai.Manzu9,
			outError:   nil,
		},
		{
			name:       "failure",
			beforeHais: copy(hai.Manzu),
			inHai:      hai.Pinzu1,
			outError:   TehaiHaiNotFoundErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tehai := tehaiImpl{hais: c.beforeHais}
			outHai, err := tehai.Remove(c.inHai)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterHais, tehai.hais)
			assert.Equal(t, c.outHai, outHai)

		})
	}
}

func TestRemoves(t *testing.T) {
	cases := []struct {
		name       string
		beforeHais []*hai.Hai
		inHais     []*hai.Hai
		afterHais  []*hai.Hai
		outHais    []*hai.Hai
		outError   error
	}{
		{
			name:       "success",
			beforeHais: copy(hai.Manzu),
			inHais:     hai.Manzu[:3],
			afterHais:  hai.Manzu[3:],
			outHais:    hai.Manzu[:3],
			outError:   nil,
		},
		{
			name:       "failure",
			beforeHais: copy(hai.Manzu),
			inHais:     hai.Pinzu[:2],
			outError:   TehaiHaiNotFoundErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tehai := tehaiImpl{hais: c.beforeHais}
			outHais, err := tehai.Removes(c.inHais)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterHais, tehai.hais)
			assert.Equal(t, c.outHais, outHais)
		})
	}
}

func TestReplace(t *testing.T) {
	cases := []struct {
		name       string
		beforeHais []*hai.Hai
		inInHai    *hai.Hai
		inOutHai   *hai.Hai
		afterHais  []*hai.Hai
		outHai     *hai.Hai
		outError   error
	}{
		{
			name:       "success",
			beforeHais: copy(hai.Manzu),
			inInHai:    hai.Souzu1,
			inOutHai:   hai.Manzu9,
			afterHais:  append(copy(hai.Manzu[:8]), hai.Souzu1),
			outHai:     hai.Manzu9,
			outError:   nil,
		},
		{
			name:       "failure",
			beforeHais: copy(hai.Manzu),
			inOutHai:   hai.Souzu1,
			outError:   TehaiHaiNotFoundErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tehai := tehaiImpl{hais: c.beforeHais}
			outHai, err := tehai.Replace(c.inInHai, c.inOutHai)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterHais, tehai.hais)
			assert.Equal(t, c.outHai, outHai)

		})
	}
}

func TestPonPairs(t *testing.T) {
	cases := []struct {
		name       string
		beforeHais []*hai.Hai
		inHai      *hai.Hai
		outPairs   [][2]*hai.Hai
	}{
		{
			name:       "success",
			beforeHais: []*hai.Hai{hai.Haku, hai.Haku},
			inHai:      hai.Haku,
			outPairs:   [][2]*hai.Hai{{hai.Haku, hai.Haku}},
		},
		{
			name:       "failureHaku1",
			beforeHais: copy(hai.Manzu),
			inHai:      hai.Haku,
			outPairs:   [][2]*hai.Hai{},
		},
		{
			name:       "failure2",
			beforeHais: []*hai.Hai{hai.Haku, hai.Haku},
			inHai:      hai.Hatsu,
			outPairs:   [][2]*hai.Hai{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tehai := tehaiImpl{hais: c.beforeHais}
			pairs, err := tehai.PonPairs(c.inHai)
			if err != nil {
				assert.Equal(t, nil, err)
				return
			}
			assert.Equal(t, c.outPairs, pairs)
		})
	}
}

func TestKanPairs(t *testing.T) {
	cases := []struct {
		name       string
		beforeHais []*hai.Hai
		inHai      *hai.Hai
		outPairs   [][3]*hai.Hai
	}{
		{
			name:       "success",
			beforeHais: []*hai.Hai{hai.Haku, hai.Haku, hai.Haku},
			inHai:      hai.Haku,
			outPairs:   [][3]*hai.Hai{{hai.Haku, hai.Haku, hai.Haku}},
		},
		{
			name:       "failure1",
			beforeHais: []*hai.Hai{},
			inHai:      hai.Haku,
			outPairs:   [][3]*hai.Hai{},
		},
		{
			name:       "failure2",
			beforeHais: []*hai.Hai{hai.Haku, hai.Haku, hai.Haku},
			inHai:      hai.Hatsu,
			outPairs:   [][3]*hai.Hai{},
		},
	}

	for _, c := range cases {
		tehai := tehaiImpl{hais: c.beforeHais}
		pairs, err := tehai.MinKanPairs(c.inHai)
		if err != nil {
			assert.Equal(t, nil, err)
			continue
		}
		assert.Equal(t, c.outPairs, pairs)
	}
}

func TestChiiPairs(t *testing.T) {
	cases := []struct {
		name       string
		beforeHais []*hai.Hai
		inHai      *hai.Hai
		outPairs   [][2]*hai.Hai
	}{
		{
			name:       "success",
			beforeHais: copy(hai.Manzu),
			inHai:      hai.Manzu4,
			outPairs:   [][2]*hai.Hai{{hai.Manzu2, hai.Manzu3}, {hai.Manzu3, hai.Manzu5}, {hai.Manzu5, hai.Manzu6}},
		},
		{
			name:       "failure1",
			beforeHais: copy(hai.Manzu),
			inHai:      hai.Pinzu1,
			outPairs:   [][2]*hai.Hai{},
		},
		{
			name:       "failure2",
			beforeHais: copy(hai.Manzu),
			inHai:      hai.Haku,
			outPairs:   [][2]*hai.Hai{},
		},
	}

	for _, c := range cases {
		tehai := tehaiImpl{hais: c.beforeHais}
		pairs, err := tehai.ChiiPairs(c.inHai)
		if err != nil {
			assert.Equal(t, nil, err)
			continue
		}
		assert.Equal(t, c.outPairs, pairs)
	}
}

func TestMachihai(t *testing.T) {
	cases := []struct {
		name       string
		beforeHais []*hai.Hai
		outHais    []*hai.Hai
	}{
		{
			name: "両面",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu2, hai.Manzu3, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Manzu7, hai.Manzu8, hai.Manzu9, hai.Haku, hai.Pinzu2, hai.Pinzu3,
				hai.Haku,
			},
			outHais: []*hai.Hai{hai.Pinzu1, hai.Pinzu4},
		},
		{
			name: "辺張",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu2, hai.Manzu3, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu8, hai.Pinzu9, hai.Souzu6, hai.Souzu6, hai.Souzu6, hai.Souzu7,
				hai.Souzu7,
			},
			outHais: []*hai.Hai{hai.Pinzu7},
		},
		{
			name: "嵌張",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Manzu7, hai.Manzu8, hai.Manzu9, hai.Haku, hai.Pinzu5, hai.Pinzu3,
				hai.Haku,
			},
			outHais: []*hai.Hai{hai.Pinzu4},
		},
		{
			name: "双碰",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Manzu7, hai.Manzu8, hai.Manzu9, hai.Hatsu, hai.Pinzu3, hai.Pinzu3,
				hai.Hatsu,
			},
			outHais: []*hai.Hai{hai.Pinzu3, hai.Hatsu},
		},
		{
			name: "単騎",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Manzu7, hai.Manzu8, hai.Manzu9, hai.Pinzu1, hai.Pinzu2, hai.Pinzu3,
				hai.Hatsu,
			},
			outHais: []*hai.Hai{hai.Hatsu},
		},
		{
			name: "延段",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Manzu7, hai.Manzu8, hai.Manzu9, hai.Pinzu1, hai.Pinzu2, hai.Pinzu3,
				hai.Pinzu4,
			},
			outHais: []*hai.Hai{hai.Pinzu1, hai.Pinzu4},
		},
		{
			name: "1/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu2, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6, hai.Pei,
				hai.Pei,
			},
			outHais: []*hai.Hai{hai.Pinzu1, hai.Pinzu4, hai.Pinzu7},
		},
		{
			name: "2/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu2, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6, hai.Pinzu7,
				hai.Pinzu8,
			},
			outHais: []*hai.Hai{hai.Pinzu2, hai.Pinzu5, hai.Pinzu8},
		},
		{
			name: "3/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu2, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6, hai.Pinzu7,
				hai.Pinzu5,
			},
			outHais: []*hai.Hai{hai.Pinzu2, hai.Pinzu5, hai.Pinzu8},
		},
		{
			name: "4/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu2, hai.Pinzu3, hai.Pinzu4, hai.Pinzu4, hai.Pinzu4, hai.Pei,
				hai.Pei,
			},
			outHais: []*hai.Hai{hai.Pinzu1, hai.Pinzu4, hai.Pei},
		},
		{
			name: "5/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu2, hai.Pinzu3, hai.Pinzu3, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5,
				hai.Pinzu6,
			},
			outHais: []*hai.Hai{hai.Pinzu1, hai.Pinzu2, hai.Pinzu4, hai.Pinzu7},
		},
		{
			name: "6/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu3, hai.Pinzu4, hai.Pinzu4, hai.Pinzu5,
				hai.Pinzu6,
			},
			outHais: []*hai.Hai{hai.Pinzu2, hai.Pinzu4, hai.Pinzu5, hai.Pinzu7},
		},
		{
			name: "7/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu5,
				hai.Pinzu6,
			},
			outHais: []*hai.Hai{hai.Pinzu4, hai.Pinzu5, hai.Pinzu7},
		},
		{
			name: "8/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6,
				hai.Pinzu7,
			},
			outHais: []*hai.Hai{hai.Pinzu2, hai.Pinzu4, hai.Pinzu5, hai.Pinzu7, hai.Pinzu8},
		},
		{
			name: "9/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6,
				hai.Pinzu8,
			},
			outHais: []*hai.Hai{hai.Pinzu7, hai.Pinzu8},
		},
		{
			name: "10/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu3, hai.Pinzu5, hai.Pinzu5, hai.Pinzu6,
				hai.Pinzu7,
			},
			outHais: []*hai.Hai{hai.Pinzu4, hai.Pinzu5, hai.Pinzu8},
		},
		{
			name: "11/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu3, hai.Pinzu5, hai.Pinzu6, hai.Pinzu7,
				hai.Pinzu8,
			},
			outHais: []*hai.Hai{hai.Pinzu4, hai.Pinzu5, hai.Pinzu8},
		},
		{
			name: "12/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu3, hai.Pinzu4, hai.Pinzu4, hai.Pinzu5,
				hai.Pinzu5,
			},
			outHais: []*hai.Hai{hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6},
		},
		{
			name: "13/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu4, hai.Pinzu4, hai.Pinzu4, hai.Pinzu5,
				hai.Pinzu5,
			},
			outHais: []*hai.Hai{hai.Pinzu3, hai.Pinzu4, hai.Pinzu5},
		},
		{
			name: "14/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu5,
				hai.Pinzu5,
			},
			outHais: []*hai.Hai{hai.Pinzu2, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6},
		},
		{
			name: "15/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu3, hai.Pinzu4, hai.Pinzu4, hai.Pinzu4,
				hai.Pinzu5,
			},
			outHais: []*hai.Hai{hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6},
		},
		{
			name: "16/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu3, hai.Pinzu5, hai.Pinzu7, hai.Pinzu7,
				hai.Pinzu7,
			},
			outHais: []*hai.Hai{hai.Pinzu4, hai.Pinzu5, hai.Pinzu6},
		},
		{
			name: "17/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu5, hai.Pinzu5,
				hai.Pinzu5,
			},
			outHais: []*hai.Hai{hai.Pinzu3, hai.Pinzu4, hai.Pinzu6},
		},
		{
			name: "18/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu3, hai.Pinzu4, hai.Pinzu4, hai.Pinzu4, hai.Pinzu4,
				hai.Pinzu5,
			},
			outHais: []*hai.Hai{hai.Pinzu2, hai.Pinzu3, hai.Pinzu5, hai.Pinzu6},
		},
		{
			name: "19/19",
			beforeHais: []*hai.Hai{
				hai.Manzu1, hai.Manzu1, hai.Manzu1, hai.Manzu4, hai.Manzu5, hai.Manzu6,
				hai.Pinzu3, hai.Pinzu4, hai.Pinzu4, hai.Pinzu4, hai.Pinzu4, hai.Pinzu5,
				hai.Pinzu6,
			},
			outHais: []*hai.Hai{hai.Pinzu2, hai.Pinzu3, hai.Pinzu5, hai.Pinzu6},
		},
	}

	for _, c := range cases {

		t.Run(c.name, func(t *testing.T) {
			tehai := tehaiImpl{c.beforeHais}
			hais, err := tehai.Machihai()
			assert.NoError(t, err)
			assert.Equal(t, c.outHais, hais)
		})
	}

}
