package cmd

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

type HashType uint8

const (
	SHA1 HashType = iota
	SHA256
)

func hash_check(hash string, hash_type HashType) error {

	hash = strings.ToLower(hash)
	hash_length := len(hash)

	switch hash_type {

	case SHA1:

		if hash_length != 40 {
			return fmt.Errorf("error: SHA1(GitHub) hash must be 40 characters long")
		}

	case SHA256:

		if hash_length != 64 {
			return fmt.Errorf("error: SHA256 hash must be 64 characters long")
		}

	default:
		return fmt.Errorf("error: unknown hash type")

	}

	for idx := 0; idx < hash_length; idx++ {

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

	// Buffer size set to 64kb
	buffer := make([]byte, 65536)

	for {

		n, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			return "", fmt.Errorf("error: an error occured whlie doing hash checksum \n%s", err)
		}

		if n == 0 {
			break
		}

		hash.Write(buffer[:n])

	}

	result := fmt.Sprintf("%x", hash.Sum(nil))

	return result, nil
}

func hex_to_RGBA(color_hex string) (color_struct color.RGBA, err error) {

	color_struct.A = 0xff
	_, err = fmt.Sscanf(color_hex, "%02x%02x%02x", &color_struct.R, &color_struct.G, &color_struct.B)

	return color_struct, err

}

func image_from_hash(hash string, filename string, hash_type HashType, size int, theme_map *map[string][]string, theme_name string) error {

	palette, exists := (*theme_map)[theme_name]

	if !exists {
		return fmt.Errorf("error: theme \"%s\" is not defined in .shavis-go.yaml file", theme_name)
	}

	if len(palette) != 16 {
		return fmt.Errorf("error: a theme must have 16 colors, got %d in \"%s\" theme in .shavis-go.yaml file", len(palette), theme_name)
	}

	var width int
	var height int

	switch hash_type {

	case SHA1:
		width, height = 8, 5

	case SHA256:
		width, height = 8, 8

	default:
		return fmt.Errorf("error: unknown hash type")

	}

	top_left := image.Point{0, 0}
	bottom_right := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{top_left, bottom_right})

	color_palette := make([]color.RGBA, 16)

	for idx := 0; idx < 16; idx++ {

		color_hex := palette[idx]
		color_struct, err := hex_to_RGBA(color_hex)

		if err != nil {
			return fmt.Errorf("error: Cannot convert color #%s to RGBA", color_hex)
		}

		color_palette[idx] = color_struct

	}

	decimal_values := make([]int, 0, len(hash))

	for _, char := range hash {

		result, _ := strconv.ParseInt(string(char), 16, 64)

		if (result < math.MinInt) || (result > math.MaxInt) {
			return fmt.Errorf("error: Integer overflow while parsing hash into decimal values")
		}

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

	src := image.Image(img)

	src = imaging.Resize(src, int(math.Pow(2, float64(size)+2)), 0, imaging.NearestNeighbor)
	err := imaging.Save(src, filename)

	if err != nil {
		return fmt.Errorf("error: Failed to save image: %v", err)
	}

	return nil
}
