package feessync

import (
	"errors"
	"sync"
	"time"

	"github.com/zhangpanyi/botcasino/config"
	"github.com/zhangpanyi/botcasino/remote"
	"github.com/zhangpanyi/botcasino/storage"

	"github.com/zhangpanyi/basebot/logger"
)

// 转账费用
var fees sync.Map
var once sync.Once

// CheckFeeStatusAsync 检查转账费用
func CheckFeeStatusAsync() {
	once.Do(func() {
		if err := syncTransferFees(); err != nil {
			panic(err)
		}

		go func() {
			timer := time.NewTimer(time.Second * 60)
			for {
				select {
				case <-timer.C:
					err := syncTransferFees()
					if err != nil {
						logger.Warnf("Failed to sync transfer fees, %v", err)
					} else {
						logger.Info("Transfer fees synchronized")
					}
				}
			}
		}()
	})
}

// GetFee 获取手续费
func GetFee(asset string) (uint32, error) {
	// 获取固定手续费
	dynamicCfg := config.GetDynamic()
	if dynamicCfg.UseFixedFee {
		switch asset {
		case storage.BitCNYSymbol:
			return dynamicCfg.FixedFeeAmount.CNY, nil
		case storage.BitUSDSymbol:
			return dynamicCfg.FixedFeeAmount.USD, nil
		default:
			return 0, errors.New("not found")
		}
	}

	// 获取动态手续费
	val, ok := fees.Load(asset)
	if !ok {
		return 0, errors.New("not found")
	}
	return val.(uint32), nil
}

// 同步转账手续费
func syncTransferFees() error {
	assets := [...]string{storage.BitCNYSymbol, storage.BitUSDSymbol}
	result, err := remote.GetFees(assets[:])
	if err != nil {
		return err
	}
	if len(assets) != len(result) {
		return errors.New("invalid result")
	}
	for idx, asset := range assets {
		if result[idx] == 0 {
			fees.Store(asset, 1)
		} else {
			fees.Store(asset, result[idx])
		}
	}
	return nil
}
