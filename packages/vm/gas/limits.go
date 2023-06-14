package gas

import (
	"errors"
	"fmt"
	"io"

	"github.com/iotaledger/wasp/packages/util/rwutil"
)

type Limits struct {
	MaxGasPerBlock         uint64 `json:"maxGasPerBlock" swagger:"desc(The maximum gas per block),required"`
	MinGasPerRequest       uint64 `json:"minGasPerRequest" swagger:"desc(The minimum gas per request),required"`
	MaxGasPerRequest       uint64 `json:"maxGasPerRequest" swagger:"desc(The maximum gas per request),required"`
	MaxGasExternalViewCall uint64 `json:"maxGasExternalViewCall" swagger:"desc(The maximum gas per external view call),required"`
}

var LimitsDefault = &Limits{
	MaxGasPerBlock:         1_000_000_000,
	MinGasPerRequest:       10_000,
	MaxGasPerRequest:       50_000_000, // 20 requests per block max
	MaxGasExternalViewCall: 50_000_000,
}

func LimitsFromBytes(data []byte) (*Limits, error) {
	return rwutil.ReaderFromBytes(data, new(Limits))
}

func (gl *Limits) IsValid() bool {
	if gl.MaxGasPerBlock == 0 {
		return false
	}
	if gl.MinGasPerRequest == 0 || gl.MinGasPerRequest > gl.MaxGasPerBlock {
		return false
	}
	if gl.MaxGasPerRequest < gl.MinGasPerRequest {
		return false
	}
	if gl.MaxGasExternalViewCall == 0 {
		return false
	}
	return true
}

func (gl *Limits) Bytes() []byte {
	return rwutil.WriterToBytes(gl)
}

func (gl *Limits) String() string {
	return fmt.Sprintf(
		"GasLimits(max/block: %d, min/req: %d, max/req: %d, max/view: %d",
		gl.MaxGasPerBlock,
		gl.MaxGasPerBlock,
		gl.MinGasPerRequest,
		gl.MaxGasExternalViewCall,
	)
}

func (gl *Limits) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	gl.MaxGasPerBlock = rr.ReadUint64()
	gl.MinGasPerRequest = rr.ReadUint64()
	gl.MaxGasPerRequest = rr.ReadUint64()
	gl.MaxGasExternalViewCall = rr.ReadUint64()
	if rr.Err == nil && !gl.IsValid() {
		rr.Err = errors.New("invalid gas limits")
	}
	return rr.Err
}

func (gl *Limits) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteUint64(gl.MaxGasPerBlock)
	ww.WriteUint64(gl.MinGasPerRequest)
	ww.WriteUint64(gl.MaxGasPerRequest)
	ww.WriteUint64(gl.MaxGasExternalViewCall)
	return ww.Err
}
