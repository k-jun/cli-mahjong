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
		beforeHais []*hai.Hai
		inHai      *hai.Hai
		afterHais  []*hai.Hai
		outError   error
	}{
		{
			beforeHais: []*hai.Hai{},
			inHai:      &hai.Haku,
			afterHais:  []*hai.Hai{&hai.Haku},
			outError:   nil,
		},
		{
			beforeHais: append(copy(hai.Manzu), hai.Souzu[:4]...),
			inHai:      &hai.Haku,
			outError:   TehaiReachMaxHaiErr,
		},
	}

	for _, c := range cases {
		tehai := tehaiImpl{hais: c.beforeHais}
		err := tehai.Add(c.inHai)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterHais, tehai.hais)
	}
}

func TestAdds(t *testing.T) {
	cases := []struct {
		beforeHais []*hai.Hai
		inHais     []*hai.Hai
		afterHais  []*hai.Hai
		outError   error
	}{
		{
			beforeHais: copy(hai.Manzu),
			inHais:     hai.Pinzu[:3],
			afterHais:  append(copy(hai.Manzu), hai.Pinzu[:3]...),
			outError:   nil,
		},
		{
			beforeHais: copy(hai.Manzu),
			inHais:     hai.Pinzu[:4],
			outError:   TehaiReachMaxHaiErr,
		},
	}

	for _, c := range cases {
		tehai := tehaiImpl{hais: c.beforeHais}
		err := tehai.Adds(c.inHais)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterHais, tehai.hais)
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		beforeHais []*hai.Hai
		inHai      *hai.Hai
		afterHais  []*hai.Hai
		outHai     *hai.Hai
		outError   error
	}{
		{
			beforeHais: copy(hai.Manzu),
			inHai:      &hai.Manzu9,
			afterHais:  hai.Manzu[:9],
			outHai:     &hai.Manzu9,
			outError:   nil,
		},
		{
			beforeHais: copy(hai.Manzu),
			inHai:      &hai.Pinzu1,
			outError:   TehaiHaiNotFoundErr,
		},
	}

	for _, c := range cases {
		tehai := tehaiImpl{hais: c.beforeHais}
		outHai, err := tehai.Remove(c.inHai)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterHais, tehai.hais)
		assert.Equal(t, c.outHai, outHai)
	}
}

func TestRemoves(t *testing.T) {
	cases := []struct {
		beforeHais []*hai.Hai
		inHais     []*hai.Hai
		afterHais  []*hai.Hai
		outHais    []*hai.Hai
		outError   error
	}{
		{
			beforeHais: copy(hai.Manzu),
			inHais:     hai.Manzu[:3],
			afterHais:  hai.Manzu[3:],
			outHais:    hai.Manzu[:3],
			outError:   nil,
		},
		{
			beforeHais: copy(hai.Manzu),
			inHais:     hai.Pinzu[:2],
			outError:   TehaiHaiNotFoundErr,
		},
	}

	for _, c := range cases {
		tehai := tehaiImpl{hais: c.beforeHais}
		outHais, err := tehai.Removes(c.inHais)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterHais, tehai.hais)
		assert.Equal(t, c.outHais, outHais)
	}
}

func TestReplace(t *testing.T) {
	cases := []struct {
		beforeHais []*hai.Hai
		inInHai    *hai.Hai
		inOutHai   *hai.Hai
		afterHais  []*hai.Hai
		outHai     *hai.Hai
		outError   error
	}{
		{
			beforeHais: copy(hai.Manzu),
			inInHai:    &hai.Souzu1,
			inOutHai:   &hai.Manzu9,
			afterHais:  append(copy(hai.Manzu[:9]), &hai.Souzu1),
			outHai:     &hai.Manzu9,
			outError:   nil,
		},
		{
			beforeHais: copy(hai.Manzu),
			inOutHai:   &hai.Souzu1,
			outError:   TehaiHaiNotFoundErr,
		},
	}

	for _, c := range cases {
		tehai := tehaiImpl{hais: c.beforeHais}
		outHai, err := tehai.Replace(c.inInHai, c.inOutHai)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterHais, tehai.hais)
		assert.Equal(t, c.outHai, outHai)
	}
}

func TestFindPonPairs(t *testing.T) {
	cases := []struct {
		beforeHais []*hai.Hai
		inHai      *hai.Hai
		outPairs   [][2]*hai.Hai
	}{
		{
			beforeHais: []*hai.Hai{&hai.Haku, &hai.Haku},
			inHai:      &hai.Haku,
			outPairs:   [][2]*hai.Hai{{&hai.Haku, &hai.Haku}},
		},
		{
			beforeHais: []*hai.Hai{},
			inHai:      &hai.Haku,
			outPairs:   [][2]*hai.Hai{},
		},
		{
			beforeHais: copy(hai.Manzu),
			inHai:      &hai.Haku,
			outPairs:   [][2]*hai.Hai{},
		},
	}

	for _, c := range cases {
		tehai := tehaiImpl{hais: c.beforeHais}
		pairs := tehai.FindPonPairs(c.inHai)
		assert.Equal(t, c.outPairs, pairs)
	}
}

func TestFindKanPairs(t *testing.T) {
	cases := []struct {
		beforeHais []*hai.Hai
		inHai      *hai.Hai
		outPairs   [][3]*hai.Hai
	}{
		{
			beforeHais: []*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku},
			inHai:      &hai.Haku,
			outPairs:   [][3]*hai.Hai{{&hai.Haku, &hai.Haku, &hai.Haku}},
		},
		{
			beforeHais: []*hai.Hai{},
			inHai:      &hai.Haku,
			outPairs:   [][3]*hai.Hai{},
		},
		{
			beforeHais: copy(hai.Manzu),
			inHai:      &hai.Haku,
			outPairs:   [][3]*hai.Hai{},
		},
	}

	for _, c := range cases {
		tehai := tehaiImpl{hais: c.beforeHais}
		pairs := tehai.FindKanPairs(c.inHai)
		assert.Equal(t, c.outPairs, pairs)
	}
}

func TestFindChiPairs(t *testing.T) {
	cases := []struct {
		beforeHais []*hai.Hai
		inHai      *hai.Hai
		outPairs   [][2]*hai.Hai
	}{
		{
			beforeHais: copy(hai.Manzu),
			inHai:      &hai.Manzu4,
			outPairs:   [][2]*hai.Hai{{&hai.Manzu2, &hai.Manzu3}, {&hai.Manzu3, &hai.Manzu5}, {&hai.Manzu5, &hai.Manzu6}},
		},
		{
			beforeHais: copy(hai.Manzu),
			inHai:      &hai.Pinzu1,
			outPairs:   [][2]*hai.Hai{},
		},
		{
			beforeHais: copy(hai.Manzu),
			inHai:      &hai.Haku,
			outPairs:   [][2]*hai.Hai{},
		},
	}

	for _, c := range cases {
		tehai := tehaiImpl{hais: c.beforeHais}
		pairs := tehai.FindChiPairs(c.inHai)
		assert.Equal(t, c.outPairs, pairs)
	}
}
