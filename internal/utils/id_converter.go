package utils

import (
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
)

// Int32ToUUID int32のIDを決定的なUUIDに変換
// Note: 本来はDBスキーマをUUIDに変更すべきだが、
// 既存のint32スキーマとの互換性のため、一時的な変換を行う
func Int32ToUUID(id int32) uuid.UUID {
	// UUID v5 (名前ベース) の代わりに、
	// int32を16バイトのUUIDに埋め込む決定的な変換を行う
	var uuidBytes [16]byte
	
	// 最初の4バイトにint32を格納（ビッグエンディアン）
	binary.BigEndian.PutUint32(uuidBytes[0:4], uint32(id))
	
	// バージョン4の形式にする（ランダムではないが）
	uuidBytes[6] = (uuidBytes[6] & 0x0f) | 0x40 // Version 4
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80 // Variant RFC4122
	
	result, _ := uuid.FromBytes(uuidBytes[:])
	return result
}

// UUIDToInt32 UUIDをint32のIDに変換
func UUIDToInt32(id uuid.UUID) (int32, error) {
	uuidBytes := [16]byte(id)
	
	// 最初の4バイトからint32を取得
	id32 := int32(binary.BigEndian.Uint32(uuidBytes[0:4]))
	
	if id32 <= 0 {
		return 0, apperror.BadRequestError("無効なIDフォーマットです").
			WithDetails("uuid", id.String()).
			WithDetails("extracted_id", id32)
	}
	
	return id32, nil
}

// MustUUIDToInt32 UUIDをint32に変換（エラー時はパニック）
// テストやデバッグ用
func MustUUIDToInt32(id uuid.UUID) int32 {
	id32, err := UUIDToInt32(id)
	if err != nil {
		panic(fmt.Sprintf("failed to convert UUID to int32: %v", err))
	}
	return id32
}
