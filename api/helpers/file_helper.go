package helpers

import (
	"context"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func ImageUploadHelper(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cld, _ := cloudinary.NewFromParams("ddeatrwxs", "641461484529384", "Pg62Uh1szh18LjcCjCHA14z_oA8")
	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: "studapp"})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}
