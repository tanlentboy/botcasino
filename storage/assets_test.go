package storage_test

import (
	"testing"

	"github.com/zhangpanyi/botcasino/storage"
)

// 测试充值
func TestDeposit(t *testing.T) {
	storage.Connect("test.db")

	var handler storage.AssetStorage
	err := handler.Deposit(1024, "BitCNY", 1000)
	if err != nil {
		t.Fatal(err)
	}
}

// 测试获取资产信息
func TestGetAsset(t *testing.T) {
	storage.Connect("test.db")

	var handler storage.AssetStorage
	asset, err := handler.GetAsset(2048, "BitCNY")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(asset)
}

// 测试获取所有资产信息
func TestGetAssets(t *testing.T) {
	storage.Connect("test.db")

	var handler storage.AssetStorage
	assets, err := handler.GetAssets(1024)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(assets)
}

// 测试资产提现
func TestWithdraw(t *testing.T) {
	storage.Connect("test.db")

	var handler storage.AssetStorage
	err := handler.Withdraw(1024, "BitCNY", 10)
	if err != nil {
		t.Fatal(err)
	}
}

// 测试冻结资产
func TestFrozenAsset(t *testing.T) {
	storage.Connect("test.db")

	var handler storage.AssetStorage
	err := handler.FrozenAsset(1024, "BitCNY", 500)
	if err != nil {
		t.Fatal(err)
	}
}

// 测试解冻资产
func TestUnfreezeAsset(t *testing.T) {
	storage.Connect("test.db")

	var handler storage.AssetStorage
	err := handler.UnfreezeAsset(1024, "BitCNY", 100)
	if err != nil {
		t.Fatal(err)
	}
}

// 测试转移冻结资产
func TestTransferFrozenAsset(t *testing.T) {
	storage.Connect("test.db")

	var handler storage.AssetStorage
	err := handler.TransferFrozenAsset(1024, 2048, "BitCNY", 100)
	if err != nil {
		t.Fatal(err)
	}
}
