package utils

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func IsBigIntZero(n *big.Int) bool {
	return len(n.Bits()) == 0
}

func PrintBlock(block *types.Block) {
	t := time.Unix(int64(block.Header().Time), 0).UTC()
	fmt.Printf("%d \t %s \t tx=%-4d \t gas=%d\n", block.Header().Number, t, len(block.Transactions()), block.GasUsed())
}

func DateToTime(dayString string, hour int, min int, sec int) (time.Time, error) {
	dateString := fmt.Sprintf("%sT%02d:%02d:%02dZ", dayString, hour, min, sec)
	return time.Parse(time.RFC3339, dateString)
}

func BigIntToHumanNumberString(i *big.Int, decimals int) string {
	return BigFloatToHumanNumberString(new(big.Float).SetInt(i), decimals)
}

func BigFloatToHumanNumberString(f *big.Float, decimals int) string {
	output := f.Text('f', decimals)
	dotIndex := strings.Index(output, ".")
	if dotIndex == -1 {
		dotIndex = len(output)
	}
	for outputIndex := dotIndex; outputIndex > 3; {
		outputIndex -= 3
		output = output[:outputIndex] + "," + output[outputIndex:]
	}
	return output
}

func NumberToHumanReadableString(value interface{}, decimals int) string {
	switch v := value.(type) {
	case int:
		i := big.NewInt(int64(v))
		return BigIntToHumanNumberString(i, decimals)
	case int64:
		i := big.NewInt(v)
		return BigIntToHumanNumberString(i, decimals)
	case big.Int:
		return BigIntToHumanNumberString(&v, decimals)
	case big.Float:
		return BigFloatToHumanNumberString(&v, decimals)
	case string:
		f, ok := new(big.Float).SetString(v)
		if !ok {
			return v
		}
		return BigFloatToHumanNumberString(f, decimals)
	default:
		return "invalid"
	}
}
