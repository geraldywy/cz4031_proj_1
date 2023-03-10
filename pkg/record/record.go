package record

import (
	"github.com/geraldywy/cz4031_proj1/pkg/consts"
	"github.com/geraldywy/cz4031_proj1/pkg/utils"
)

type Record interface {
	Serialize() []byte
	TConst() string
	AvgRating() float32
	NumVotes() int32
}

var _ Record = (*recordImpl)(nil)

func NewRecord(tconst string, avgRating float32, numVotes int32) Record {
	return &recordImpl{
		tconst:        tconst,
		averageRating: avgRating,
		numVotes:      numVotes,
	}
}

func NewRecordFromBytes(buf []byte) Record {
	if buf == nil {
		return nil
	}
	x := 1
	if buf[x] == 0 {
		x++
	}
	return &recordImpl{
		tconst:        string(buf[x:11]),
		averageRating: utils.Float32FromBytes(utils.SliceTo4ByteArray(buf[11:15])),
		numVotes:      utils.Int32FromBytes(utils.SliceTo4ByteArray(buf[15:19])),
	}
}

type recordImpl struct {
	tconst        string // fixed size string, size 10 ascii characters only
	averageRating float32
	numVotes      int32
}

func (r *recordImpl) Serialize() []byte {
	buf := make([]byte, consts.RecordSize)
	buf[0] = consts.RecordSize
	j := 1
	if len(r.tconst) < 10 {
		buf = append(buf, 0)
		j++
	}
	for i := range r.tconst {
		buf[j] = r.tconst[i]
		j += 1
	}
	for _, b := range utils.Float32ToBytes(r.averageRating) {
		buf[j] = b
		j += 1
	}
	for _, b := range utils.Int32ToBytes(r.numVotes) {
		buf[j] = b
		j += 1
	}

	return buf
}

func (r *recordImpl) TConst() string {
	return r.tconst
}

func (r *recordImpl) AvgRating() float32 {
	return r.averageRating
}

func (r *recordImpl) NumVotes() int32 {
	return r.numVotes
}
