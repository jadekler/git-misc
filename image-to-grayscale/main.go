package main

import (
    "bufio"
    "fmt"
    "image"
    "image/png"
    "image/draw"
    "log"
    "os"
    "time"
    "reflect"

    _ "image/jpeg"
)

var bench_total float64

func benchmark(comment string, fun func() interface{}) interface{} {
    t := time.Now()
    result := fun()
    delta_ms := float64(time.Now().Sub(t).Nanoseconds()) / 1000 / 1000
    bench_total += delta_ms
    fmt.Printf("%s took %g ms\n", comment, delta_ms)
    return result
}

func printSummary() {
    fmt.Println("---")
    fmt.Printf("Total time: %g ms\n", bench_total)
}


func openImage(path string) image.Image {
    // Open the file.
    file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Decode the image.
    return benchmark("Image decode", func() (img interface{}) {
        img, err := png.Decode(file)
        if err != nil {
            log.Fatal(err)
        }
        return
    }).(image.Image)
}

func saveImage(image image.Image, path string) {
    // Create or truncate the file.
    out_file, err := os.Create(path)
    if err != nil {
        log.Fatal(err)
    }
    defer out_file.Close()

    // Save the image
    benchmark("Image save", func() (interface{}) {
        err = png.Encode(out_file, image)
        if err != nil {
            log.Fatal(err)
        }
        return ""
    })
}

func newGray(bounds image.Rectangle) *image.Gray {
    return benchmark("Creating an empty image", func() interface{} {
        return image.NewGray(bounds)
    }).(*image.Gray)
}

func convertDraw(img image.Image, dst *image.Gray) {
    benchmark("Drawing the image", func() interface{} {
        draw.Draw(dst, dst.Bounds(), img, image.ZP, draw.Src)
        return nil
    })
}

func convertLoop(img image.Image, dst *image.Gray) {
    benchmark("Converting pixels", func() interface{} {
        log.Println("\n"+reflect.TypeOf(img))

        // We're ignoring alpha here because it's arguable that pixels should
        // be multiplied by its value. It would make more sense to output
        // Gray+Alpha if the input image is RGBA.
        rgba := img.(*image.RGBA64)
        for i := 0; i < len(dst.Pix); i++ {
            src_i := i * 4
            r := uint32(rgba.Pix[src_i])
            g := uint32(rgba.Pix[src_i+1])
            b := uint32(rgba.Pix[src_i+2])
            dst.Pix[i] = uint8((r*299 + g*587 + b*114) / 1000)

        }
        return nil
    })
}

func convertToPNG(file string) {
    reader, err := os.Open(file)
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()

    m, _, err := image.Decode(reader)
    if err != nil {
        log.Fatal(err)
    }

    fw, err := os.Create("zoidberg.png")
    defer fw.Close()
    if err != nil {
        log.Fatal(err)
    }
    writer := bufio.NewWriter(fw)

    png.Encode(writer, m)

    writer.Flush()
}

func main() {
    for i := 0; i < 2; i++ {
        bench_total = 0

        convertToPNG("zoidberg.jpg")
        image := openImage("zoidberg.png")

        grayscale := newGray(image.Bounds())
        if i == 0 {
            convertDraw(image, grayscale)
        } else {
            convertLoop(image, grayscale)
        }
        saveImage(grayscale, "go_output.png")

        printSummary()
        fmt.Println()
    }
}