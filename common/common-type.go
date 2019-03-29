package common

const(
	FP_AUDIO = "audio/"
	FP_OTHER = "other/"

	FT_AUDIO_WAV = "audio/wav"
	FT_AUDIO_MP3 = "audio/mp3"
	FT_AUDIO_OPUS = "audio/opus"
	FT_AUDIO_PCM = "audio/pcm"
)

type FSUploadParam struct {
	FileName string `form:"tex" json:"file_name" binding:"required"`
	Expire int64	`form:"tex" json:"expire" binding:"omitempty"`
	Type string `form:"tex" json:"type" binding:"required,eq=audio|eq=text|eq=video|eq=raw|eq=other"`
}

type FSUploadResult struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}

type FileInfo struct {
	CreateTime int64
	Name string
	Expire int64
	Type string
}

func (f *FileInfo) IsExpire(now int64) bool {
	return now < f.CreateTime + f.Expire
}