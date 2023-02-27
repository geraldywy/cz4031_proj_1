package storage

import (
	"github.com/geraldywy/cz4031_proj1/pkg/record"
	"reflect"
	"testing"
)

func Test_storageImpl_ReadRecord(t *testing.T) {
	type fields struct {
		store []block
	}
	tests := []struct {
		name    string
		fields  fields
		ptr     *StoragePointer
		want    record.Record
		wantErr bool
	}{
		{
			name: "Simple read within a block",
			fields: fields{
				store: []block{
					nil,
					nil,
					[]byte{28, 0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					nil,
					nil,
					nil,
				},
			},
			ptr: &StoragePointer{
				BlockPtr:  2,
				RecordPtr: 1,
			},
			want:    record.NewRecordFromBytes([]byte{0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			wantErr: false,
		},
		{
			name: "Read across blocks",
			fields: fields{
				store: []block{
					nil,
					nil,
					[]byte{9, 0, 116, 116, 48, 48, 48, 48, 48},
					[]byte{19, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					nil,
					nil,
				},
			},
			ptr: &StoragePointer{
				BlockPtr:  2,
				RecordPtr: 1,
			},
			want:    record.NewRecordFromBytes([]byte{0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storageImpl{
				store: tt.fields.store,
			}
			got, err := s.ReadRecord(tt.ptr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadRecord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storageImpl_InsertRecord(t *testing.T) {
	type fields struct {
		store       []block
		spaceUsed   int
		maxCapacity int
		blockSize   int
	}
	type args struct {
		record record.Record
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *StoragePointer
		wantErr bool
	}{
		{
			name: "Basic insert within a block",
			fields: fields{
				store: []block{
					nil,
					nil,
					[]byte{28, 0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0,
						0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				spaceUsed:   100,
				maxCapacity: 1000,
				blockSize:   50,
			},
			args: args{
				record: record.NewRecordFromBytes([]byte{0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			},
			want: &StoragePointer{
				BlockPtr:  2,
				RecordPtr: 28,
			},
			wantErr: false,
		},
		{
			name: "Insert across blocks",
			fields: fields{
				store: []block{
					nil,
					nil,
					[]byte{28, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157,
						0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				spaceUsed:   100,
				maxCapacity: 1000,
				blockSize:   30,
			},
			args: args{
				record: record.NewRecordFromBytes([]byte{0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			},
			want: &StoragePointer{
				BlockPtr:  2,
				RecordPtr: 28,
			},
			wantErr: false,
		},
		{
			name: "Insert at new block",
			fields: fields{
				store: []block{
					nil,
					nil,
					[]byte{28, 0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				spaceUsed:   500,
				maxCapacity: 1000,
				blockSize:   28,
			},
			args: args{
				record: record.NewRecordFromBytes([]byte{0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			},
			want: &StoragePointer{
				BlockPtr:  3,
				RecordPtr: 1,
			},
			wantErr: false,
		},
		{
			name: "Insert at the start",
			fields: fields{
				store:       []block{},
				spaceUsed:   0,
				maxCapacity: 1000,
				blockSize:   21,
			},
			args: args{
				record: record.NewRecordFromBytes([]byte{0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			},
			want: &StoragePointer{
				BlockPtr:  0,
				RecordPtr: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storageImpl{
				store:       tt.fields.store,
				spaceUsed:   tt.fields.spaceUsed,
				maxCapacity: tt.fields.maxCapacity,
				blockSize:   tt.fields.blockSize,
			}
			got, err := s.InsertRecord(tt.args.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil && tt.want != nil) || !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertRecord() got = %v, want %v", got, tt.want)
			}
			rec, err := s.ReadRecord(got)
			if err != nil {
				t.Errorf("Failed to read back record inserted, got err %v", err)
			}
			if !reflect.DeepEqual(rec, tt.args.record) {
				t.Errorf("Failed to read back record inserted, got = %v, want %v", rec, tt.args.record)
			}
		})
	}
}

func Test_storageImpl_DeleteRecord(t *testing.T) {
	type fields struct {
		store                 []block
		spaceUsed             int
		maxCapacity           int
		blockSize             int
		lastRecordInsertedPtr *StoragePointer
	}
	type args struct {
		ptr *StoragePointer
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantStore []block
	}{
		{
			name: "delete only record",
			fields: fields{
				store: []block{
					nil,
					nil,
					[]byte{24, 0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0},
				},
				spaceUsed:   24,
				maxCapacity: 30,
				blockSize:   30,
				lastRecordInsertedPtr: &StoragePointer{
					BlockPtr:  2,
					RecordPtr: 1,
				},
			},
			args: args{
				ptr: &StoragePointer{
					BlockPtr:  2,
					RecordPtr: 1,
				},
			},
			wantErr: false,
			wantStore: []block{
				nil,
				nil,
			},
		},
		{
			name: "delete one record within a block",
			fields: fields{
				store: []block{
					nil,
					nil,
					[]byte{47, 0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 2, 1},
				},
				spaceUsed:   47,
				maxCapacity: 500,
				blockSize:   47,
				lastRecordInsertedPtr: &StoragePointer{
					BlockPtr:  2,
					RecordPtr: 24,
				},
			},
			args: args{
				ptr: &StoragePointer{
					BlockPtr:  2,
					RecordPtr: 1,
				},
			},
			wantErr: false,
			wantStore: []block{
				nil,
				nil,
				[]byte{24, 0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		{
			name: "delete one record, update across blocks",
			fields: fields{
				store: []block{
					nil,
					nil,
					[]byte{28, 0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 116, 116, 48},
					[]byte{20, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 2, 1},
				},
				spaceUsed:   48,
				maxCapacity: 100,
				blockSize:   28,
				lastRecordInsertedPtr: &StoragePointer{
					BlockPtr:  2,
					RecordPtr: 24,
				},
			},
			args: args{
				ptr: &StoragePointer{
					BlockPtr:  2,
					RecordPtr: 1,
				},
			},
			wantErr: false,
			wantStore: []block{
				nil,
				nil,
				[]byte{24, 0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		{
			name: "delete one record stored across blocks",
			fields: fields{
				store: []block{
					nil,
					nil,
					[]byte{28, 0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 116, 116, 48},
					[]byte{20, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 2, 1},
				},
				spaceUsed:   48,
				maxCapacity: 100,
				blockSize:   28,
				lastRecordInsertedPtr: &StoragePointer{
					BlockPtr:  2,
					RecordPtr: 24,
				},
			},
			args: args{
				ptr: &StoragePointer{
					BlockPtr:  2,
					RecordPtr: 24,
				},
			},
			wantErr: false,
			wantStore: []block{
				nil,
				nil,
				[]byte{24, 0, 116, 116, 48, 48, 48, 48, 48, 50, 55, 64, 179, 51, 51, 0, 0, 3, 157, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storageImpl{
				store:                 tt.fields.store,
				spaceUsed:             tt.fields.spaceUsed,
				maxCapacity:           tt.fields.maxCapacity,
				blockSize:             tt.fields.blockSize,
				lastRecordInsertedPtr: tt.fields.lastRecordInsertedPtr,
			}
			if err := s.DeleteRecord(tt.args.ptr); (err != nil) != tt.wantErr {
				t.Errorf("DeleteRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(s.store, tt.wantStore) {
				t.Errorf("DeleteRecord() store assertion failed got: %v, want %v", s.store, tt.wantStore)
			}
		})
	}
}
