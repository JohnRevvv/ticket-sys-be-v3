package models

type UploadResponse struct {
	FileName string `json:"file_name"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	URL      string `json:"url"`
}