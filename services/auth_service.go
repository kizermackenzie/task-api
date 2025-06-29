package services

type AuthService interface {
	Register(dto RegisterDTO) (*AuthResponseDTO, error)
	Login(dto LoginDTO) (*AuthResponseDTO, error)
	RefreshToken(dto RefreshTokenDTO) (*AuthResponseDTO, error)
	GetUserProfile(userID uint) (*UserResponseDTO, error)
}