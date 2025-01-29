// This file is part of Crop.
// Copyright (C) 2025 Enindu Alahapperuma
//
// Crop is free software: you can redistribute it and/or modify it under the
// terms of the GNU General Public License as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// Crop is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
// A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with
// Crop. If not, see <https://www.gnu.org/licenses/>.

// Crop is a simple command-line application for cropping images.
//
// Usage:
//
//	crop [path] [width] [height]
package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("usage: %s [path] [width] [height]\n", os.Args[0])
		return
	}

	inputFilePath := os.Args[1]

	targetWidth, err := strconv.ParseInt(os.Args[2], 10, 0)
	if err != nil {
		fmt.Printf("define target width: %q\n", err)
		return
	}

	targetHeight, err := strconv.ParseInt(os.Args[3], 10, 0)
	if err != nil {
		fmt.Printf("define target height: %q\n", err)
		return
	}

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Printf("open file: %q\n", err)
		return
	}

	defer inputFile.Close()

	inputImage, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Printf("decode image: %q\n", err)
		return
	}

	inputWidth := inputImage.Bounds().Dx()
	inputHeight := inputImage.Bounds().Dy()
	inputAspectRatio := float64(inputWidth) / float64(inputHeight)
	targetAspectRation := float64(targetWidth) / float64(targetHeight)
	cropWidth := inputWidth
	cropHeight := inputHeight

	if inputAspectRatio > targetAspectRation {
		cropWidth = int(float64(cropHeight) * targetAspectRation)
	} else {
		cropHeight = int(float64(cropWidth) / targetAspectRation)
	}

	left := (inputWidth - cropWidth) / 2
	top := (inputHeight - cropHeight) / 2
	right := left + cropWidth
	bottom := top + cropHeight
	cropRectangle := image.Rect(left, top, right, bottom)
	croppedImage := image.NewRGBA(cropRectangle)

	draw.Draw(croppedImage, croppedImage.Bounds(), inputImage, cropRectangle.Min, draw.Src)

	resizedImage := image.NewRGBA(image.Rect(0, 0, int(targetWidth), int(targetHeight)))

	draw.ApproxBiLinear.Scale(resizedImage, resizedImage.Bounds(), croppedImage, croppedImage.Bounds(), draw.Over, nil)

	inputFileDirectory := filepath.Dir(inputFilePath)
	inputFileName := filepath.Base(inputFilePath)
	inputFileExtension := filepath.Ext(inputFilePath)

	outputFile, err := os.Create(fmt.Sprintf("%s/%s-%dx%d.png", inputFileDirectory, strings.TrimSuffix(inputFileName, inputFileExtension), targetWidth, targetHeight))
	if err != nil {
		fmt.Printf("create file: %q\n", err)
		return
	}

	defer outputFile.Close()

	err = png.Encode(outputFile, resizedImage)
	if err != nil {
		fmt.Printf("encode image: %q\n", err)
		return
	}
}
