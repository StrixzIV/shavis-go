package cmd

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
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
			return fmt.Errorf("Error: SHA1(GitHub) hash must be 40 characters long.")
		}

	case "SHA256":

		if len(hash) != 64 {
			return fmt.Errorf("Error: SHA256 hash must be 64 characters long.")
		}

	}

	for idx := 0; idx < len(hash); idx++ {

		character := hash[idx]

		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return fmt.Errorf("Error: Invalid hashsum in hash string\n%s\n%s", hash, strings.Repeat(" ", idx)+"↑")
		}

	}

	return nil

}

func hex_to_RGBA(color_hex string) (color_struct color.RGBA, err error) {

	color_struct.A = 0xff
	_, err = fmt.Sscanf(color_hex, "%02x%02x%02x", &color_struct.R, &color_struct.G, &color_struct.B)

	return color_struct, err

}

func image_from_hash(hash string, filename string, width int, height int, size int, palette []string) error {

	if size > 8 {
		return fmt.Errorf("Error: size must be an integer between 1 to 8")
	}

	top_left := image.Point{0, 0}
	bottom_right := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{top_left, bottom_right})

	var color_palette []color.RGBA

	for _, color_hex := range palette {

		color_struct, err := hex_to_RGBA(color_hex)

		if err != nil {
			return fmt.Errorf("Error: Cannot convert color #%s to RGBA", color_hex)
		}

		color_palette = append(color_palette, color_struct)

	}

	var decimal_values []int

	for _, char := range hash {
		result, _ := strconv.ParseInt(string(char), 16, 64)
		decimal_values = append(decimal_values, int(result))
	}

	fmt.Println(decimal_values)

	var decimal_values_mat [][]int

	for i := 0; i < len(decimal_values); i += 8 {

		end := i + 8

		if end > len(decimal_values) {
			end = len(decimal_values)
		}

		decimal_values_mat = append(decimal_values_mat, decimal_values[i:end])

	}

	fmt.Println(decimal_values_mat)

	for row_idx, row := range decimal_values_mat {
		for col_idx, value := range row {
			img.Set(col_idx, row_idx, color_palette[value])
		}
	}

	f, _ := os.Create(filename)
	png.Encode(f, img)

	fmt.Println("Main image done")

	src, err := imaging.Open(filename)

	if err != nil {
		return fmt.Errorf("Error: Failed to open image: %v", err)
	}

	src = imaging.Resize(src, int(math.Pow(2, float64(size)+2)), 0, imaging.NearestNeighbor)
	err = imaging.Save(src, filename)

	if err != nil {
		return fmt.Errorf("Error: Failed to save image: %v", err)
	}

	fmt.Println("Resized image done")

	return nil
}
