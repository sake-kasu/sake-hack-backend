package presentation

import "github.com/sake-kasu/sake-hack-backend/internal/apperror"

// maxTokenLength はトークンの最大長
const maxTokenLength = 255

// validateLikeToken はX-Like-Tokenヘッダーの値をバリデーションする
func validateLikeToken(token string) error {
	verr := apperror.NewValidationError("いいねリクエストが不正です")

	if token == "" {
		verr = verr.AddField("X-Like-Token", "トークンは必須です")
	} else if len(token) > maxTokenLength {
		verr = verr.AddField("X-Like-Token", "トークンは255文字以内である必要があります")
	}

	if verr.HasErrors() {
		return verr
	}
	return nil
}
