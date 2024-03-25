package cmd

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

func hash_check(hash string, hash_type string) error {

	hash = strings.ToLower(hash)

	switch hash_type {

	case "SHA1":

		if len(hash) != 40 {
			return fmt.Errorf("error: SHA1(GitHub) hash must be 40 characters long")
		}

	case "SHA256":

		if len(hash) != 64 {
			return fmt.Errorf("error: SHA256 hash must be 64 characters long")
		}

	}

	for idx := 0; idx < len(hash); idx++ {

		character := hash[idx]

		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return fmt.Errorf("error: Invalid hashsum in hash string\n%s\n%s", hash, strings.Repeat(" ", idx)+"↑")
		}

	}

	return nil

}

func filedata_to_hash(filename string) (string, error) {

	hash := sha256.New()
	file, err := os.Open(filename)

	if err != nil {
		return "", fmt.Errorf("error: an error occured while opening file \n%s", err)
	}

	defer file.Close()

	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("error: an error occured whlie doing hash checksum \n%s", err)
	}

	result := fmt.Sprintf("%x", hash.Sum(nil))

	return result, nil
}

func hex_to_RGBA(color_hex string) (color_struct color.RGBA, err error) {

	color_struct.A = 0xff
	_, err = fmt.Sscanf(color_hex, "%02x%02x%02x", &color_struct.R, &color_struct.G, &color_struct.B)

	return color_struct, err

}

func image_from_hash(hash string, filename string, width int, height int, size int, theme_map map[string][]string, theme_name string) error {

	palette, exists := theme_map[theme_name]

	if !exists {
		return fmt.Errorf("error: theme \"%s\" is not defined in .shavis-go.yaml file", theme_name)
	}

	if len(palette) != 16 {
		return fmt.Errorf("error: a theme must have 16 colors, got %d in \"%s\" theme in .shavis-go.yaml file", len(palette), theme_name)
	}

	top_left := image.Point{0, 0}
	bottom_right := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{top_left, bottom_right})

	var color_palette []color.RGBA

	for _, color_hex := range palette {

		color_struct, err := hex_to_RGBA(color_hex)

		if err != nil {
			return fmt.Errorf("error: Cannot convert color #%s to RGBA", color_hex)
		}

		color_palette = append(color_palette, color_struct)

	}

	var decimal_values []int

	for _, char := range hash {
		result, _ := strconv.ParseInt(string(char), 16, 64)
		decimal_values = append(decimal_values, int(result))
	}

	var decimal_values_mat [][]int

	for i := 0; i < len(decimal_values); i += 8 {

		end := i + 8

		if end > len(decimal_values) {
			end = len(decimal_values)
		}

		decimal_values_mat = append(decimal_values_mat, decimal_values[i:end])

	}

	for row_idx, row := range decimal_values_mat {
		for col_idx, value := range row {
			img.Set(col_idx, row_idx, color_palette[value])
		}
	}

	f, _ := os.Create(filename)
	png.Encode(f, img)

	src, err := imaging.Open(filename)

	if err != nil {
		return fmt.Errorf("error: Failed to open image: %v", err)
	}

	src = imaging.Resize(src, int(math.Pow(2, float64(size)+2)), 0, imaging.NearestNeighbor)
	err = imaging.Save(src, filename)

	if err != nil {
		return fmt.Errorf("error: Failed to save image: %v", err)
	}

	return nil
}
