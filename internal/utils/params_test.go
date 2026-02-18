package utils

import "testing"

func TestGetOffsetOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		offset   *int32
		expected int32
	}{
		{
			name:     "nilの場合はデフォルト値0を返す",
			offset:   nil,
			expected: 0,
		},
		{
			name:     "値が設定されている場合はその値を返す",
			offset:   int32Ptr(10),
			expected: 10,
		},
		{
			name:     "0が設定されている場合は0を返す",
			offset:   int32Ptr(0),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetOffsetOrDefault(tt.offset)
			if result != tt.expected {
				t.Errorf("GetOffsetOrDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetLimitOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		limit    *int32
		expected int32
	}{
		{
			name:     "nilの場合はデフォルト値20を返す",
			limit:    nil,
			expected: 20,
		},
		{
			name:     "値が設定されている場合はその値を返す",
			limit:    int32Ptr(50),
			expected: 50,
		},
		{
			name:     "最小値1が設定されている場合は1を返す",
			limit:    int32Ptr(1),
			expected: 1,
		},
		{
			name:     "最大値100が設定されている場合は100を返す",
			limit:    int32Ptr(100),
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLimitOrDefault(tt.limit)
			if result != tt.expected {
				t.Errorf("GetLimitOrDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}
