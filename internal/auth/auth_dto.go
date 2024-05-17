package auth

import "1-cat-social/internal/user"

type NurseLoginRequest struct {
	NIP      string `json:"nip" binding:"required,len=13,number"`
	Password string `json:"password" binding:"required,min=5,max=33"`
}

type NuseLoginResponse struct {
	UserId      string `json:"userId"`
	NIP         string `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

func FormatNurseLoginResponse(nurse user.User, token string) NuseLoginResponse {
	return NuseLoginResponse{
		UserId:      nurse.ID,
		NIP:         nurse.NIP,
		Name:        nurse.Name,
		AccessToken: token,
	}
}
