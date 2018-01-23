package storage_test

import (
	"log"
	"testing"

	"github.com/zhangpanyi/botcasino/envelopes/algorithm"
	"github.com/zhangpanyi/botcasino/storage"
)

// 测试创建红包
func TestNewRedEnvelope(t *testing.T) {
	storage.Connect("test.db")

	number := 1
	envelopes, err := algorithm.Generate(10000, uint32(number))
	if err != nil {
		t.Fatal(err)
	}

	redEnvelope := &storage.RedEnvelope{
		Asset:  "bitCNY",
		Amount: 100,
		Number: uint32(number),
	}
	handler := storage.RedEnvelopeStorage{}
	redEnvelope, err = handler.NewRedEnvelope(redEnvelope, envelopes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(redEnvelope)
}

// 测试获取红包信息
func TestGetRedEnvelope(t *testing.T) {
	storage.Connect("test.db")

	handler := storage.RedEnvelopeStorage{}
	redEnvelope, received, err := handler.GetRedEnvelope(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(redEnvelope, received)
}

// 测试领取红包
func TestReceiveRedEnvelope(t *testing.T) {
	storage.Connect("test.db")

	handler := storage.RedEnvelopeStorage{}
	for i := 0; i < 100; i++ {
		amount, number, err := handler.ReceiveRedEnvelope(1, int32(i), "zpy")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(amount, number)
	}
}

// 测试获取极端红包
func TestGetTwoTxtremes(t *testing.T) {
	storage.Connect("test.db")

	handler := storage.RedEnvelopeStorage{}
	min, max, err := handler.GetTwoTxtremes(100002)
	if err != nil {
		log.Fatalln(err)
	}
	t.Log(min, max)
}

// 测试遍历红包
func TestForeachRedEnvelopes(t *testing.T) {
	storage.Connect("test.db")

	handler := storage.RedEnvelopeStorage{}
	err := handler.ForeachRedEnvelopes(100043, func(data *storage.RedEnvelope) {
		log.Println(data)
	})
	if err != nil {
		log.Fatalln(err)
	}
}
