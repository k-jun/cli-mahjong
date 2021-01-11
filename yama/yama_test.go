package yama

import (
	"mahjong/hai"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_ = New()
}

func TestTumo(t *testing.T) {
	cases := []struct {
		beforeYamaHai []*hai.Hai
		beforeWanHai  []*hai.Hai
		outHai        *hai.Hai
		outError      error
	}{
		{
			beforeYamaHai: all[:122],
			beforeWanHai:  all[122:],
			outHai:        &hai.Manzu1,
			outError:      nil,
		},
		{
			beforeYamaHai: []*hai.Hai{},
			beforeWanHai:  all[122:],
			outHai:        nil,
			outError:      YamaNoMoreHaiErr,
		},
	}

	for _, c := range cases {
		yama := yamaImpl{yamaHai: c.beforeYamaHai, wanHai: c.beforeWanHai}
		outHai, err := yama.Tumo()
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, &hai.Manzu1, outHai)
	}
}

func TestKanDora(t *testing.T) {
	cases := []struct {
		beforeWanHai   []*hai.Hai
		afterWanHai    []*hai.Hai
		afterOmoteDora []*hai.Hai
		afterUraDora   []*hai.Hai
		outError       error
	}{
		{
			beforeWanHai:   all[122:],
			afterWanHai:    all[124:],
			afterOmoteDora: []*hai.Hai{&hai.Sha},
			afterUraDora:   []*hai.Hai{&hai.Pei},
			outError:       nil,
		},
		{
			beforeWanHai: []*hai.Hai{},
			outError:     YamaNoMoreHaiErr,
		},
	}

	for _, c := range cases {

		yama := yamaImpl{wanHai: c.beforeWanHai}
		err := yama.KanDora()
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterWanHai, yama.wanHai)
		assert.Equal(t, c.afterOmoteDora, yama.omoteDora)
		assert.Equal(t, c.afterUraDora, yama.uraDora)
	}

}
