package attachment

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"nbientry/web"
	"nbientry/web/common"
	"os"
	"strings"
)

var _ = web.Migrate(Attachment{})

var UPLOAD_PATH = "/src/attachment"

type Attachment struct {
	NbiId    string `gorm:"primaryKey"`
	Location string
	Mime     string
	Name     string
}

func (a *Attachment) Exists() bool {
	if a.Location == "" {
		return false
	}

	loc := strings.ReplaceAll(a.Location, "/src/", "./src/")

	stat, err := os.Stat(loc)
	if err != nil || stat.Size() == 0 {
		return false
	}

	return true
}

func (a *Attachment) SaveBytes(file_data []byte) error {
	if !strings.Contains(a.Name, ".") {
		return fmt.Errorf("not a valid filename")
	}

	// Try to convert to an image
	img := common.ImageUpload{Link: "blob"}
	err := img.SaveBytes(file_data)

	if err == nil {
		// Attachment is an image; save as WebP
		a.Location = img.Location
		a.Mime = "image/webp"
		return nil
	}

	file_size := len(file_data)

	sha512_hash := crypto.SHA512.New()
	_, _ = sha512_hash.Write(file_data)
	sha512 := hex.EncodeToString(sha512_hash.Sum(nil))

	file_dir := fmt.Sprintf("%s/%s/%d", sha512[0:2], sha512[2:4], file_size)
	parts := strings.Split(a.Name, ".")
	file_name := sha512 + "." + parts[len(parts)-1]
	file_path := fmt.Sprintf("%s/%s", file_dir, file_name)
	dest_path := fmt.Sprintf(".%s%s/", UPLOAD_PATH, file_dir)
	fmt.Println(file_name)

	err = os.MkdirAll(dest_path, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(dest_path+file_name, file_data, 0644)
	if err != nil {
		return err
	}

	a.Location = UPLOAD_PATH + file_path

	web.Save(a)
	return nil
}
