package BitExport

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

type Bits struct {
	bytes     []byte
	bit_count int
}

func (bits *Bits) PushBit(val int) {
	if (bits.bit_count / 8) >= len(bits.bytes) {
		bits.bytes = append(bits.bytes, 0)
	}

	offset := bits.bit_count % 8

	if val == 1 {
		current_byte := (bits.bytes[bits.bit_count/8] | (1 << (offset)))
		bits.bytes[bits.bit_count/8] = current_byte
	}

	bits.bit_count++
}

func BitCount(obj reflect.StructField) int {
	rf_tag := obj.Tag.Get("bits")
	if rf_tag != "" {
		val, _ := strconv.Atoi(rf_tag)
		return val
	}
	return int(obj.Type.Size()) * 8
}

func GetFieldBytes(obj reflect.Value) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, obj.Interface())
	return buf.Bytes()
}

func isPointer(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Pointer
}

func ByteCount(obj interface{}) int {
	bit_count := float64(0)
	rf_type := reflect.TypeOf(obj)
	for i := 0; i < rf_type.NumField(); i++ {
		rf_field := rf_type.Field(i)
		bit_count += float64(BitCount(rf_field))
	}
	return int(math.Ceil(bit_count / 8.0))
}

func ToBytes(obj interface{}) []byte {
	var bits Bits
	rf_obj := reflect.New(reflect.ValueOf(obj).Type()).Elem()
	rf_obj.Set(reflect.ValueOf(obj))
	rf_type := reflect.TypeOf(obj)

	for i := 0; i < rf_obj.NumField(); i++ {
		rf_field := rf_obj.Field(i)
		rf_field = reflect.NewAt(rf_field.Type(), unsafe.Pointer(rf_field.UnsafeAddr())).Elem()
		field_bytes := GetFieldBytes(rf_field)
		abs_byte := 0

		for abs_bit := 0; abs_bit < BitCount(rf_type.Field(i)); abs_bit++ {
			bit := int((field_bytes[abs_bit/8] >> (abs_bit % 8)) & 1)
			bits.PushBit(bit)

			// ?????
			if abs_bit%8 == 0 && abs_bit != 0 {
				abs_byte++
			}
		}
	}

	return bits.bytes
}

func FromBytes(byte_src []byte, dest interface{}) interface{} {
	if !isPointer(dest) {
		log.Fatal("Must Pass a Pointer to FromBytes")
	}
	rf_obj := reflect.New(reflect.ValueOf(dest).Type()).Elem()
	rf_obj.Set(reflect.ValueOf(dest))
	rf_obj = reflect.Indirect(rf_obj)
	rf_type := reflect.TypeOf(dest)

	current_byte := 0
	abs_bit := 0

	for i := 0; i < rf_obj.NumField(); i++ {
		var field_bytes Bits
		rf_field := rf_obj.Field(i)
		rf_field = reflect.NewAt(rf_field.Type(), unsafe.Pointer(rf_field.UnsafeAddr())).Elem()
		for j := 0; j < BitCount(rf_type.Elem().Field(i)); j++ {
			bit := int((byte_src[abs_bit/8] >> (abs_bit % 8)) & 1)
			field_bytes.PushBit(bit)
			if abs_bit%8 == 0 && abs_bit != 0 {
				current_byte++
			}
			abs_bit++
		}

		arr_type := reflect.New(reflect.ArrayOf(int(rf_type.Elem().Field(i).Type.Size()), reflect.TypeOf(byte(0)))).Elem().Type()
		test := reflect.NewAt(arr_type, unsafe.Pointer(rf_field.UnsafeAddr())).Elem()
		for j, by := range field_bytes.bytes {
			test.Index(j).Set(reflect.ValueOf(by))
		}
	}
	return dest
}

// Should rewrite this so that the two could be bit packed together
func MultipleToBytes(objs ...interface{}) []byte {
	bytes := make([]byte, 0)

	for _, obj := range objs {
		bytes = append(bytes, ToBytes(obj)...)
	}
	return bytes
}
