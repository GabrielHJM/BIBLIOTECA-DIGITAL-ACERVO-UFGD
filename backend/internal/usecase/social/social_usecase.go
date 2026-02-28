package social

import (
	"biblioteca-digital-api/internal/domain/social"
)

type socialUseCase struct {
	likeRepo    social.LikeRepository
	commentRepo social.ComentarioRepository
	shareRepo   social.ShareLogRepository
	messageRepo social.MensagemRepository
}

func NewSocialUseCase(
	lr social.LikeRepository,
	cr social.ComentarioRepository,
	sr social.ShareLogRepository,
	mr social.MensagemRepository,
) social.UseCase {
	return &socialUseCase{
		likeRepo:    lr,
		commentRepo: cr,
		shareRepo:   sr,
		messageRepo: mr,
	}
}

func (u *socialUseCase) ToggleLike(usuarioID, materialID int) (bool, error) {
	liked, err := u.likeRepo.HasLiked(usuarioID, materialID)
	if err != nil {
		return false, err
	}

	if liked {
		if err := u.likeRepo.RemoveLike(usuarioID, materialID); err != nil {
			return false, err
		}
		return false, nil
	}

	if err := u.likeRepo.AddLike(usuarioID, materialID); err != nil {
		return false, err
	}
	return true, nil
}

func (u *socialUseCase) AddComment(c social.Comentario) error {
	return u.commentRepo.AddComentario(c)
}

func (u *socialUseCase) RegisterShare(s social.ShareLog) error {
	return u.shareRepo.AddShareLog(s)
}

func (u *socialUseCase) SendMessage(m social.Mensagem) error {
	return u.messageRepo.SendMensagem(m)
}
