package utils

// GetOffsetOrDefault offset値を取得、nilの場合はデフォルト値(0)を返す
func GetOffsetOrDefault(offset *int32) int32 {
	if offset != nil {
		return *offset
	}
	return 0
}

// GetLimitOrDefault limit値を取得、nilの場合はデフォルト値(20)を返す
func GetLimitOrDefault(limit *int32) int32 {
	if limit != nil {
		return *limit
	}
	return 20
}

// StringPtr 文字列のポインタを返す
func StringPtr(s string) *string {
	return &s
}

// Float32Ptr float32のポインタを返す
func Float32Ptr(f float32) *float32 {
	return &f
}

// Int32Ptr int32のポインタを返す
func Int32Ptr(i int32) *int32 {
	return &i
}
